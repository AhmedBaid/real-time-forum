export const container = document.querySelector(".main");

export function Navigate(url) {
  history.pushState({}, "", url);
}

export const spanError = document.querySelector(".error");
export const spanMessage = document.getElementById("message");
export const registerPage = `<div class="register">
            <div class="error"></div>
            <form id="form" method="post">
                <h2>Register</h2>
                <input type="text" placeholder="Username" id="Username">
                <input type="text" placeholder="Email" id="Email">
                <input type="text" placeholder="FirstName" id="FirstName">
                <input type="text" placeholder="LastName" id="LastName">
                <input type="number" placeholder="Age" id="Age">
                <select id="Gender">
                    <option value="Male">Male</option>
                    <option value="Female">Female</option>
                </select>
                <input type="password" placeholder="Password" id="Password">
                <button type="submit">Register</button>
                <h4>Already have an account <a class="lienLogin">login</a></h4>
            </form>
        </div>`;

export const loginPage = `<div class="login">
            <div class="error"></div>
            <form id="form" method="post">
                <h2>Login</h2>
                <input type="text" placeholder="Username"  id="username">
                <input type="password" placeholder="Password"  id="password">
                <button type="submit">Login</button>
                <h4>Don't have an account <a class="lienRegister">register</a></h4>
            </form>
        </div>`;

export const Header = `<h1>RT-<span>FO</span>RUM</h1>
    <div class="buttons">
    <button class="create">Create Post</button>
    <button class="logout">Logout</button>
    </div>`;

const categories = [
  { id: 1, name: "General" },
  { id: 2, name: "Technology" },
  { id: 3, name: "Health" },
  { id: 4, name: "Sports" },
  { id: 5, name: "Entertainment" },
  { id: 6, name: "Education" },
];
export const PostForm = `
  <div class="Post-form">
    <form method="post">
      <h2>Create Post</h2>
      <input type="text" id="title" name="title" placeholder="Title" />
      <textarea id="content" name="content" placeholder="Description"></textarea>
      <div class="categories">
        ${categories
          .map(
            (category) => `
              <input type="checkbox" id="tag${category.id}" class="tag-checkk" name="tags" value="${category.id}" hidden />
              <label for="tag${category.id}" class="tagCC">${category.name}</label>
            `
          )
          .join("")}
      </div>
      <button type="submit">Create Post</button>
    </form>
  </div>`;

export async function isLogged() {
  let response = await fetch("/isLogged");
  return response.ok;
}

export const postContainer = `
 <div class="post-combine">
      <div class="post-card" id="post-{{.Id}}">
        <div class="first-part">
          <div class="post-header">
            <div class="user-info">
              <img src="https://robohash.org/{{.Username}}.png?size=50x50" class="avatar" />
              <span class="username">{{.Username}}</span>
            </div>
            <span class="post-time">{{.TimeFormatted}}</span>
          </div>

          <h2 class="post-title">{{.Title}}</h2>


          {{if .ImagePath}}
          <img src="{{.ImagePath}}" class="image" />
          {{end}}

          <p class="post-description">{{.Description}}</p>
          <div class="post-tags">
            {{range .Categories}}
            <span class="tag">{{.Name}}</span>
            {{end}}
          </div>
          <div class="post-reactions">
            <form action="/reaction" method="post">
              <div class="reaction">
                <span class="span-like  {{if eq .UserReactionPosts 1}}active-like{{end}}">{{.TotalLikes}}</span>
                <button name="reaction" value="1" class="like-btn  {{if eq .UserReactionPosts 1}}active-like{{end}}"
                  type="submit">
                  <i class="fa-solid fa-thumbs-up"></i>
                </button>
              </div>
              <div class="reaction">
                <span
                  class="span-dislike {{if eq .UserReactionPosts -1}}active-dislike{{end}}">{{.TotalDislikes}}</span>
                <button name="reaction" value="-1"
                  class="dislike-btn  {{if eq .UserReactionPosts -1}}active-dislike{{end}}" type="submit">
                  <i class="fa-solid fa-thumbs-down"></i>
                </button>
              </div>
              <div class="reaction">
                <span>{{.TotalComments}}</span>
                <input type="checkbox" class="hidd" id="commentshow-{{.Id}}" />
                <label for="commentshow-{{.Id}}" class="comment-icon"><i class="fa-solid fa-comment"></i></label>
                <style>
                  #post-{{.Id}}:has(#commentshow-{{.Id}}:checked) .second-part {
                    display: flex;
                  }
                 </style>
              </div>
              {{if eq $.UserActive .Username}}
              <button type="submit" formaction="/deletePost" formmethod="post" name="postID" value="{{.Id}}" class="delete-btn" title="Delete post" onclick="return confirm('Are you sure you want to delete this post?');">
                <i class="fa-solid fa-trash-can"></i>
              </button>
              {{end}}
              <input type="hidden" name="postID" value="{{.Id}}" />
            </form>
        
          </div>
        </div>
        <div class="second-part" id="post-{{.Id}}">
          {{if $.Session}}
          <div class="comment">
            <form action="/comment" method="post">
              <input type="hidden" name="postID" value="{{.Id}}" />
              <img src="https://robohash.org/{{$.UserActive}}.png?size=50x50" />
              <input type="text" name="comment" placeholder="Add Comment" required /><br />
              <button type="submit">Add</button>
            </form>
          </div>
          {{end}}
          {{if gt (len .Comments) 0}}
          <div class="commentaires">
            <h2 class="comment-title">Comments</h2>
            {{range .Comments}}
            <div class="comments">
              <img src="https://robohash.org/{{.Username}}.png?size=50x50" />
              <div class="comment-content">
                <p class="user"><strong>{{.Username}}</strong></p>
                <p class="comm">{{.Comment}}</p>
                <div class="comment-actions">
                  <span class="time">{{.TimeFormattedComment}}</span>
                  <div class="comment-reactions">
                    <form action="/CommentsLike " method="post">
                      <div class="reactionComment">
                        <span
                          class="span-like  {{if eq .UserReactionComment 1}}active-like{{end}}">{{.TotalLikes}}</span>
                        <button name="reaction" value="1"
                          class="comment-like-btn {{if eq .UserReactionComment 1}}active-like{{end}}" type="submit">
                          <i class="fa-solid fa-thumbs-up"></i>
                        </button>
                      </div>
                      <div class="reactionComment">
                        <span
                          class="span-dislike {{if eq .UserReactionComment -1}}active-dislike{{end}}">{{.TotalDislikes}}</span>
                        <button name="reaction" value="-1"
                          class="comment-dislike-btn {{if eq .UserReactionComment -1}}active-dislike{{end}}"
                          type="submit">
                          <i class="fa-solid fa-thumbs-down"></i>
                        </button>
                      </div>
                      <input type="hidden" name="userId" value="{{$userId}}" />
                      <input type="hidden" name="commentID" value="{{.Id}}" />
                    </form>
                  </div>
                </div>
              </div>
            </div>
            {{end}}
          </div>
          {{else}}
          <h1 class="messageErr">No Commentaires 🤷‍♂️</h1>
          {{end}}
        </div>
      </div>
    </div>`;
