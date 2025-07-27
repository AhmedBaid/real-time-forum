import { container, Navigate, registerPage } from "../config.js";

export function home() {
    container.innerHTML = ""
    container.innerHTML = "welcome to the home page";
}
