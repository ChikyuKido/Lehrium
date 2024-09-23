        document.getElementById('register-form').addEventListener('submit', async (event) => {
            event.preventDefault();

            const email = document.getElementById('email').value;
            const password = document.getElementById('password').value;
            const untisName = document.getElementById('untisName').value;

            const loginData = {
                email: email,
                password: password
            };

            try {
                const response = await fetch('http://localhost:8080/auth/register', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify(loginData)
                });

                if (response.ok) {
                    const result = await response.json();
                    window.location.replace("/auth/succesfullRegister");
                    console.log('Success:', result);
                } else {
                    console.error('Error:', response.statusText);
                }
            } catch (error) {
                console.error('Error:', error);
            }
        });
