import { socket  , Currentusername} from "./home.js";


function throttle(func, time, option = { leading: false, trailing: false }) {
    let wait = false
    return (...args) => {

        if (!wait) {
            if (option.leading) {
                func.apply(this, args)
            }
            wait = true
            setTimeout(() => {
                if (!option.leading && option.trailing) {
                    func.apply(this, args)
                }
                wait = false
            }, time);
        }
    }
}


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
  chatDiv.dataset.idU = receiverId
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

  async function loadMessages(scroll) {
    let res = await fetch(`/messages?receiver=${receiverId}&offset=${offset}`);
    if (!res.ok) throw new Error(await res.text());
    let messages = await res.json();
console.log(messages);

    if (messages.length === 0) return;
let oldScrollHeight = messagesBox.scrollHeight;
    messages.forEach(msg => {
      let div = document.createElement("div");
      div.className = `msg ${msg.receiver === receiverId ? "right" : "left"}`;
      div.innerHTML = `
        <p>${msg.message}</p>
        <span class="time">${msg.senderUsername}-${new Date(msg.time).toLocaleTimeString()}</span>
      `;
      messagesBox.prepend(div);
      /* 
          if (msg.receiver === receiverId) {
            fetch(`/mark-read/${msg.id}`, { method: "POST" });
          } */
    });

    if (scroll) {
      messagesBox.scrollTop = messagesBox.scrollHeight;
    } else {
       messagesBox.scrollTop = messagesBox.scrollHeight - oldScrollHeight;
    }

    console.log(messages);

    offset += 10;
  }

  await loadMessages(true);

  messagesBox.addEventListener("scroll", throttle(async () => {
    if (messagesBox.scrollTop === 0) {
      await loadMessages(false);
    }

  }, 200, {leading:false ,trailing: true }));

  form.addEventListener("submit", (ev) => {
    offset+=1
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
      <p>${input.value}</p>
      <span class="time"> ${Currentusername}-${new Date().toLocaleTimeString()}</span>
    `;
    messagesBox.appendChild(div);

    input.value = "";
    messagesBox.scrollTop = messagesBox.scrollHeight;
  });
}
