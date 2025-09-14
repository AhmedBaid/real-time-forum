import { Navigate } from "../config.js";
import { showToast } from "../helpers/showToast.js";
import { timeFormat } from "../helpers/timeFormat.js";
import { login } from "./login.js";

export async function HandleComments(e) {
  e.preventDefault();
  let form = e.target;

  let post_id = Number(form.querySelector("[name='postID']").value);
  let Comment = form.querySelector("[name='comment']").value;

  if (post_id === 0 || Comment.trim() === "" || Comment.length < 3) {
    showToast("error", "Comment must be at least 3 characters.");
    return;
  }
  try {
    const response = await fetch("/createComment", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ post_id, Comment }),
    });

    const data = await response.json();
    if (!response.ok) {

      if (response.status === 401) {
        Navigate("/login");
        login();
        showToast("error", data.message);
        return
      }
      showToast("error", data.message);zfùlvj
      return;
    }


    const commentaires = form
      .closest(".second-part")
      .querySelector(".commentaires");

    let div = document.createElement("div");
    div.className = "comments";
    let j = document.querySelector(".messageErr")
    if (j) {
      j.innerHTML = ""
    }
    div.innerHTML = `
      <img src="https://robohash.org/${data.data.Username}.png?size=50x50" />
      <div class="comment-content">
        <div>
          <p class="user"><strong>${data.data.Username}</strong></p>
          <p class="comm"></p>
        </div>
        <div class="comment-actions">
          <span class="time">${timeFormat(data.data.time)}</span>
        </div>
      </div>
    `;
    let p = div.querySelector(".comm")
    p.textContent = data.data.Comment

    form.reset();

    commentaires.prepend(div);
    const totalSpan = form.closest(".post-card").querySelector(".totalComnts");
    if (totalSpan) {
      let current = Number(totalSpan.textContent);
      totalSpan.textContent = current + 1;
    }
  } catch (error) {
    showToast("error", error.message);
  }
}
