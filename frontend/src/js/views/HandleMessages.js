export async function HandleMessages(e) {
    let username = e.currentTarget.dataset.username;

    let main = document.querySelector(".main");

    let chatArea = document.querySelector(".chat-area");
    if (!chatArea) {
        chatArea = document.createElement("div");
        chatArea.className = "chat-area";
        main.appendChild(chatArea);
    }

    if (document.getElementById(`chat-${username}`)) {
        return;
    }


    let chatDiv = document.querySelector(".chat-box");
    if (!chatDiv) {

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

  <div class="chat-messages">
    <div class="msg left">
      <p>Hey! I loved your feedback on my designs</p>
      <span class="time">Yesterday</span>
    </div>
    <div class="msg right">
      <p>Your work is always inspiring!</p>
      <span class="time">Yesterday</span>
    </div>
  </div>

  <form class="chat-form">
    <input type="text" placeholder="Type a message..." />
    <button type="submit">➤</button>
  </form>
`;

        chatArea.appendChild(chatDiv);

        chatDiv.querySelector(".close-btn").onclick = () => {
            chatDiv.remove();
        };

    }

    chatDiv.innerHTML = `
  <div class="chat-header">
    <span class="avatar">👤</span>
    <div>
      <strong>${username}</strong>
      <div class="status">Online</div>
    </div>
    <button class="close-btn">✖</button>
  </div>

  <div class="chat-messages">
    <div class="msg left">
      <p>Hey! I loved your feedback on my designs</p>
      <span class="time">Yesterday</span>
    </div>
    <div class="msg right">
      <p>Your work is always inspiring!</p>
      <span class="time">Yesterday</span>
    </div>
  </div>

  <form class="chat-form">
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

    form.onsubmit = (ev) => {
        ev.preventDefault();
        let input = form.querySelector("input");
        if (input.value.trim() === "") return;

        messagesBox.innerHTML += `
    <div class="msg right">
      <p>${input.value}</p>
      <span class="time">Now</span>
    </div>
  `;

        input.value = "";
        messagesBox.scrollTop = messagesBox.scrollHeight;
    };

}