function submitComment(postID) {

    const comment = $("#comment").val();
    //check if the comment is empty or only contains spaces
    if (comment === "" || !comment.trim()) {
        $("#error").html("Comment cannot be empty");
        $("#error").show();
        $("#error").removeClass("d-none");
        $("#error").fadeOut(5000);
        return;
    }
    const data = JSON.stringify({
        PostID: postID,
        Comment: comment
    });
    $.ajax({
        url: "/createCommentAction",
        type: "POST",
        contentType: "application/json",
        data: data,  // Send the JSON data object
        success: function (response) {
            // Assuming commentText contains the text of the comment you want to add
            $("#comment").val("");  // Clear the comment text area
            addCommentToPage(response);
        },
        error: function (error) {
            $("#error").html(error.responseText);  // Display error message
            $("#error").show();
            $("#error").removeClass("d-none");
            $("#error").fadeOut(5000);
        }
    });
}

// Function to add a comment to the page
function addCommentToPage(comment) {

    // Create a new div element for the comment
    const newCommentDiv = document.createElement("div");
    newCommentDiv.className = "card mb-4";
    newCommentDiv.id = comment.ID;

    newCommentDiv.innerHTML =
`<div class="card-body">
    <p class="card-text">` + comment.Content + `</p>
</div>
<div class="card-footer text-muted">
    <div class="d-flex flex-row gap-2">
        <div style="margin-right: auto">
            Posted by You on ` + comment.CreatedOn + `
        </div>

        <button type="submit" onclick="likeDislikeComment('` + comment.ID + `','like')"
                id="like-btn-` + comment.ID + `"
                class="btn btn-outline-success btn-sm">
            <p id="like-count-` + comment.ID + `" class="d-inline">
                ` + comment.Like.CountLikes + `
            </p>
            <i class="bi bi-hand-thumbs-up"></i>
        </button>

        <button type="submit" onclick="likeDislikeComment('` + comment.ID + `','dislike')"
                id="dislike-btn-` + comment.ID + `" class="btn btn-outline-danger btn-sm">
            <p id="dislike-count-` + comment.ID + `" class="d-inline">
                ` + comment.Like.CountDislikes + `
            </p>
            <i class="bi bi-hand-thumbs-down"></i>
        </button>
    </div>
</div>`;

    // Append the new comment div to the container
    $("#commentsContainer").append(newCommentDiv);
}

// Example usage
const comment = {
    ID: "123",
    Content: "This is a new comment",
    CreatedBy: "User",
    CreatedOn: "2024-07-06",
    Like: {
        IsLiked: false,
        IsDisliked: false
    }
};

//addCommentToPage(comment);
$("input").on("keypress", function () {
    $("#error").hide();
    $("#success").hide();
    $("#error").addClass("d-none");
    $("#success").addClass("d-none");
});

function likeDislikeComment(commentID, likeDislike) {
    const data = JSON.stringify({
        ID: commentID,
        isLike: likeDislike
    });
    $.ajax({
        url: "/likeOrDislikeComment",
        type: "POST",
        contentType: "application/json",
        data: data,  // Send the JSON data object
        success: function (response) {
            const likeBtn = $('#like-btn-' + commentID);
            const dislikeBtn = $('#dislike-btn-' + commentID);

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

            updateCommentCounters(commentID);
        },
        error: function (error) {
            alert(error.responseText);
            $("#error").html(error.responseText);  // Display error message
            $("#error").show();
            $("#error").removeClass("d-none");
            $("#error").fadeOut(5000);
        }
    });
}

function updateCommentCounters(commentID) {
    $.ajax({
        url: "/getCommentLikeDislikeCount",
        type: "POST",
        data: JSON.stringify({
            ID: commentID
        }),
        success: function (response) {
            $('#like-count-' + commentID).text(response.likes);
            $('#dislike-count-' + commentID).text(response.dislikes);
        },
        error: function (error) {
            alert(error.responseText);
        }
    });
}
