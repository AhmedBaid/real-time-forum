import { container, registerPage } from "./config.js";

export function register() {
    container.innerHTML = ""
    container.innerHTML = registerPage;
    let form = document.querySelector("form");
    form.addEventListener("submit", HandleRegister);
}
async function HandleRegister(e) {
    e.preventDefault()
    let errMsg = document.querySelector(".error")
    errMsg.innerHTML = "";
    const data = {
        username: document.getElementById("Username").value,
        email: document.getElementById("Email").value,
        firstName: document.getElementById("FirstName").value,
        lastName: document.getElementById("LastName").value,
        age: parseInt(document.getElementById("Age").value),
        gender: document.getElementById("Gender").value,
        password: document.getElementById("Password").value,
    };
    if (!data.username || !data.email || !data.firstName || !data.lastName || !data.age) {
        errMsg.innerHTML = "All fields are required";
        return;
    }
    try {
        let response = await fetch("/register", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify(data)
        });
        if (!response.ok) {
            errMsg.innerHTML = "Registration failed. Please try again.";
            return;
        }
        const result = await response.json();
        if (result.message) {
            errMsg.innerHTML = result.message;
            return;
        }
        errMsg.innerHTML = "Registration successful! You can now log in.";
    } catch (error) {
        errMsg.innerHTML = "An error occurred: " + error.message;
        return
    }
}