function loginAction() {
    event.preventDefault(); // Prevent the default form submission
    var email = document.getElementById('email').value;
    var password = document.getElementById('password').value;
    $.ajax({
        url: "/loginAction",
        type: "POST",
        data: JSON.stringify({
            email: email,
            password: password
        }),
        contentType: "application/json",
        success: function (data) {
            window.location.href = "/";
        },
        error: function (data) {
            $("#error").html(data.responseText); // Assuming your server sends plain text error messages
            $("#error").show();
        }
    });
}

function registerAction() {
    event.preventDefault(); // Prevent the default form submission
    var email = document.getElementById('email').value;
    var username = document.getElementById('username').value;
    var password = document.getElementById('password').value;
    $.ajax({
        url: "/registerAction",
        type: "POST",
        data: JSON.stringify({
            email: email,
            username: username,
            password: password
        }),
        contentType: "application/json",
        success: function (data) {
            window.location.href = "/";
        },
        error: function (data) {
            $("#error").html(data.responseText); // Assuming your server sends plain text error messages
            $("#error").show();
        }
    });
}




