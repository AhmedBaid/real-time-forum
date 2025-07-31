import { container, Header, Navigate, PostForm } from "../config.js";
import { login } from "./login.js";
import { timeFormat } from "../helpers/timeFormat.js";
import { HandleComments } from "./createComments.js";





async function fetchPosts() {
  try {
    let res = await fetch("/getPosts");
    if (!res.ok) {
      let errormsg = await res.text();
      throw new Error(errormsg);
    }
    return await res.json();
  } catch (error) {
    console.log(error);
  }
}

export async function home() {
  let header = document.createElement("header");
  let Postform = document.createElement("div");
  // posts
  let allPost = document.createElement("div");
  allPost.className = "allPost";
  let obj = await fetchPosts();
  for (const post of obj.data.Posts) {
    const postCombine = document.createElement("div");
    postCombine.className = "post-combine";

    const postId = post.id;
    const commentToggleId = `commentshow-${postId}`;
    const commentsSectionId = `comments-section-${postId}`;

    const commentsHTML =
      post.comments  !== null
        ? `
      <div class="commentaires">
        <h2 class="comment-title">Comments</h2>
        ${post.comments
          .map(
            (comment) => `
          <div class="comments">
            <img src="https://robohash.org/${
              comment.Username
            }.png?size=50x50" />
            <div class="comment-content">
              <p class="user"><strong>${comment.Username}</strong></p>
              <p class="comm">${comment.Comment}</p>
              <div class="comment-actions">
                <span class="time">${timeFormat(comment.time)}</span>
                <div class="comment-reactions">
                  <div class="reactionComment">
                    <span class="span-like ${
                      comment.UserReactionComment === 1 ? "active-like" : ""
                    }">
                      ${comment.TotalLikes}
                    </span>
                    <button class="comment-like-btn ${
                      comment.UserReactionComment === 1 ? "active-like" : ""
                    }" data-id="${comment.Id}" data-reaction="1" type="button">
                      <i class="fa-solid fa-thumbs-up"></i>
                    </button>
                  </div>
                  <div class="reactionComment">
                    <span class="span-dislike ${
                      comment.UserReactionComment === -1 ? "active-dislike" : ""
                    }">
                      ${comment.TotalDislikes}
                    </span>
                    <button class="comment-dislike-btn ${
                      comment.UserReactionComment === -1 ? "active-dislike" : ""
                    }" data-id="${comment.Id}" data-reaction="-1" type="button">
                      <i class="fa-solid fa-thumbs-down"></i>
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </div>
        `
          )
          .join("")}
      </div>
    `
        : `<h1 class="messageErr">No Commentaires 🤷‍♂️</h1>`;

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
          <form action="/reaction" method="post">
            <div class="reaction">
              <span class="span-like ${
                post.userReactionPosts === 1 ? "active-like" : ""
              }">${post.totalLikes}</span>
              <button name="reaction" value="1" class="like-btn ${
                post.userReactionPosts === 1 ? "active-like" : ""
              }" type="submit">
                <i class="fa-solid fa-thumbs-up"></i>
              </button>
            </div>

            <div class="reaction">
              <span class="span-dislike ${
                post.userReactionPosts === -1 ? "active-dislike" : ""
              }">${post.totalDislikes}</span>
              <button name="reaction" value="-1" class="dislike-btn ${
                post.userReactionPosts === -1 ? "active-dislike" : ""
              }" type="submit">
                <i class="fa-solid fa-thumbs-down"></i>
              </button>
            </div>

            <div class="reaction">
              <span>${post.totalComments}</span>
              <input type="checkbox" class="hidd" id="${commentToggleId}" />
              <label for="${commentToggleId}" class="comment-icon">
                <i class="fa-solid fa-comment"></i>
              </label>
            </div>

            <input type="hidden" name="postID" value="${postId}" />
          </form>
        </div>
      </div>

      <div class="second-part" id="${commentsSectionId}" style="display: none;">
        <!-- Form Add Comment -->
        <div class="comment">
          <form  method="post" id="formComment">
            <input type="hidden" name="postID" value="${postId}" />
            <img src="https://robohash.org/${
              obj.data.UserActive
            }.png?size=50x50" />
            <input type="text" name="comment" placeholder="Add Comment" required />
            <button type="submit" >Add</button>
          </form>
          <span class="error"></span>
        </div>

        ${commentsHTML}
      </div>
    </div>
  `;

    // 💡 Toggle logic (JS based)
    setTimeout(() => {
      const toggle = document.getElementById(commentToggleId);
      const section = document.getElementById(commentsSectionId);

      if (toggle && section) {
        toggle.addEventListener("change", () => {
          section.style.display = toggle.checked ? "flex" : "none";
        });
      }
    }, 0);

    allPost.appendChild(postCombine);
  }

  // end post container
  container.innerHTML = "";
  header.innerHTML = Header;
  Postform.innerHTML = PostForm;

  container.appendChild(header);
  container.appendChild(Postform);
  container.append(allPost);
  const logoutButton = header.querySelector(".logout");
  let createButton = header.querySelector(".create");

let formComment = document.getElementById("formComment")

formComment.addEventListener("submit", HandleComments)
  logoutButton.addEventListener("click", Logout);
  createButton.addEventListener("click", () => {
    const postForm = document.querySelector(".Post-form");
    if (postForm.style.display === "none" || postForm.style.display === "") {
      postForm.style.display = "block";
    } else {
      postForm.style.display = "none";
    }
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
