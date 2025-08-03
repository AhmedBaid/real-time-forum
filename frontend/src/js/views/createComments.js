import { timeFormat } from "../helpers/timeFormat.js";

export async function HandleComments(e) {
  e.preventDefault();
  let errMsg = document.querySelector(".error");
  errMsg.innerHTML = "";
  let post_id = Number(document.getElementsByName("postID")[0].value);
  let Comment = document.getElementsByName("comment")[0].value;

  const response = await fetch("/createComment", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ post_id, Comment }),
  });

  const data = await response.json();
  if (!response.ok) {
    errMsg.innerHTML = data.message;
    console.log(data.message, 3);

    return;
  }

  let commentaires = document.querySelector(".commentaires");

  let div = document.createElement("div");
  div.className = "comments";

  div.innerHTML = `
 <img src="https://robohash.org/redaanniz.png?size=50x50" />
            <div class="comment-content">
              <p class="user"><strong>redaanniz</strong></p>
              <p class="comm">${Comment}</p>
              <div class="comment-actions">
                <span class="time">${timeFormat(new Date())}</span>
              </div>
 </div>
`;
  commentaires.prepend(div);
}
