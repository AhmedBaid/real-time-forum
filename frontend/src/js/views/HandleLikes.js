export async function HandleLikes(e) {
  e.preventDefault();

  const form = e.target.closest("form");
  const errorDiv = document.querySelector(".error");
  const errorMessage = document.getElementById("message");

  errorDiv.style.display = "none";
  errorMessage.textContent = "";

  const like =
    e.submitter?.value ||
    form.querySelector("button[type='submit']:focus")?.value;
  const postId = Number(form.querySelector("[name='postID']").value);

  try {
    const res = await fetch("/ReactionHandler", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ like, postId }),
    });

    if (!res.ok) {
      const er = await res.text();
      throw new Error(er);
    }

    const data = await res.json();

    let spanlike = form.querySelector(".span-like");
    let spandislike = form.querySelector(".span-dislike");

    console.log("spanlike", spanlike);
    console.log("spandislike", spandislike);
console.log(data);

    spanlike.textContent = data.TotalLike;
    spandislike.textContent = data.TotalDislikes;
  } catch (error) {
    console.log(error);
    errorMessage.textContent = error.message;
    errorDiv.style.display = "flex";
  }
}
