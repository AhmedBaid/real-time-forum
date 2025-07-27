import { container } from "../config.js";

export function home(data) {
    container.innerHTML = "";
    const username = data.username;
    container.innerHTML = `Welcome ${username}!`;
}

