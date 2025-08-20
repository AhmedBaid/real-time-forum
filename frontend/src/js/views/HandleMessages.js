import {  socket} from "./home.js";

export async function HandleMessages(e) {
  let username = e.currentTarget.dataset.username;
  let recieverId = Number(e.currentTarget.dataset.id);

  let main = document.querySelector(".main");
  let chatArea = document.querySelector(".chat-area");
  if (!chatArea) {
    chatArea = document.createElement("div");
    chatArea.className = "chat-area";
    main.appendChild(chatArea);
  }

  if (document.getElementById(`chat-${username}`)) return;

  try {
    // Fetch messages history
    let res = await fetch(`/getMessages?reciever=${recieverId}`);
    if (!res.ok) throw new Error(await res.text());
    let messages = await res.json();

    // Create chat box
    let chatDiv = document.createElement("div");
    chatDiv.className = "chat-box";
    chatDiv.id = `chat-${username}`;
    chatDiv.innerHTML = `
      <div class="chat-header">
        <span class="avatar">👤</span>
        <div>
          <strong>${username}</strong>
          <div class="status">Online</div>
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

    // Close button
    chatDiv.querySelector(".close-btn").onclick = () => {
      chatDiv.remove();
    };

    let form = chatDiv.querySelector(".chat-form");
    let messagesBox = chatDiv.querySelector(".chat-messages");

    // Display message history
    messagesBox.innerHTML = "";
    messages.forEach(msg => {
      messagesBox.innerHTML += `
        <div class="msg ${msg.reciever === recieverId ? "right" : "left"}">
          <p>${msg.message}</p>
          <span class="time">${new Date(msg.time).toLocaleTimeString()}</span>
        </div>
      `;
    });
    messagesBox.scrollTop = messagesBox.scrollHeight;

    // Form submit
    form.addEventListener("submit", async (ev) => {
      ev.preventDefault();
      try {
        let input = form.querySelector("input");
        if (input.value.trim() === "") return;

        let msg = {
         
          reciever: recieverId,
          message: input.value,
          time: new Date().toISOString()
        };

        // Send to backend
        let resp = await fetch("/sendMessage", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify(msg),
        });
        if (!resp.ok) throw new Error(await resp.text());

        // Append to chat box
        messagesBox.innerHTML += `
          <div class="msg right">
            <p>${input.value}</p>
            <span class="time">Now</span>
          </div>
        `;
        input.value = "";
        messagesBox.scrollTop = messagesBox.scrollHeight;

       

      } catch (error) {
        console.log(error, 5555);
      }
    });

    // Live WebSocket updates
   

  } catch (error) {
    console.log(error);
  }
}
