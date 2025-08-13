
import { timeFormat } from "./timeFormat.js";

export  function renderCommentsStyled(section, comments) {
  const commentsContainer = section.querySelector(".comments-list");

  const commentsHTML =
    comments && comments.length > 0
      ? `
      <h1 class="comment-title">Comments</h1>
      <hr/>
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