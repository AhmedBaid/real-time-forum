import { loadPage } from "./loadPage.js";
import { Navigate } from "./config.js";

async function isloged() {
    const response = await fetch("/isloged");
    const data = await response.json();
    if (response.ok) {
        Navigate("/");
        loadPage();
    } else {
        if (location.pathname === "/login" || location.pathname === "/register") {
            loadPage();
        } else {
            Navigate("/login");
            loadPage();
        }
    }
}

isloged()

window.onpopstate = () => {
    loadPage(location.pathname);
};