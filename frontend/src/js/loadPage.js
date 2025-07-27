import { container } from "./config.js";
import { register } from "./pages/register.js";
import { login } from "./pages/login.js";
import { home } from "./pages/home.js";

export function loadPage(data) {    
    if (location.pathname == "/login") {
        login()
        return
    } else if (location.pathname == "/") {
        home(data)
        return
    } else if (location.pathname == "/register") {
        register()
        return
    } else {
        container.innerHTML = `<h2>404 Not Found</h2>`;
        return
    }
}
