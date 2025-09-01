import { socket } from "./home.js";

export async function HandleMessages(e) {
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

  let offset = 0;
  let chatDiv = document.createElement("div");
  chatDiv.className = "chat-box";
  chatDiv.id = `chat-${username}`;
  chatDiv.innerHTML = `
    <div class="chat-header">
      <img src="https://robohash.org/${username}.png?size=50x50" class="avatar" />
      <div><strong>${username}</strong></div>
      <button class="close-btn">✖</button>
    </div>
    <div class="chat-messages"></div>
    <form class="chat-form" method="post">
      <input type="text" placeholder="Type a message..." />
      <button type="submit">➤</button>
    </form>
  `;
  chatArea.appendChild(chatDiv);

  chatDiv.querySelector(".close-btn").onclick = () => {
    chatDiv.remove();
  };

  let form = chatDiv.querySelector(".chat-form");
  let messagesBox = chatDiv.querySelector(".chat-messages");

  async function loadMessages(prepend = false) {
    let res = await fetch(`/messages?receiver=${receiverId}&offset=${offset}`);
    if (!res.ok) throw new Error(await res.text());
    let messages = await res.json();



    if (messages.length === 0) return;

    let oldScrollHeight = messagesBox.scrollHeight;

    messages.forEach(msg => {
      let div = document.createElement("div");
      div.className = `msg ${msg.receiver === receiverId ? "right" : "left"}`;
      div.innerHTML = `
        <p></p>
        <span class="time"></span>
      `;
      let p = div.querySelector("p")
      let span = div.querySelector("span")
      span.textContent = new Date(msg.time).toLocaleTimeString()
      p.textContent = msg.message
      if (prepend) {
        messagesBox.prepend(div);
      } else {
        messagesBox.appendChild(div);
      }

      if (msg.receiver === receiverId) {
        fetch(`/mark-read/${msg.id}`, { method: "POST" });
      }
    });

    if (!prepend) {
      messagesBox.scrollTop = messagesBox.scrollHeight;
    } else {
      messagesBox.scrollTop = messagesBox.scrollHeight - oldScrollHeight;
    }

    offset += 10;
  }

  await loadMessages(false);

  messagesBox.addEventListener("scroll", async () => {
    if (messagesBox.scrollTop === 0) {
      await loadMessages(true);
    }
  });

  form.addEventListener("submit", (ev) => {
    ev.preventDefault();
    let input = form.querySelector("input");
    if (input.value.trim() === "") return;

    socket.send(JSON.stringify({
      type: "message",
      receiver: receiverId,
      message: input.value,
    }));

    let div = document.createElement("div");
    div.className = "msg right";
    div.innerHTML = `
      <p></p>
      <span class="time"></span>
    `;
    let p  =   div.querySelector("p")
    let span  =   div.querySelector("span")
span.textContent=new Date().toLocaleTimeString()
    p.textContent=input.value
    messagesBox.appendChild(div);

    input.value = "";
    messagesBox.scrollTop = messagesBox.scrollHeight;
  });
}
