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
    try {
        const response = await fetch("/login", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({ username, password })
        })
        
        const data = await response.json();
        
        if (!response.ok) {
            errMsg.innerHTML = data.message;
            return;
        }
        console.log(1);
        
        Navigate("/");
        loadPage("/", data);
    } catch (err) {
        console.log(2);
        
        errMsg.innerHTML = err.message;
        return;
    }
}