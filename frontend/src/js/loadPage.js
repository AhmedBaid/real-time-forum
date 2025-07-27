import { container } from "./config.js";
import { register } from "./pages/register.js";
import { login } from "./pages/login.js";
import { home } from "./pages/home.js";

export function loadPage(data=null) {
    if (location.pathname == "/login") {
        return login()
    } else if (location.pathname == "/") {
        return home(data)
    } else if (location.pathname == "/register") {
        return register()
    } else {
        return container.innerHTML = `<h2>404 Not Found</h2>`;
    }
}
