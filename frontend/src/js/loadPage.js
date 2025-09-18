import { container } from "./config.js";
import { register } from "./views/register.js";
import { login } from "./views/login.js";
import { home } from "./views/home.js";
import { notfound } from "./views/notfound.js";

export function loadPage() {
    if (location.pathname == "/login") {
        return login()
    } else if (location.pathname == "/") {
        return home()
    } else if (location.pathname == "/register") {
        return register()
    } else {
        return notfound()
    }
}
