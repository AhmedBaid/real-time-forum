export async function loadUnreadNotifications() {
  const res = await fetch('/unread-messages'); 
  if (!res.ok) return;

  const data = await res.json();

  let chatbox = document.querySelector(`.chat-box[data-id-u="${data.from}"]`);
  if (chatbox) return; 

  const userElement = document.querySelector(`.users[data-id="${data.from}"] .notification`);
  if (userElement) {
    userElement.textContent = data.message; 
  }
}
