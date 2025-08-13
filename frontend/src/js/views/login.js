import { container, errorDiv, errorMessage, loginPage, Navigate } from "../config.js";
import { loadPage } from "../loadPage.js";
import { home } from "./home.js";

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
    let username = document.getElementById("username").value;
    let password = document.getElementById("password").value;

    const response = await fetch("/login", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({ username, password })
    })

    const data = await response.json();
    if (!response.ok) {
        errorDiv.style.display = "flex";
        errorMessage.textContent = data.message;
        return;
    }

    errorDiv.style.display = "flex";
    errorDiv.style.backgroundColor = "#04e17a";
    errorMessage.textContent = "User logged in successfully";
    Navigate("/");
    home(data.data);
}