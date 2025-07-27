import { container } from "./config.js";
import { register } from "./register.js";
import { login } from "./login.js";

export function loadPage() {    
    if (location.pathname == "/login") {
        login()
        return
    } else if (location.pathname == "/") {
        const name = data?.message || "Guest";
        container.innerHTML = `<h2>Welcome ${name}!</h2>`;
        return
    }else if (location.pathname == "/register") {
        register()
        return
    } else {
        container.innerHTML = `<h2>404 Not Found</h2>`;
        return
    }
}
