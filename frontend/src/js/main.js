import { loadPage } from "./loadPage";
async function isloged() {
    try {
        const response = await fetch("/isloged");
        if (!response.ok) {
            Navigate("/login")
        }
        const json = await response.json();
        loadPage(json)
    } catch (error) {
        console.error(error.message);
    }
}

function Navigate(url) {
    history.pushState(null, null, url)
}
isloged()