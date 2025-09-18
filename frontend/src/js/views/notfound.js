import { container, Navigate, notFound } from "../config.js";
import { home } from "./home.js";

export function notfound() {
    container.innerHTML = notFound;
    let linkhome = document.querySelector(".linkHome")
    linkhome.addEventListener("click", (e) => {
        e.preventDefault()
        Navigate("/")
        home()
    })
}