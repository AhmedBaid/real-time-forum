import { loadPage } from "./loadPage.js";
import { Navigate } from "./config.js";

async function isloged() {
    const response = await fetch("/isloged");
    const data = await response.json();
    if (response.ok) {
        if (location.pathname === "/login" || location.pathname === "/register") {
            console.log(2222);
            
            Navigate("/");
            loadPage(data);
        } else {
            loadPage(data);
            console.log(1);
            
        }
    } else {
        if (location.pathname === "/login" || location.pathname === "/register") {
            console.log("Already logged in, redirecting to home page.");
            
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