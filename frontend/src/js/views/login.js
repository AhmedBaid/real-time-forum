import { container, errorDiv, errorMessage, loginPage, Navigate, successDiv, successMessage } from "../config.js";
import { showToast } from "../helpers/showToast.js";
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
        showToast("error", data.message);
        return;
    }

    showToast("success", "Logged in successfully");
    Navigate("/");
    home(data.data);
}