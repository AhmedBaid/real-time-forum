export const container = document.querySelector(".main");

export function Navigate(url) {
    console.log(3);
    
    history.pushState({}, "", url)
}

export const registerPage = `<div class="register">
            <div class="error"></div>
            <form id="form" method="post">
                <h2>Register</h2>
                <input type="text" placeholder="Username" name="Username">
                <input type="text" placeholder="Email" name="Email">
                <input type="text" placeholder="FirstName" name="FirstName">
                <input type="text" placeholder="LastName" name="LastName">
                <input type="text" placeholder="Age" name="Age">
                <select name="Gender">
                    <option value="Male">Male</option>
                    <option value="Female">Female</option>
                </select>
                <input type="text" placeholder="Password" name="Password">
                <button type="submit">Register</button>
                <h4>Already have an account <a href="/login">login</a></h4>
            </form>
        </div>`


export const loginPage = `<div class="login">
            <div class="error"></div>
            <form id="form" method="post">
                <h2>Login</h2>
                <input type="text" placeholder="Username" name="Username" class="username">
                <input type="text" placeholder="Password" name="Password" class="password">
                <button type="submit">Login</button>
                <h4>Don't have an account <a href="/register">register</a></h4>
            </form>
        </div>`
