import {
  container,
  Header,
  isLogged,
  Navigate,
  PostForm,
} from "../config.js";
import { login } from "./login.js";
import { timeFormat } from "../helpers/timeFormat.js";
import { HandleComments } from "./createComments.js";
import { HandleLikes } from "./HandleLikes.js";
import { fetchComments, fetchPosts } from "../helpers/api.js";
import { renderCommentsStyled } from "../helpers/randerComments.js";
import { createPost } from "./createPost.js";
import { HandleMessages } from "./HandleMessages.js";
import { sortUsers } from "../helpers/sortUsers.js";
import { showToast } from "../helpers/showToast.js";
export let currentUserId = null;
export let Currentusername = null;
export let offset = { nbr: 0 };
let id = null;


// get the current user
async function fetchCurrentUserId() {
  try {
    const res = await fetch("/api/current-user");
    if (!res.ok) throw new Error("Failed to fetch user ID");
    const data = await res.json();
    currentUserId = data.userId;
    Currentusername = data.username;
  } catch (error) {
    console.error("Error fetching current user ID:", error);
  }
}

export let socket = null;
// websocket andler
async function connectWebSocket() {
  console.log("hadii tanya");

  socket = new WebSocket("ws://localhost:8080/ws");

  socket.onopen = () => {
    console.log("WebSocket connected");
    fetchCurrentUserId();
  };

  socket.onmessage = async (e) => {
    let data = JSON.parse(e.data);

    const logged = await isLogged()
    if (!logged) {
      Navigate("/login")
      login()
      socket.close()
      return
    }
    switch (data.type) {
      case "online":
        setTimeout(async () => {
          const aside = document.querySelector(".aside2")
          await sortUsers(aside)
          let el = document.querySelector(`.users`);
          if (el) {
            setUserOnline(data.userId);
          }
        }, 0);
        break;
      case "offline":
        setTimeout(() => {
          let el = document.querySelector(`.users`);
          if (el) {
            setUserOffline(data.userId);
          }
        }, 0);

        break;
      case "message":
        setTimeout(async () => {
          const aside = document.querySelector(".aside2")
          await sortUsers(aside)
          appendMessage(data);
        }, 0);
        break;
      case "notification":
        setTimeout(async () => {
          let aside = document.querySelector(".aside2");
          await sortUsers(aside);

          let chatbox = document.querySelector(
            `.chat-box[data-id-u="${data.from}"]`
          );
          if (chatbox) {
            return;
          }

          let notif = document.querySelector(".notifIcon");
          notif.innerHTML = `<i class="fa-solid fa-bell bell-icon" id="bellIcon"></i>`;

          const user = document.querySelector(
            `.users[data-id="${data.from}"] .text-wrapper .notification`
          );

          if (user) {
            user.innerHTML = "new Message";
          }

          setTimeout(() => {
            notif.innerHTML = "";
          }, 2000);
        }, 0);

        break;

      case "online_list":
        setTimeout(() => {

          console.log(data);
          let el = document.querySelector(`.users`);
          if (el) {
            data.users.forEach((id) => setUserOnline(Number(id)));
          }

        }, 0);

        break;
      case "typing":
        setTimeout(() => {
          const typing = document.querySelector(
            `.users[data-id="${data.senderId}"] .text-wrapper .typing  `
          );
          typing.style.display = "block";

          let chatBox = document.querySelector(
            `#chat-${data.senderUsername} .chatTyping`
          );
          if (chatBox) {
            chatBox.style.display = "block";
            const str = chatBox.querySelector("strong");
            str.textContent = data.senderUsername + " typing";
          }
        }, 0);
        clearTimeout(id);
        id = setTimeout(() => {
          const typing = document.querySelector(
            `.users[data-id="${data.senderId}"] .text-wrapper .typing  `
          );
          typing.style.display = "none";

          let chatBox = document.querySelector(
            `#chat-${data.senderUsername} .chatTyping`
          );
          if (chatBox) {
            chatBox.style.display = "none";
            const str = chatBox.querySelector("strong");
            str.textContent = "";
          }
        }, 500);
        break;
      case "stopTyping":
        const typing = document.querySelector(
          `.users[data-id="${data.senderId}"] .text-wrapper .typing  `
        );
        typing.style.display = "none";

        let chatBox = document.querySelector(
          `#chat-${data.senderUsername} .chatTyping`
        );
        if (chatBox) {
          chatBox.style.display = "none";
          const str = chatBox.querySelector("strong");
          str.textContent = "";
        }
        break;
    }
  };

  socket.onerror = (err) => console.error("WebSocket error:", err);
  socket.onclose = () => {
    console.log("WebSocket disconnected");
    setTimeout(connectWebSocket, 5000);
  };
}


