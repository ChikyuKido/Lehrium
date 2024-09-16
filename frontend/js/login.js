document.getElementById('login-form').addEventListener('submit', async (event) => {
    event.preventDefault();

    const email = document.getElementById('email').value;
    const password = document.getElementById('password').value;
    const untisName = document.getElementById('untisName').value;

    const loginData = {
        email: email,
        password: password,
        untisName: untisName
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
            console.log('Response JSON:', result);
        } else {
            const errorText = await response.text();
            console.error('Error:', response.statusText);
            console.error('Response Text:', errorText);
        }
    } catch (error) {
        console.error('Fetch Error:', error);
    }
});
