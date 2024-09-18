document.getElementById('login-form').addEventListener('submit', async (event) => {
    event.preventDefault();

    const email = document.getElementById('email').value;
    const password = document.getElementById('password').value;
    const rememberMe = document.getElementById('remember_me').checked;

    const loginData = {
        email: email,
        password: password,
        rememberme: rememberMe
    };

    try {
        const response = await fetch('http://localhost:8080/auth/login', {
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
            window.location.replace("/succesfullLogin");
        } else {
            const errorText = await response.text();
            console.error('Error:', response.statusText);
            console.error('Response Text:', errorText);
        }
    } catch (error) {
        console.error('Fetch Error:', error);
    }
});
