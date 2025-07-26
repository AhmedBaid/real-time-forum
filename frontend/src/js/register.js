import { container ,registerPage} from "./config.js";

export function register() {
    container.innerHTML=""
    container.innerHTML = registerPage;
}