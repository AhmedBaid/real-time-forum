import {errorMessage ,errorDiv  }  from "../config.js" 


export  async function fetchPosts() {
  try {
    let res = await fetch("/getPosts");
    if (!res.ok) {
      let errormsg = await res.text();
      throw new Error(errormsg);
    }
    return await res.json();
  } catch (error) {
    errorMessage.textContent = error;
    errorDiv.style.display = "flex";
  }
}



export async function fetchComments(postID) {
  try {
    const res = await fetch(`/getComments?id=${postID}`);
    if (!res.ok) {
      const errormsg = await res.text();
      throw new Error(errormsg);
    }

    return await res.json();
  } catch (error) {
    console.error("Error fetching comments:", error.message);
    return null;
  }
}
