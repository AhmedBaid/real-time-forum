import { container, Header, Navigate, PostForm } from "../config.js";
import { login } from "./login.js";
import { timeFormat } from "../helpers/timeFormat.js";
import { HandleComments } from "./createComments.js";
import { HandleLikes } from "./HandleLikes.js";

const errorDiv = document.querySelector(".error");
const errorMessage = document.getElementById("message");

errorDiv.style.display = "none";
errorMessage.textContent = "";

async function fetchPosts() {
  try {
    let res = await fetch("/getPosts");
    if (!res.ok) {
      let errormsg = await res.text();
      throw new Error(errormsg);
    }
    return await res.json();
  } catch (error) {
    errorMessage.textContent = error;
    errorDiv.style.display = "flex";
  }
}

function renderCommentsStyled(section, comments) {
  const commentsContainer = section.querySelector(".comments-list");

  const commentsHTML =
    comments && comments.length > 0
      ? `
      <h2 class="comment-title">Comments</h2>
      <div class="commentaires">
        ${comments
          .map(
            (comment) => `
          <div class="comments">
            <img src="https://robohash.org/${
              comment.Username
            }.png?size=50x50" />
            <div class="comment-content">
              <div>
                <p class="user"><strong>${comment.Username}</strong></p>
                <p class="comm">${comment.Comment}</p>
              </div>
              <div class="comment-actions">
                <span class="time">${timeFormat(comment.time)}</span>
              </div>
            </div>
          </div>
        `
          )
          .join("")}
      </div>
    `
      : `<div class="commentaires"><h1 class="messageErr">No Commentaires 🤷‍♂️</h1></div>`;

  commentsContainer.innerHTML = commentsHTML;
}

async function fetchComments(postID) {
  try {
    const res = await fetch(`/getComments?id=${postID}`);
    if (!res.ok) {
      const errormsg = await res.text();
      throw new Error(errormsg);
    }

    return await res.json();
  } catch (error) {
    console.error("Error fetching comments:", error.message);
    return null;
  }
}

export async function home() {
  let header = document.createElement("header");
  let Postform = document.createElement("div");
  let allPost = document.createElement("div");
  allPost.className = "allPost";
  let obj = await fetchPosts();

  for (const post of obj.data.Posts) {
    const postCombine = document.createElement("div");
    postCombine.className = "post-combine";

    const postId = post.id;
    const commentToggleId = `commentshow-${postId}`;
    const commentsSectionId = `comments-section-${postId}`;

    postCombine.innerHTML = `
    <div class="post-card" id="post-${postId}">
      <div class="first-part">
        <div class="post-header">
          <div class="user-info">
            <img src="https://robohash.org/${
              post.username
            }.png?size=50x50" class="avatar" />
            <span class="username">${post.username}</span>
          </div>
          <span class="post-time">${timeFormat(post.time)}</span>
        </div>

        <h2 class="post-title">${post.title}</h2>
        <p class="post-description">${post.description}</p>

        <div class="post-tags">
          ${post.categories
            .map((cat) => `<span class="tag">${cat.name}</span>`)
            .join("")}
        </div>

        <div class="post-reactions">
          <form method="post" class="likesForm">
            <div class="reaction">
              <span class="span-like ${
                post.userReactionPosts === 1 ? "active-like" : ""
              }">${post.totalLikes}</span>
              <button name="reaction1" value="1" class="like-btn ${
                post.userReactionPosts === 1 ? "active-like" : ""
              }" type="submit">
                <i class="fa-solid fa-thumbs-up"></i>
              </button>
            </div>

            <div class="reaction">
              <span class="span-dislike ${
                post.userReactionPosts === -1 ? "active-dislike" : ""
              }">${post.totalDislikes}</span>
              <button name="reaction2" value="-1" class="dislike-btn ${
                post.userReactionPosts === -1 ? "active-dislike" : ""
              }" type="submit">
                <i class="fa-solid fa-thumbs-down"></i>
              </button>
            </div>

            <div class="reaction">
              <span class="totalComnts">${post.totalComments}</span>
              <input type="checkbox" class="hidd" value="${
                post.id
              }" id="${commentToggleId}" />
              <label for="${commentToggleId}" class="comment-icon">
                <i class="fa-solid fa-comment"></i>
              </label>
            </div>

            <input type="hidden" name="postID" value="${postId}" />
          </form>
        </div>
      </div>

      <div class="second-part" id="${commentsSectionId}" style="display: none;">
        <div class="comment">
          <form method="post" id="${postId}" class="formComment">
            <input type="hidden" name="postID" value="${postId}" />
            <img src="https://robohash.org/${
              obj.data.UserActive
            }.png?size=50x50" />
            <input type="text" name="comment" placeholder="Add Comment" required />
            <button type="submit">Add</button>
          </form>
        </div>
        <div class="comments-list"></div>
      </div>
    </div>
  `;

    setTimeout(() => {
      const toggle = document.getElementById(commentToggleId);
      const section = document.getElementById(commentsSectionId);

      if (toggle && section) {
        toggle.addEventListener("change", async () => {
          section.style.display = toggle.checked ? "flex" : "none";

          if (toggle.checked) {
            const data = await fetchComments(toggle.value);

            if (!data) return;

            renderCommentsStyled(section, data.data);
            const countSpan = section.closest(".post-card").querySelector(".totalComnts");
            if (countSpan) {
             countSpan.innerHTML=   data.data.length
            }
          }
        });
      }
    }, 0);

    allPost.appendChild(postCombine);
  }

  container.innerHTML = "";
  header.innerHTML = Header;
  Postform.innerHTML = PostForm;

  container.appendChild(header);
  container.appendChild(Postform);
  container.append(allPost);

  const logoutButton = header.querySelector(".logout");
  let createButton = header.querySelector(".create");

  document.querySelectorAll(".formComment").forEach((form) => {
    form.addEventListener("submit", HandleComments);
  });

  document.querySelectorAll(".likesForm").forEach((form) => {
    form.addEventListener("submit", HandleLikes);
  });

  logoutButton.addEventListener("click", Logout);

  createButton.addEventListener("click", () => {
    const postForm = document.querySelector(".Post-form");
    postForm.style.display =
      postForm.style.display === "none" || postForm.style.display === ""
        ? "block"
        : "none";
  });
}

async function Logout(e) {
  e.preventDefault();
  const response = await fetch("/logout", {
    method: "POST",
  });
  if (!response.ok) {
    console.log("Logout failed");
  }
  const data = await response.json();
  Navigate("/login");
  login();
}
