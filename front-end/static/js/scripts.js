function loginAction() {
    event.preventDefault(); // Prevent the default form submission
    const email = document.getElementById('email').value;
    const password = document.getElementById('password').value;
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
    const email = document.getElementById('email').value;
    const username = document.getElementById('username').value;
    const password = document.getElementById('password').value;

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
            ID: postID,
            isLike: isLike
        }),
        contentType: "application/json",
        success: function (response) { //data is returnin as string
            // Assuming the server response contains the updated like/dislike status
            // data should include fields like data.isLiked and data.isDisliked
            // data should include fields like data.isLiked and data.isDisliked
            const likeBtn = $('#like-btn-' + postID);
            const dislikeBtn = $('#dislike-btn-' + postID);

            if (response === "liked") {
                likeBtn.addClass('custom-hover-like');
                dislikeBtn.removeClass('custom-hover-dislike');
            } else if (response === "disliked") {
                likeBtn.removeClass('custom-hover-like');
                dislikeBtn.addClass('custom-hover-dislike');
            } else {
                likeBtn.removeClass('custom-hover-like');
                dislikeBtn.removeClass('custom-hover-dislike');
            }
            updateCounters(postID);
        },
        error: function (response) {
            alert(response.responseText);
            $("#error").html(response.responseText); // Assuming your server sends plain text error messages
            $("#error").show();
            $("#error").removeClass("d-none");
            $("#error").fadeOut(5000);
        }
    });
}

function updateCounters(postID) {
    $.ajax({
        url: "/getPostLikesAndDislikesCount",
        type: "POST",
        data: JSON.stringify({
            ID: postID
        }),
        contentType: "application/json",
        success: function (data) {
            debugger;
            $('#like-count-' + postID).text(data.likes);
            $('#dislike-count-' + postID).text(data.dislikes);
        },
        error: function (data) {
            alert(data.responseText);
        }
    });
}



