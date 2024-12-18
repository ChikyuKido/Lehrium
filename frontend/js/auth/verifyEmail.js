const API_BASE_URL = "{{ .ApiBaseUrl }}"

document.addEventListener('DOMContentLoaded', () => {
    const urlParams = new URLSearchParams(window.location.search);
    const uuid = urlParams.get('uuid');

    if (uuid) {
        console.log('UUID:', uuid);
        fetch(API_BASE_URL+'/api/v1/auth/verifyEmail?uuid=' + uuid, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({uuid})
        })
            .then(response => response.json())
            .then(data => {
                console.log('Success:', data);
            })
            .catch((error) => {
                console.error('Error:', error);
            });
    } else {
        console.error('UUID not found in the URL');
    }
});