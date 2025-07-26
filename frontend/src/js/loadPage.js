import { container, registerPage, loginPage } from "./config.js";
export function loadPage(path, data = null) {
    console.log(path, "efef");

    if (path == "login") {
        container.innerHTML = loginPage;
        return
    } else if (path == "/") {
        const name = data?.message || "Guest";
        container.innerHTML = `<h2>Welcome ${name}!</h2>`;
        return
    } if (path == "register") {
        container.innerHTML = registerPage;
        return
    } else {
        container.innerHTML = `<h2>404 Not Found</h2>`;
        return
    }
}
