export function showToast(type, message) {
  const toast = document.querySelector(`.${type}`);
  console.log(toast);
  
  const span = toast.querySelector("span");
  const closeBtn = toast.querySelector("i");

  toast.style.display = "flex"
  span.textContent = message;

  setTimeout(() => {
    toast.style.display = "none"
    span.textContent = "";
  }, 2000);

  closeBtn.onclick = () => {
    toast.style.display = "none"
  };
}
