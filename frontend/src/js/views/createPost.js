import { Navigate } from "../config.js";
import { renderPost } from "../helpers/renderPosts.js";
import { showToast } from "../helpers/showToast.js";
import { HandleComments } from "./createComments.js";
import { HandleLikes } from "./HandleLikes.js";
import { login } from "./login.js";

export async function createPost(e) {
  e.preventDefault();
  const overlay = document.querySelector(".overlay");
  const title = document.querySelector(".title").value;
  const content = document.querySelector(".content").value;
  const categories = Array.from(
    document.querySelectorAll("input[name='tags']:checked")
  ).map(tag => Number(tag.value));
  console.log(categories);

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
      overlay.remove();
      showToast("error", data.message);
      Navigate("/login")
      login()
      return;
    }
    showToast("error", data.message);
    return;
  }
  overlay.remove();
  console.log("efefef");

  showToast("success", "Post created successfully");
  Navigate("/");
  const allPost = document.querySelector(".allPost");
  const newPost = data.data;
  const userActiveId = newPost.username;
  const postDiv = renderPost(newPost, userActiveId);
  allPost.prepend(postDiv);

  const newPostForm = postDiv.querySelector(".formComment");
  if (newPostForm) {
    newPostForm.addEventListener("submit", HandleComments);
  }

  const newLikesForm = postDiv.querySelector(".likesForm");
  if (newLikesForm) {
    newLikesForm.addEventListener("submit", HandleLikes);
  }
}
