import { errorMessage, errorDiv, Navigate } from "../config.js";
import { showToast } from "./showToast.js";
import { login } from "../views/login.js";

export async function fetchPosts() {
  try {
    let res = await fetch("/getPosts");
    if (!res.ok) {
       if (res.status === 401) {
        showToast("error", "you are not authorized");
        Navigate("/login");
        login();
        return;
      }
      let errormsg = await res.text();
      throw new Error(errormsg);
    }
    return await res.json();
  } catch (error) {
    showToast("error", error.message);
  }
}

export async function fetchComments(postID) {
  try {
    const res = await fetch(`/getComments?id=${postID}`);
    if (!res.ok) {
      if (res.status === 401) {
        showToast("error", "you are not authorized");
        Navigate("/login");
        login();
        return;
      }
      let errormsg = await res.text();
      throw new Error(errormsg);
    }

    return await res.json();
  } catch (error) {
    showToast("error", error.message);
  }
}

export async function fetchUsers() {
  try {
    const res = await fetch(`/getUsers`);
    if (!res.ok) {
       if (res.status === 401) {
        showToast("error", "you are not authorized");
        Navigate("/login");
        login();
        return;
      }
      let errormsg = await res.text();
      throw new Error(errormsg);
    }

    return await res.json();
  } catch (error) {
    showToast("error", error.message);
  }
}
