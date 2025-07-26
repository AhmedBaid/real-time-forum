import { loadPage } from "./loadPage.js";
import { Navigate } from "./config.js";
async function isloged() {
    try {
        const response = await fetch("/isloged");
        console.log(response);
        
        if (!response.ok) {            
            Navigate("login")
            loadPage("login")
            return
        }
        const data = await response.json();
        console.log(data);
        
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