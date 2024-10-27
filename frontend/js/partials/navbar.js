document.addEventListener("DOMContentLoaded",function () {
    const authButtons = document.getElementById('authButtons')
    const settingsLink = document.getElementById('settingsLink')
    fetch("/api/v1/auth/pingAuth").then(response => {
        if(response.status === 200) {
            settingsLink.classList.remove('is-hidden');
        }else {
            authButtons.classList.remove('is-hidden');
        }
    })
})