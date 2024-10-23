document.getElementById('login-form').addEventListener('submit', async (event) => {
    event.preventDefault();

    const email = document.getElementById('email').value;
    const password = document.getElementById('password').value;
    const rememberMe = document.getElementById('remember_me').checked;

    const loginData = {
        email: email,
        password: password,
        rememberMe: rememberMe
    };

    try {
        const response = await fetch('http://localhost:8080/api/v1/auth/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            credentials: 'include',
            body: JSON.stringify(loginData)
        });

        if (response.ok) {
            const result = await response.json();
            localStorage.setItem("login_token", result.token);
            console.log('Response JSON:', result);
            window.location.replace("/");
        } else {
            const errorText = await response.text();
            displayErrorMessage(`Login failed: ${errorText}`);
        }
    } catch (error) {
        displayErrorMessage(`Fetch Error: ${error.message}`);
    }
});

function displayErrorMessage(message) {
    let existingNotification = document.getElementById('error-notification');
    if (existingNotification) {
        existingNotification.remove();
    }

    const errorNotification = document.createElement('div');
    errorNotification.id = 'error-notification';
    errorNotification.className = 'notification is-danger';
    errorNotification.innerHTML = `
        <button class="delete" onclick="this.parentElement.remove()"></button>
        ${message}
    `;
    const form = document.getElementById('login-form');
    form.parentNode.insertBefore(errorNotification, form);
}
