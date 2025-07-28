import { container, Header, Navigate, PostForm } from "../config.js";
import { login } from "./login.js";

export function home() {
    let header = document.createElement("header")
    let Postform = document.createElement("div")
    container.innerHTML = "";
    header.innerHTML = Header;
    Postform.innerHTML = PostForm;

    container.appendChild(header);
    container.appendChild(Postform);
    const logoutButton = header.querySelector(".logout");
    let createButton = header.querySelector(".create");
    logoutButton.addEventListener("click", Logout)
    createButton.addEventListener("click", () => {
        const postForm = document.querySelector(".Post-form");
        if (postForm.style.display === "none" || postForm.style.display === "") {
            postForm.style.display = "block";
        } else {
            postForm.style.display = "none";
        }
    });
}
async function Logout(e) {
    e.preventDefault();
    const response = await fetch("/logout", {
        method: "POST",
    });
    if (!response.ok) {
        console.log("Logout failed");
    }
    const data = await response.json();

    Navigate("/login");
    login();
}