// messages realtime
function appendMessage(msg) {
  let chatBoxSender = document.getElementById(`chat-${msg.senderUsername}`);
  let chatBoxreciever = document.getElementById(`chat-${msg.receiverUsername}`);




  let chatBox = chatBoxSender ? chatBoxSender : chatBoxreciever;
  if (!chatBox) {
    console.log(msg.senderUsername);
    console.log(msg.receiverUsername);



    return;
  }

  let messagesBox = chatBox.querySelector(".chat-messages");
  let div = document.createElement("div");
  div.className = `msg ${msg.sender === currentUserId ? "right" : "left"}`;

  let p = document.createElement("p");
  p.innerHTML = msg.message;

  let span = document.createElement("span");
  span.className = "time";
  span.innerHTML =
    msg.senderUsername + " - " + new Date(msg.time).toLocaleString();

  div.appendChild(p);
  div.appendChild(span);
  messagesBox.appendChild(div);

  messagesBox.scrollTop = messagesBox.scrollHeight;
  offset.nbr += 1;
}

// online handler
function setUserOnline(userId) {
  let el = document.querySelector(`.users[data-id="${userId}"] .online`);
  if (el) el.style.backgroundColor = "green";
}
//offline handler
function setUserOffline(userId) {
  let el = document.querySelector(`.users[data-id="${userId}"] .online`);
  if (el) el.style.backgroundColor = "grey";
}

export async function home() {
  let logged = await isLogged()
  if (!logged) {
    Navigate("/login")
    login()
    socket.close()
    return
  }
  let header = document.createElement("header");
  let parentContainer = document.createElement("div");
  parentContainer.className = "parentContainer";
  let allPost = document.createElement("div");
  let aside = document.createElement("div");
  aside.className = "aside2";
  allPost.className = "allPost";

  let obj = await fetchPosts();

  if (obj.data.Posts && obj.data.Posts.length > 0) {
    allPost.innerHTML = ""; // Clear previous content, including "No posts yet"
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
              const countSpan = section
                .closest(".post-card")
                .querySelector(".totalComnts");
              if (countSpan) {
                countSpan.innerHTML = !data.data ? 0 : data.data.length;;
              }
            }
          });
        }
      }, 0);
    }

  } else {
    allPost.innerHTML = "<h2>No posts yet</h2>";
  }

  container.innerHTML = "";
  header.innerHTML = Header;
  parentContainer.appendChild(allPost);
  parentContainer.appendChild(aside);
  container.appendChild(header);
  container.append(parentContainer);
  await sortUsers(aside);
  const logoutButton = header.querySelector(".logout");
  let createButton = header.querySelector(".create");

  document.querySelectorAll(".users").forEach((user) => {
    user.addEventListener("click", HandleMessages);
  });

  document.querySelectorAll(".formComment").forEach((form) => {
    form.addEventListener("submit", HandleComments);
  });

  document.querySelectorAll(".likesForm").forEach((form) => {
    form.addEventListener("submit", HandleLikes);
  });

  logoutButton.addEventListener("click", Logout);

  createButton.addEventListener("click", () => {
    const postForm = document.querySelector(".Post-form");

    if (!postForm) {
      Navigate("/createpost");
      const overlay = document.createElement("div");
      overlay.className = "overlay";

      const injecthtml = document.createElement("div");
      injecthtml.className = "Post-form";
      injecthtml.innerHTML = PostForm;

      overlay.appendChild(injecthtml);
      document.body.appendChild(overlay);



      const form = injecthtml.querySelector("form");
      form.addEventListener("submit", createPost);

      overlay.addEventListener("click", async (e) => {
        if (e.target === overlay) {
          overlay.remove();
          Navigate("/");
        }
      });
    }
  });

  await connectWebSocket();
}


export async function Logout(e) {
  e.preventDefault();
  Navigate("/logout");
  const response = await fetch("/logout", {
    method: "POST",
  });
  if (!response.ok) {


    showToast("error", "Failed to logout ");
    Navigate("/login");
    login();
    return;
  }
  showToast("success", "Logged out successfully");
  socket.close();
  Navigate("/login");
  login();
}


