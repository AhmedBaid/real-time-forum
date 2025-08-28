import { socket } from "./home.js";

export async function HandleMessages(e) {
  let username = e.currentTarget.dataset.username;
  let receiverId = Number(e.currentTarget.dataset.id);

  let main = document.querySelector(".main");
  let chatArea = document.querySelector(".chat-area");
  if (!chatArea) {
    chatArea = document.createElement("div");
    chatArea.className = "chat-area";
    main.appendChild(chatArea);
  }
chatArea.innerHTML=""
  if (document.getElementById(`chat-${username}`)) return;

  try {
    let res = await fetch(`/messages?receiver=${receiverId}`);
    if (!res.ok) throw new Error(await res.text());
    let messages = await res.json();

    let chatDiv = document.createElement("div");
    chatDiv.className = "chat-box";
    chatDiv.id = `chat-${username}`;
    chatDiv.innerHTML = `
      <div class="chat-header">

<img src="https://robohash.org/${username}.png?size=50x50" class="avatar" />

      <div>
          <strong>${username}</strong>
        </div>
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

    messagesBox.innerHTML = "";
    messages.forEach(msg => {
      messagesBox.innerHTML += `
        <div class="msg ${msg.receiver === receiverId ? "right" : "left"}">
          <p>${msg.message}</p>
          <span class="time">${new Date(msg.time).toLocaleTimeString()}</span>
        </div>
      `;
      if (msg.receiver === receiverId) {
        fetch(`/mark-read/${msg.id}`, { method: "POST" });
      }
    });
    messagesBox.scrollTop = messagesBox.scrollHeight;

    form.addEventListener("submit", (ev) => {
      ev.preventDefault();
      let input = form.querySelector("input");
      if (input.value.trim() === "") return;

      socket.send(JSON.stringify({
        type: "message",
        receiver: receiverId,
        message: input.value,
      }));

      messagesBox.innerHTML += `
        <div class="msg right">
          <p>${input.value}</p>
          <span class="time">${new Date().toLocaleTimeString()}</span>
        </div>
      `;
      input.value = "";
      messagesBox.scrollTop = messagesBox.scrollHeight;
    });

  } catch (error) {
    console.log(error);
  }
}