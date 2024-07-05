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

{/* <form action="/likeOrDislikePost" method="post" class="d-inline">
<input type="hidden" name="postID" value="{{$post.ID}}">
<input type="hidden" name="isLike" value="like">
<button type="submit" class="btn btn-outline-success btn-sm
                {{if $post.Like.IsLiked}}
                custom-hover-like
                {{end}}">
    <i class="bi bi-hand-thumbs-up"></i>
</button>
</form> */}

function likeDislikePost(postID, isLike) {
    event.preventDefault(); //  the default form submission
    $.ajax({
        url: "/likeOrDislikePost",
        type: "POST",
        data: JSON.stringify({
            postID: postID,
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

