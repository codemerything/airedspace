let container = document.querySelector(".app-container");

if (window.location.pathname.includes("/search")) {
    container.classList.add("-search");
}