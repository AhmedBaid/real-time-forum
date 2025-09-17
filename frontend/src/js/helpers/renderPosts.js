// helpers/renderPost.js
import { timeFormat } from "./timeFormat.js";
import { fetchComments } from "./api.js";
import { renderCommentsStyled } from "./randerComments.js";
import { HandleComments } from "../views/createComments.js";
import { HandleLikes } from "../views/HandleLikes.js";

export function renderPost(post, userActiveId) {
  const postId = post.id;
  const commentToggleId = `commentshow-${postId}`;
  const commentsSectionId = `comments-section-${postId}`;

  const div = document.createElement("div");
  div.className = "post-card";
  div.id = `post-${postId}`;
  div.innerHTML = `
      <div class="first-part">
      <div class="post-header">
        <div class="user-info">
          <img src="https://robohash.org/${post.username}.png?size=50x50" class="avatar" />
          <span class="username">${post.username}</span>
        </div>
        <span class="post-time">${timeFormat(post.time)}</span>
      </div>
      <h2 class="post-title">${post.title}</h2>
      <p class="post-description">${post.description}</p>
      <div class="post-tags">
        ${post.categories.map(cat => `<span class="tag">${cat}</span>`).join("")}
      </div>
      <div class="post-reactions">
        <form method="post" class="likesForm">
          <div class="reaction">
            <span class="span-like ${post.userReactionPosts === 1 ? "active-like" : ""}">${post.totalLikes || 0}</span>
            <button name="reaction1" value="1" class="like-btn ${post.userReactionPosts === 1 ? "active-like" : ""}" type="submit">
              <i class="fa-solid fa-thumbs-up"></i>
            </button>
          </div>
          <div class="reaction">
            <span class="span-dislike ${post.userReactionPosts === -1 ? "active-dislike" : ""}">${post.totalDislikes || 0}</span>
            <button name="reaction2" value="-1" class="dislike-btn ${post.userReactionPosts === -1 ? "active-dislike" : ""}" type="submit">
              <i class="fa-solid fa-thumbs-down"></i>
            </button>
          </div>
          <div class="reaction">
            <span class="totalComnts">${post.totalComments || 0}</span>
            <input type="checkbox" class="hidd" value="${post.id}" id="${commentToggleId}" />
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
          <img src="https://robohash.org/${userActiveId}.png?size=50x50" />
          <input type="text" name="comment" placeholder="Add Comment" required />
          <button type="submit">Add</button>
        </form>
      </div>
      <div class="comments-list"></div>
    </div>
`;

  const toggle = div.querySelector(`#${commentToggleId}`);
  const section = div.querySelector(`#${commentsSectionId}`);
  if (toggle && section) {
    toggle.addEventListener("change", async () => {
      section.style.display = toggle.checked ? "flex" : "none";
      if (toggle.checked) {
        const data = await fetchComments(toggle.value);
        if (!data) return;
        renderCommentsStyled(section, data.data);
        const countSpan = section.closest(".post-card").querySelector(".totalComnts");
        if (countSpan) countSpan.innerHTML = data.datax ? data.data.length : 0 ;
      }
    });
  }

  document.querySelectorAll(".formComment").forEach((form) => {
    form.addEventListener("submit", HandleComments);
  });

  document.querySelectorAll(".likesForm").forEach((form) => {
    form.addEventListener("submit", HandleLikes);
  });

  return div;
}
