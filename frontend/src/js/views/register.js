import { container, Navigate, registerPage } from "../config.js";
import { loadPage } from "../loadPage.js";
import { home } from "./home.js";

export function register() {
    container.innerHTML = ""
    container.innerHTML = registerPage;
    let form = document.querySelector("form");
    let lien = document.querySelector(".lienLogin");
    lien.addEventListener("click", (e) => {
        e.preventDefault();
        Navigate("/login");
        loadPage();
    });
    form.addEventListener("submit", HandleRegister);
}
async function HandleRegister(e) {
    e.preventDefault()
    let errMsg = document.querySelector(".error")
    errMsg.innerHTML = "";
    const fields = {
        username: document.getElementById("Username").value,
        email: document.getElementById("Email").value,
        firstName: document.getElementById("FirstName").value,
        lastName: document.getElementById("LastName").value,
        age: parseInt(document.getElementById("Age").value),
        gender: document.getElementById("Gender").value,
        password: document.getElementById("Password").value,
    };
    let response = await fetch("/register", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(fields)
    });
    const data = await response.json();

    if (!response.ok) {
        errMsg.innerHTML = data.message;
        return;
    }
    errMsg.innerHTML = "Registration successful!!";
    Navigate("/");
    home(data.data)
}