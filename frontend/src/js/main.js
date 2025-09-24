import { loadPage } from "./loadPage.js";
import { Navigate } from "./config.js";

async function isloged() {
    const response = await fetch("/isloged");
    if (response.ok) {
        if (location.pathname === "/login" || location.pathname === "/register") {
            Navigate("/");
            return loadPage();
        }
        return loadPage();
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
    Navigate("/")
    return
 
};