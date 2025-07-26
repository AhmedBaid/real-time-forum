import { container, loginPage, Navigate } from "./config.js";

export async function login() {
    container.innerHTML = ""
    container.innerHTML = loginPage;
    let form = document.querySelector("form");
    form.addEventListener("submit", HandleLogin);
}

async function HandleLogin(e) {
    e.preventDefault();
    let errMsg = document.querySelector(".error")
    errMsg.innerHTML = "";
    let username = document.querySelector(".username").value;
    let password = document.querySelector(".password").value;
    if (!username || !password) {
        errMsg.innerHTML = "Username and Password are required";
        return;
    }
    if (username.length < 3 || password.length < 3) {
        errMsg.innerHTML = "Username and Password must be at least 3 characters long";
        return;
    }
    try {
        const response = await fetch("/login", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({ username, password })
        })
        if (!response.ok) {
            errMsg.innerHTML = "Invalid Username or Password";
            return;
        }
        const data = await response.json();
        if (data.message) {
            errMsg.innerHTML = data.message;
            return;
        }
        Navigate("/");
        loadPage("/", data);
    } catch (err) {
        errMsg.innerHTML = err.message;
        return;
    }
}