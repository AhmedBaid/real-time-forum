import { container } from "./config.js";
import { register } from "./views/register.js";
import { login } from "./views/login.js";
import { home } from "./views/home.js";

export function loadPage(data) {
    if (location.pathname == "/login") {
        return login()
    } else if (location.pathname == "/") {
        return home()
    } else if (location.pathname == "/register") {
        return register()
    } else {
        return container.innerHTML = `<h2>404 Not Found</h2>`;
    }
}
