import { Navigate } from "../config.js";
import { showToast } from "../helpers/showToast.js";
import { timeFormat } from "../helpers/timeFormat.js";
import { home } from "./home.js";
import { login } from "./login.js";

export async function createPost(e) {
  e.preventDefault();
  const overlay = document.querySelector(".overlay");
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
      overlay.remove();
      showToast("error", data.message);
      Navigate("/login");
      login();
      return;
    }
    showToast("error", data.message);
    return;
  }
  overlay.remove();

  const allPost = document.querySelector(".allPost");

  const newPost = data.data;
  console.log(newPost, newPost.categories);

  let postId = newPost.id;
  let commentToggleId = `toggle-${postId}`;
  let commentsSectionId = `comments-${postId}`;
  const div = document.createElement("div");
  div.className = "post-card";
  div.id = `post-${postId}`;
  div.innerHTML = `
      <div class="first-part">
      <div class="post-header">
        <div class="user-info">
          <img src="https://robohash.org/${newPost.username}.png?size=50x50" class="avatar" />
          <span class="username">${newPost.username}</span>
        </div>
        <span class="post-time">${timeFormat(newPost.time)}</span>
      </div>
      <h2 class="post-title">${newPost.title}</h2>
      <p class="post-description">${newPost.description}</p>
      <div class="post-tags">
        ${newPost.categories.map((cat) => `<span class="tag">${cat}</span>`).join("")}
      </div>
      <div class="post-reactions">
        <form method="post" class="likesForm">
          <div class="reaction">
            <span class="span-like">0</span>
            <button name="reaction1" value="1" class="like-btn" type="submit">
              <i class="fa-solid fa-thumbs-up"></i>
            </button>
          </div>
          <div class="reaction">
            <span class="span-dislike">0</span>
            <button name="reaction2" value="-1" class="dislike-btn" type="submit">
              <i class="fa-solid fa-thumbs-down"></i>
            </button>
          </div>
          <div class="reaction">
            <span class="totalComnts">0</span>
            <input type="checkbox" class="hidd" value="${postId}" id="${commentToggleId}" />
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
          <img src="https://robohash.org/${newPost.username}.png?size=50x50" />
          <input type="text" name="comment" placeholder="Add Comment" required />
          <button type="submit">Add</button>
        </form>
      </div>
      <div class="comments-list"></div>
    </div>
`;

  allPost.prepend(div);

  Navigate("/");
  showToast("success", "Post created successfully");
}
