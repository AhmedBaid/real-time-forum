import { timeFormat } from "../helpers/timeFormat.js";

export async function HandleComments(e) {
  e.preventDefault();
  let form = e.target;

  const errorDiv = document.querySelector(".error");
  const errorMessage = document.getElementById("message");

  errorDiv.style.display = "none";
  errorMessage.textContent = "";

  let post_id = Number(form.querySelector("[name='postID']").value);
  let Comment = form.querySelector("[name='comment']").value;

  if (post_id === 0 || Comment.trim() === "" || Comment.length < 3) {
    errorMessage.textContent = "Comment must be at least 3 characters.";
    errorDiv.style.display = "flex";
    return;
  }
  try {
    const response = await fetch("/createComment", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ post_id, Comment }),
    });

    if (!response.ok) {
      let err = await response.text();
      throw new Error(err);
    }

    const data = await response.json();

    const commentaires = form
      .closest(".second-part")
      .querySelector(".commentaires");

    let div = document.createElement("div");
    div.className = "comments";
   let j  =   document.querySelector(".messageErr")
   if (j) {
    j.innerHTML=""
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
let p  = div.querySelector(".comm")
p.textContent=data.data.Comment

    form.reset();

    commentaires.prepend(div);
    const totalSpan = form.closest(".post-card").querySelector(".totalComnts");
    if (totalSpan) {
      let current = Number(totalSpan.textContent);
      totalSpan.textContent = current + 1;
    }
  } catch (error) {
    errorMessage.textContent = error.message;
    errorDiv.style.display = "flex";
  }
}
