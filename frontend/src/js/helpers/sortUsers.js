import {  fetchUsers} from "../helpers/api.js";
export async function sortUsers( aside) {
  let users = await fetchUsers();

  users = users.data.sort((a, b) => {
    const aHasMsg = !!a.lastMessageTime;
    const bHasMsg = !!b.lastMessageTime;

    if (aHasMsg && bHasMsg) {
      return new Date(b.lastMessageTime) - new Date(a.lastMessageTime);
    }
    if (aHasMsg) return -1;
    if (bHasMsg) return 1;
    return a.username.localeCompare(b.username);
  });

  for (const user of users) {
    const div = document.createElement("div");
    div.className = "users";
    div.dataset.username = user.username;
    div.dataset.id = user.id;
    div.innerHTML = `
      <img src="https://robohash.org/${user.username}.png?size=50x50" class="avatar" />
        <div class="text-wrapper">

      <span class="username">${user.username}</span>
        <span class="notification"></span>
<span class="typing"> 
<strong>typing</strong>
        <span class="dots">
            <span class="dot">.</span>
            <span class="dot">.</span>
            <span class="dot">.</span>
        </span></span>

  </div>
      <span class="online">.</span>
    `;
    aside.appendChild(div);
  }

}