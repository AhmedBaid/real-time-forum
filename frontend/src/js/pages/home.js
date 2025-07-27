import { container } from "../config.js";

export function home(data = null) {
    container.innerHTML = "";
    const username = data?.username || "Guest";
    container.innerHTML = `Welcome ${username}!`;
}

