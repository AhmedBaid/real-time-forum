import { container } from "../config.js";

export function home(data) {
    const username = data.username;
    let header = document.createElement("header")
    header.className = "header";
    container.innerHTML = "";
    header.innerHTML = `<h1>Forum Welcome ${username}</h1>
    <form method="post">
        <button class="logout">Logout</button>
    </form>`;
    container.appendChild(header);
    const logoutButton = header.querySelector(".logout");
    logoutButton.addEventListener("click", Logout)
}
function Logout(e) {
    e.preventDefault();
    let response = fetch("/logout")
    if (!response.ok) {
        console.error("Logout failed");
        return;
    }
    home("/login")
}
