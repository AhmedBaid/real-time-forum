export async function loadUnreadNotifications() {

  const res = await fetch('/unread-messages'); 
  if (!res.ok) return;
  const data = await res.json();
let chatbox = document.querySelector(`.chat-box[data-id-u="${data.from}"]`);

if (chatbox) {
  return
}
  data.forEach(msg => {
    const userElement = document.querySelector(`.users[data-id="${msg.sender}"] .notification`);
    if (userElement) {
      userElement.textContent = `${msg.count} new Message(s)`;
    }
  });
}
