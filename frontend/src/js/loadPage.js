import { container } from "./config.js";
import { register } from "./register.js";
import { login } from "./login.js";

export function loadPage(path, data = null) {    
    if (path == "login") {
        console.log("Loading login page");
        
        login()
        return
    } else if (path == "/") {
        console.log(555);
        
        const name = data?.message || "Guest";
        container.innerHTML = `<h2>Welcome ${name}!</h2>`;
        return
    } if (path == "register") {
        register()
        return
    } else {
        container.innerHTML = `<h2>404 Not Found</h2>`;
        return
    }
}
