import { container } from "../config.js";

export function home(data) {    
    container.innerHTML = ""
    container.innerHTML = `welcome ${data.username}!`;
}
