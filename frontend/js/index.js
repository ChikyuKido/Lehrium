window.onload = function() {
    const token = localStorage.getItem('login_token');
    if (token !== null) {
        console.log('Token found:', token);
        window.location.replace("/succesfullLogin");
    } else {
        console.log('No token found, staying on the login page.');
    }
};
