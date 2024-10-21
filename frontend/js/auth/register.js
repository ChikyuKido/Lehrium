document.getElementById('register-form').addEventListener('submit', async (event) => {
    event.preventDefault();

    const email = document.getElementById('email').value;
    const password = document.getElementById('password').value;
    const confirm_password = document.getElementById('confirm-password').value;
    const untisName = document.getElementById('untisName').value;

    if (password !== confirm_password) {
        displayErrorMessage('Passwords do not match.');
        return;
    }

    const loginData = {
        email: email,
        password: password,
        untisName: untisName
    };

    try {
        const response = await fetch('http://localhost:8080/api/v1/auth/register', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(loginData)
        });

        if (response.ok) {
            const result = await response.json();
            window.location.replace("/");
            console.log('Success:', result);
        } else {
            displayErrorMessage('Registration failed: ' + response.statusText);
        }
    } catch (error) {
        displayErrorMessage('Error: ' + error.message);
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

    const form = document.getElementById('register-form');
    form.parentNode.insertBefore(errorNotification, form);
}
