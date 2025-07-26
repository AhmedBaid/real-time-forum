import { loadPage } from "./loadPage.js";
import { Navigate } from "./config.js";

async function isloged() {
    try {
        const response = await fetch("/isloged");

        if (!response.ok) {
            console.log(1);

            Navigate("login");
            loadPage("login");
            return;
        }
        console.log(2);

        const data = await response.json();

        Navigate("/");
        loadPage("/", data);
    } catch (error) {
        console.error("Error checking login:", error.message);
        Navigate("login");
        loadPage("login");
    }
}

isloged()

window.onpopstate = () => {
    loadPage(location.pathname);
};