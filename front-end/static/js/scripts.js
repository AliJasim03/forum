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
            $("#error").removeClass("d-none");
            $("#error").fadeOut(5000);
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
            debugger;
            window.location.href = "/";
        },
        error: function (data) {
            debugger;
            $("#error").html(data.responseText); // Assuming your server sends plain text error messages
            $("#error").show();
            $("#error").removeClass("d-none");
            $("#error").fadeOut(5000);
        }
    });
}

function likeDislikePost(postID, isLike) {
    event.preventDefault(); // Prevent the default form submission
    $.ajax({
        url: "/likeOrDislikePost",
        type: "POST",
        data: JSON.stringify({
            postID: postID,
            isLike: isLike
        }),
        contentType: "application/json",
        success: function (data) {
            // Assuming the server response contains the updated like/dislike status
            // data should include fields like data.isLiked and data.isDisliked
            const likeBtn = document.getElementById('like-btn-' + postID);
            const dislikeBtn = document.getElementById('dislike-btn-' + postID);

            if (data.isLiked) {
                likeBtn.classList.add('custom-hover-like');
                dislikeBtn.classList.remove('custom-hover-dislike');
            } else if (data.isDisliked) {
                likeBtn.classList.remove('custom-hover-like');
                dislikeBtn.classList.add('custom-hover-dislike');
            } else {
                likeBtn.classList.remove('custom-hover-like');
                dislikeBtn.classList.remove('custom-hover-dislike');
            }
        },
        error: function (data) {
            alert(data.responseText);
            $("#error").html(data.responseText); // Assuming your server sends plain text error messages
            $("#error").show();
            $("#error").removeClass("d-none");
            $("#error").fadeOut(5000);
        }
    });
}



function likeDislikeComment(commentID, isLike) {
    event.preventDefault(); // Prevent the default form submission
    $.ajax({
        url: "/likeOrDislikeComment",
        type: "POST",
        data: JSON.stringify({
            commentID: commentID,
            isLike: isLike
        }),
        contentType: "application/json",
        success: function (data) {
        },
        error: function (data) {
            alert(data.responseText);
            $("#error").html(data.responseText); // Assuming your server sends plain text error messages
            $("#error").show();
            $("#error").removeClass("d-none");
            $("#error").fadeOut(5000);
        }
    });
}

function createComment(postID) {
    event.preventDefault(); // Prevent the default form submission
    var comment = $("#comment").val();
    $.ajax({
        url: "/createComment",
        type: "POST",
        data: JSON.stringify({
            postID: postID,
            comment: comment
        }),
        contentType: "application/json",
        success: function (data) {
            alert("Comment created successfully");
        },
        error: function (data) {
            alert(data.responseText);
            $("#error").html(data.responseText); // Assuming your server sends plain text error messages
            $("#error").show();
            $("#error").removeClass("d-none");
            $("#error").fadeOut(5000);
        }
    });
}

