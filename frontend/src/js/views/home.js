import { container, Header, Navigate, PostForm } from "../config.js";
import { login } from "./login.js";
import { timeFormat } from "../helpers/timeFormat.js";
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
  console.log(obj.data.Posts[0].comments);

  for (const post of obj.data.Posts) {
   if  ( post.comments === null)   post.comments =  []
    const postCombine = document.createElement("div");
    postCombine.className = "post-combine";
    postCombine.innerHTML = `
    <div class="post-card" id="post-${post.id}">
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
              <input type="checkbox" class="hidd" id="commentshow-${post.id}" />
              <label for="commentshow-${
                post.id
              }" class="comment-icon"><i class="fa-solid fa-comment"></i></label>
              <style>
                #post-${post.id}:has(#commentshow-${
      post.id
    }:checked) .second-part {
                  display: flex;
                }
              </style>
            </div>

            <input type="hidden" name="postID" value="${post.id}" />
          </form>
        </div>
      </div>

      <div class="second-part" id="post-${post.id}">

             <div class="comment">
            <form action="/comment" method="post">
              <input type="hidden" name="postID" value="{{.Id}}" />
              <img src="https://robohash.org/{{$.UserActive}}.png?size=50x50" />
              <input type="text" name="comment" placeholder="Add Comment" required /><br />
              <button type="submit">Add</button>
            </form>
          </div>


             <div class="commentaires">
      <h2 class="comment-title">Comments</h2>
      
      ${post.comments.map(
        (comment) => `
        <div class="comments">
          <img src="https://robohash.org/${comment.Username}.png?size=50x50" />
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
      ).join("")}
    </div>

          
      </div>
    </div>
`;
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
