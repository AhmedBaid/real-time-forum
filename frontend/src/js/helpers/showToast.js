import { errorDiv, errorMessage, successDiv, successMessage } from "../config.js";

export function showToast(type, message) {
  const toast = document.querySelector(`.${type}`);
  const span = toast.querySelector(".toast");
  const closeBtn = toast.querySelector("i");
  errorDiv.style.display = "none";
  errorMessage.textContent = "";
  successDiv.style.display = "none";
  successMessage.textContent = "";

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
