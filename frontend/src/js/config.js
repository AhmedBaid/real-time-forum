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
                <input type="text" placeholder="Username" name="Username" class="username">
                <input type="text" placeholder="Password" name="Password" class="password">
                <button type="submit">Login</button>
                <h4>Don't have an account <a class="lienRegister">register</a></h4>
            </form>
        </div>`
