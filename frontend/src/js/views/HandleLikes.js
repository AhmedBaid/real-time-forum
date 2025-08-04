export async function HandleLikes(e) {
  e.preventDefault();
  let form = e.target.closest("form");

  const errorDiv = document.querySelector(".error");
  const errorMessage = document.getElementById("message");

  errorDiv.style.display = "none";
  errorMessage.textContent = "";

  let like = e.submitter.value;
  let postId = Number(form.querySelector("[name='postID']").value);

  try {
    
    let res = await fetch("/ReactionHandler", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ like, postId }),
    });
    
    if (!res.ok) {
      let er = await res.text();
      throw new Error(er);
    }
    
    const data = await res.json();


    console.log(data );
    console.log(form);

    spandislike = form.querySelector(".span-dislike");
    spanlike = form.querySelector(".span-like"); 



    




  } catch (error) {
    console.log(error);
    
    errorMessage.textContent = error.message;
    errorDiv.style.display = "flex";
  }
}
