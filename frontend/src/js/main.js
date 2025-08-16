import { loadPage } from "./loadPage.js";
import { Navigate } from "./config.js";

async function isloged() {
    const response = await fetch("/isloged");
    const data = await response.json();
    if (response.ok) {
        if (location.pathname === "/login" || location.pathname === "/register") {
            Navigate("/");
            return loadPage(data);
        }
        return loadPage(data);
    } else {
        if (location.pathname === "/login" || location.pathname === "/register") {
            return loadPage();
        }
        Navigate("/login");
        return loadPage();
    }
}
isloged()

window.onpopstate = () => {
    loadPage(location.pathname);
    const content = document.getElementById("content")
content.innerHTML = `
<div class="login">
    <h1>Zone01</h1>
        <form  id="form">
            <input type="text"  placeholder="username/email" id="username">
            <input type="password"  placeholder="password" id="password">
            <span id="error"></span>
            <button type="submit">login</button>
        </form>
</div>
`
};