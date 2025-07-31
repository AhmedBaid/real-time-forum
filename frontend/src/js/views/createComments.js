export async function HandleComments(e) {

    e.preventDefault();
    let errMsg = document.querySelector(".error")
    errMsg.innerHTML = "";
    let post_id =Number( document.getElementsByName("postID")[0].value);
    let Comment = document.getElementsByName("comment")[0].value;
    
    const response = await fetch("/createComment", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({ post_id , Comment })
    })

    const data = await response.json();
    if (!response.ok) {
        errMsg.innerHTML = data.message;
        console.log(data.message ,3);
        
        return;
    }
}