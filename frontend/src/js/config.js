export const container = document.querySelector(".main");

export function Navigate(url) {
    history.pushState({}, "", url)
}

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
        </div>`

export const loginPage = `<div class="login">
            <div class="error"></div>
            <form id="form" method="post">
                <h2>Login</h2>
                <input type="text" placeholder="Username"  id="username">
                <input type="password" placeholder="Password"  id="password">
                <button type="submit">Login</button>
                <h4>Don't have an account <a class="lienRegister">register</a></h4>
            </form>
        </div>`

export const Header = `<h1>RT-<span>FO</span>RUM</h1>
    <div class="buttons">
    <button class="create">Create Post</button>
    <button class="logout">Logout</button>
    </div>`


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
          .join('')}
      </div>
      <button type="submit">Create Post</button>
    </form>
  </div>`;


export async function isLogged() {
    let response = await fetch("/isLogged");
    return response.ok;
}