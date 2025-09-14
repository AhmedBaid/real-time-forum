export async function loadUnreadNotifications() {
  const res = await fetch('/unread-messages');
  if (!res.ok) return;

  const data = await res.json();

  let chatbox = document.querySelector(`.chat-box[data-id-u="${data.from}"]`);
  if (chatbox) {

    return;
  }
  let notif = document.querySelector(".notifIcon");
  if (notif) {

    notif.innerHTML = ` <i class="fa-solid fa-bell bell-icon" id="bellIcon"></i>`
  }
  const userElement = document.querySelector(`.users[data-id="${data.from}"] .notification`);
  if (userElement) {
    userElement.textContent = data.message;
  }
}
