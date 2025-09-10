import { errorDiv, errorMessage } from "../config.js";
import { showToast } from "../helpers/showToast.js";

export async function HandleLikes(e) {
  e.preventDefault();

  const form = e.target.closest("form");

  errorDiv.style.display = "none";
  errorMessage.textContent = "";

  const like =
    e.submitter?.value ||
    form.querySelector("button[type='submit']:focus")?.value;
  const postId = Number(form.querySelector("[name='postID']").value);

  try {
    const res = await fetch("/ReactionHandler", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ like, postId }),
    });
    if (!res.ok) {
      const er = await res.text();
      throw new Error(er);
    }
    const data = await res.json();

    let spanlike = form.querySelector(".span-like");
    let spandislike = form.querySelector(".span-dislike");
    spanlike.textContent = data.data.TotalLike;
    spandislike.textContent = data.data.TotalDislikes;
    let likebtn = form.querySelector("[name='reaction1']");
    let deslikebtn = form.querySelector("[name='reaction2']");


    if (data.data.userReactionPosts === 1) {
      likebtn.classList.add("active-like");
      deslikebtn.classList.remove("active-dislike");
    } else {
      likebtn.classList.remove("active-like");
    }

    if (data.data.userReactionPosts === -1) {

      deslikebtn.classList.add("active-dislike");
      likebtn.classList.remove("active-like");
    } else {
      deslikebtn.classList.remove("active-dislike");
    }

  } catch (error) {
    showToast("error", error.message);
  }
}
