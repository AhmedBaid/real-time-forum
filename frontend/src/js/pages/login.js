import { container, loginPage, Navigate } from "../config.js";
import { loadPage } from "../loadPage.js";

export async function login() {
    container.innerHTML = ""
    container.innerHTML = loginPage;
    let lien = document.querySelector(".lienRegister");
    lien.addEventListener("click", (e) => {
        e.preventDefault();
        Navigate("/register");
        loadPage();
    });
    let form = document.querySelector("form");
    form.addEventListener("submit", HandleLogin);
}

async function HandleLogin(e) {
    e.preventDefault();
    let errMsg = document.querySelector(".error")
    errMsg.innerHTML = "";
    let username = document.getElementById("username").value;
    let password = document.getElementById("password").value;
    try {
        const response = await fetch("/login", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({ username, password })
        })

        const data = await response.json();
        Navigate("/");
        loadPage();
    } catch (err) {
        errMsg.innerHTML = err.message;
        return;
    }
}