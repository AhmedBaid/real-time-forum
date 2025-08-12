import { container, Navigate, PostForm } from "../config.js";
import { home } from "./home.js";

export async function createPost(e) {
    e.preventDefault();
    console.log("createPost function called");
    const errorDiv = document.querySelector(".error");
    const errorMessage = document.getElementById("message");
    errorDiv.style.display = "none";
    errorMessage.textContent = "";
    const title = document.querySelector(".title").value;
    const content = document.querySelector(".content").value;
    const categories = Array.from(
        document.querySelectorAll("input[name='tags']:checked")
    ).map(tag => tag.value);
    // console.log(title, content, categories);
    const response = await fetch("/createpost", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({
            title : title,
            description: content,
            categories : categories,
        }),
    });
    const data = await response.json();
    if (!response.ok) {
        errorDiv.style.display = "flex";
        errorMessage.textContent =data.message;
        return;
    }
    const postForm = document.querySelector(".Post-form");    
    postForm.style.display = "none";
    container.style.opacity = "1";
    errorDiv.style.display = "flex";
    errorDiv.style.backgroundColor = "#04e17a";
    errorMessage.textContent = "Post created successfully";
    Navigate("/");
    home();
}