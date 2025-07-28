import { container, isLogged, Navigate } from "../config.js";
import { login } from "./login.js";

export function home(data) {
    const username = data.username;
    let header = document.createElement("header")
    container.innerHTML = "";
    header.innerHTML = `<h1>RT-<span>FO</span>RUM</h1>
    <h3>Welcome, ${username}!</h3>
    <button class="logout">Logout</button>`;
    container.appendChild(header);
    const logoutButton = header.querySelector(".logout");
    logoutButton.addEventListener("click", Logout)
}
async function Logout(e) {
    e.preventDefault();
    const response = await fetch("/logout", {
        method: "POST",
    });
    if (!response.ok) {
        console.log("Logout failed");
    }
    const data = await response.json();
    console.log(data);
    
    console.log(data.message);
    Navigate("/login");
    login();
}
