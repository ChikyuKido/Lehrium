window.onload = function() {
    const token = localStorage.getItem('login_token'); // Get the token from localStorage

    if (token !== null) { // Check if the token is not null
        console.log('Token found:', token);
        window.location.replace("/succesfullLogin");
    } else {
        console.log('No token found, staying on the login page.');
    }
};
