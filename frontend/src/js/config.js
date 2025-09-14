export const container = document.querySelector(".main");

export function Navigate(url) {
  history.pushState({}, "", url);
}

export const errorDiv = document.querySelector(".error");
export const successDiv = document.querySelector(".success");
export const errorMessage = document.getElementById("errmessage");
export const successMessage = document.getElementById("successmessage");


export const registerPage = `<div class="register">
            <form id="form" method="post" class="register-form">
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
            <form id="form" method="post" class="login-form">
                <h2>Login</h2>
                <input type="text" placeholder="Username"  id="username">
                <input type="password" placeholder="Password"  id="password">
                <button type="submit">Login</button>
                <h4>Don't have an account <a class="lienRegister">register</a></h4>
            </form>
        </div>`;

export const Header = `<h1>RT-<span class="logo">FO</span>RUM</h1>
    <span class="notifIcon"></span>
    <div class="buttons">
    <button class="create"><span class="span-create">Create Post</span><i class="fa-solid fa-plus"></i></button>
    <button class="logout"><span class="span-logout">Logout</span><i class="fa-solid fa-right-from-bracket"></i></button>
    </div>`;

const categories = [
  { id: 1, name: "Sport" },
  { id: 6, name: "Music" },
  { id: 3, name: "Movies" },
  { id: 5, name: "Gym" },
  { id: 2, name: "Technology" },
  { id: 7, name: "Science" },
  { id: 4, name: "Culture" },
  { id: 8, name: "Politics" },
];
export const PostForm = `
    <form method="post" class="post-form">
      <h2>Create Post</h2>
      <input type="text" id="title" class="title" placeholder="Title" />
      <textarea id="content" class="content" placeholder="Description"></textarea>
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
    </form>`;

export async function isLogged() {
  let response = await fetch("/isloged", {
    credentials: "include"
  });
  return response.ok;
}
