import { container, errorDiv, errorMessage, Header, Navigate, PostForm } from "../config.js";
import { login } from "./login.js";
import { timeFormat } from "../helpers/timeFormat.js";
import { HandleComments } from "./createComments.js";
import { HandleLikes } from "./HandleLikes.js";
import { fetchComments, fetchPosts, fetchUsers } from "../helpers/api.js";
import { renderCommentsStyled } from "../helpers/randerComments.js"
import { createPost } from "./createPost.js";
import { HandleMessages } from "./HandleMessages.js";
export let socket = new WebSocket("ws://localhost:8080/ws");
socket.onmessage = (e) => {
  let data = JSON.parse(e.data);

  switch (data.type) {
    case "online":
      setUserOnline(data.userId);
      break;

    case "offline":
      setUserOffline(data.userId);
      break;
    case "messages":
      appendMessage(data)
    case "online_list":
      data.users.forEach((id) => setUserOnline(id));
      break;
  }
};

function appendMessage(msg) {
  let chatBox = document.getElementById(`chat-${msg.sender}`);
  if (!chatBox) return;

  let messagesBox = chatBox.querySelector(".chat-messages");
  messagesBox.innerHTML += `
    <div class="msg ${msg.sender === currentUserId ? "right" : "left"}">
      <p>${msg.message}</p>
      <span class="time">${new Date(msg.time).toLocaleTimeString()}</span>
    </div>
  `;
  messagesBox.scrollTop = messagesBox.scrollHeight;
}
function setUserOnline(userId) {
  let el = document.querySelector(`.users[data-id="${userId}"] .online`);
  if (el) el.style.backgroundColor = "green";
}

function setUserOffline(userId) {
  let el = document.querySelector(`.users[data-id="${userId}"] .online`);
  if (el) el.style.backgroundColor = "red";
}

export async function home() {
  let header = document.createElement("header");
  let Postform = document.createElement("div");
  Postform.className = "Post-form";
  let parentContainer = document.createElement("div")
  parentContainer.className = "parentContainer"
  let allPost = document.createElement("div");
  let aside = document.createElement("div")
  aside.className = "aside2"
  allPost.className = "allPost";
  let users = await fetchUsers()

  for (const user of users.data) {
    const div = document.createElement("div")

    div.className = "users"
    div.dataset.username = user.username
    div.dataset.id = user.id
    div.innerHTML = `
  <img src="https://robohash.org/${user.username
      }.png?size=50x50" class="avatar" />
  <span class="username">${user.username}</span>
  <span class="online">.</span>
  `
    aside.appendChild(div)
  }

  let obj = await fetchPosts();

  for (const post of obj.data.Posts) {
    const postId = post.id;
    const commentToggleId = `commentshow-${postId}`;
    const commentsSectionId = `comments-section-${postId}`;

    allPost.innerHTML += `
    <div class="post-card" id="post-${postId}">
      <div class="first-part">
        <div class="post-header">
          <div class="user-info">
            <img src="https://robohash.org/${post.username
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
              <span class="span-like ${post.userReactionPosts === 1 ? "active-like" : ""
      }">${post.totalLikes}</span>
              <button name="reaction1" value="1" class="like-btn ${post.userReactionPosts === 1 ? "active-like" : ""
      }" type="submit">
                <i class="fa-solid fa-thumbs-up"></i>
              </button>
            </div>

            <div class="reaction">
              <span class="span-dislike ${post.userReactionPosts === -1 ? "active-dislike" : ""
      }">${post.totalDislikes}</span>
              <button name="reaction2" value="-1" class="dislike-btn ${post.userReactionPosts === -1 ? "active-dislike" : ""
      }" type="submit">
                <i class="fa-solid fa-thumbs-down"></i>
              </button>
            </div>

            <div class="reaction">
              <span class="totalComnts">${post.totalComments}</span>
              <input type="checkbox" class="hidd" value="${post.id
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
            <img src="https://robohash.org/${obj.data.UserActive
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
              countSpan.innerHTML = data.data.length
            }
          }
        });
      }
    }, 0);

  }

  container.innerHTML = "";
  header.innerHTML = Header;
  Postform.innerHTML = PostForm;

  //end  fetch posts 

  parentContainer.appendChild(allPost)
  parentContainer.appendChild(aside)

  container.appendChild(header);
  document.body.appendChild(Postform);
  container.append(parentContainer);


  const logoutButton = header.querySelector(".logout");
  let createButton = header.querySelector(".create");


  document.querySelectorAll(".users").forEach((user) => {
    user.addEventListener("click", HandleMessages)
  })

  document.querySelectorAll(".formComment").forEach((form) => {
    form.addEventListener("submit", HandleComments);
  });

  document.querySelectorAll(".likesForm").forEach((form) => {
    form.addEventListener("submit", HandleLikes);
  });

  logoutButton.addEventListener("click", Logout);

  createButton.addEventListener("click", () => {
    Navigate("/createpost");
    const postForm = document.querySelector(".Post-form");
    postForm.style.display =
      postForm.style.display === "none" || postForm.style.display === ""
        ? "block"
        : "none";
    container.style.opacity = postForm.style.display === "block" ? "0.2" : "1";
  });
  let form = document.querySelector(".post-form");
  form.addEventListener("submit", createPost);



}

async function Logout(e) {
  e.preventDefault();
  const response = await fetch("/logout", {
    method: "POST",
  });
  if (!response.ok) {
    errorDiv.style.display = "flex";
    errorMessage.textContent = "logout failed";
    return;
  }
  errorDiv.style.display = "flex";
  errorDiv.style.backgroundColor = "#04e17a";
  errorMessage.textContent = "logout successfully";
  Navigate("/login");
  login();
}
