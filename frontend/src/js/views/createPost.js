import { Navigate } from "../config.js";
import { showToast } from "../helpers/showToast.js";
import { home } from "./home.js";
import { login } from "./login.js";

export async function createPost(e) {
    e.preventDefault();
    const postForm = document.querySelector(".Post-form");
    let overlay = document.querySelector(".overlay");
    const title = document.querySelector(".title").value;
    const content = document.querySelector(".content").value;
    const categories = Array.from(
        document.querySelectorAll("input[name='tags']:checked")
    ).map(tag => Number(tag.value));
    console.log({ title, content, categories });

    const response = await fetch("/createpost", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({
            title: title,
            description: content,
            categories: categories,
        }),
    });
    const data = await response.json();
    if (!response.ok) {
        if (response.status === 401) {
            postForm.remove();
            overlay.remove();
            showToast("error", data.message);
            Navigate("/login");
            login();
            return;
        }
        showToast("error", data.message);
        return;
    }
    postForm.remove();
    overlay.remove();
    showToast("success", "Post created successfully");
    Navigate("/");
    home();
}
