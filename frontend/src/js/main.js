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
};