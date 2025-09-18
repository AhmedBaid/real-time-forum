import { isLogged, Navigate } from "../config.js";
import { showToast } from "../helpers/showToast.js";
import { offset, socket, Currentusername, currentUserId } from "./home.js";
import { login } from "./login.js";

function throttle(func, time, option = { leading: false, trailing: false }) {
  let wait = false;
  return (...args) => {
    if (!wait) {
      if (option.leading) {
        func.apply(this, args);
      }
      wait = true;
      setTimeout(() => {
        if (!option.leading && option.trailing) {
          func.apply(this, args);
        }
        wait = false;
      }, time);
    }
  };
}

export async function HandleMessages(e) {

  offset.nbr = 0
  const user = document.querySelector(`.users[data-id="${e.currentTarget.dataset.id}"] .text-wrapper .notification`)
  if (user.textContent !== "") {
    user.innerHTML = ""
  }

  let username = e.currentTarget.dataset.username;
  let receiverId = Number(e.currentTarget.dataset.id);

  let main = document.querySelector(".main");
  let chatArea = document.querySelector(".chat-area");
  if (!chatArea) {
    chatArea = document.createElement("div");
    chatArea.className = "chat-area";
    main.appendChild(chatArea);
  }
  chatArea.innerHTML = "";
  if (document.getElementById(`chat-${username}`)) return;

  let chatDiv = document.createElement("div");
  chatDiv.className = "chat-box";
  chatDiv.dataset.idU = receiverId;
  chatDiv.id = `chat-${username}`;
  chatDiv.innerHTML = `
    <div class="chat-header">
      <img src="https://robohash.org/${username}.png?size=50x50" class="avatar" />
      <div><strong>${username}</strong></div>
      <button class="close-btn">✖</button>
    </div>
    <div class="chat-messages"></div>
    <span class="chatTyping"><strong></strong>
    <span class="dots2">
    <span class="d1"></span>
    <span class="d2"></span>
    <span class="d3"></span>

    </span>
    </span>
    <form class="chat-form" method="post">
      <input type="text" placeholder="Type a message..." id="input"/>
      <button type="submit">➤</button>
    </form>
  `;
  chatArea.appendChild(chatDiv);

  chatDiv.querySelector(".close-btn").onclick = () => {
    offset.nbr = 0;
    chatDiv.remove();
  };

  let form = chatDiv.querySelector(".chat-form");
  let messagesBox = chatDiv.querySelector(".chat-messages");

  let inputt = form.querySelector("#input");
  let typingTimeout;

  inputt.addEventListener("input", async () => {

    const logged = await isLogged()
    if (!logged) {
      await fetch("/logout", {
        method: "POST",
        credentials: "include"
      });
      login()
      socket.close()
      return
    }
    clearTimeout(typingTimeout);

    socket.send(
      JSON.stringify({
        type: "typing",
        senderUsername: Currentusername,
        senderId: currentUserId,
        receiver: receiverId,
      })
    );

    typingTimeout = setTimeout(() => {
      socket.send(
        JSON.stringify({
          type: "stopTyping",
          senderUsername: Currentusername,
          senderId: currentUserId,
          receiver: receiverId,
        })
      );
    }, 300);
  });

  async function loadMessages(scroll) {
    let res = await fetch(
      `/messages?receiver=${receiverId}&offset=${offset.nbr}`
    );
    if (!res.ok) {
      if (res.status === 401) {
        await fetch("/logout", {
          method: "POST",
          credentials: "include"
        });
        socket.close()
        showToast("error", "you are not authorized");
        Navigate("/login");
        login();
        return;
      }
    }
    let messages = await res.json();

    if (messages.length === 0) return;
    let oldScrollHeight = messagesBox.scrollHeight;
    messages.forEach((msg) => {
      let div = document.createElement("div");
      div.className = `msg ${msg.receiver === receiverId ? "right" : "left"}`;
      div.innerHTML = `
        <p>${msg.message}</p>
        <span class="time">${msg.senderUsername}-${new Date(
        msg.time
      ).toLocaleString()}</span>
      `;
      messagesBox.prepend(div);
    });

    if (scroll) {
      messagesBox.scrollTop = messagesBox.scrollHeight;
    } else {
      messagesBox.scrollTop = messagesBox.scrollHeight - oldScrollHeight;
    }

    offset.nbr += 10;
  }

  await loadMessages(true);

  messagesBox.addEventListener(
    "scroll",
    throttle(
      async () => {
        if (messagesBox.scrollTop === 0) {
          await loadMessages(false);
        }
      },
      200,
      { leading: false, trailing: true }
    )
  );

  form.addEventListener("submit", async (ev) => {
    offset.nbr += 1;
    ev.preventDefault();
    let logged = await isLogged()
    if (!logged) {
      await fetch("/logout", {
        method: "POST",
        credentials: "include"
      });
      socket.close()
      login()
      return
    }
    let input = form.querySelector("input");
    if (input.value.trim() === "") return;

    socket.send(
      JSON.stringify({
        type: "message",
        receiver: receiverId,
        message: input.value,
      })
    );
    input.value = "";

    messagesBox.scrollTop = messagesBox.scrollHeight;
  });
}
