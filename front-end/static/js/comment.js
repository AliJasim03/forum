function submitComment(postID) {

    var comment = $("#comment").val();
    var data = JSON.stringify({
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
    var newCommentDiv = document.createElement("div");
    newCommentDiv.className = "card mb-4";
    newCommentDiv.id = comment.ID;

    // Create the card body
    var cardBodyDiv = document.createElement("div");
    cardBodyDiv.className = "card-body";

    // Create the card text
    var cardTextP = document.createElement("p");
    cardTextP.className = "card-text";
    cardTextP.textContent = comment.Content;

    // Append card text to card body
    cardBodyDiv.appendChild(cardTextP);

    // Create the card footer
    var cardFooterDiv = document.createElement("div");
    cardFooterDiv.className = "card-footer text-muted";
    cardFooterDiv.textContent = "Posted by " + comment.CreatedBy + " on " + comment.CreatedOn;

    // Create the hidden inputs and buttons
    var colDiv = document.createElement("div");
    colDiv.className = "col";

    var hiddenCommentID = document.createElement("input");
    hiddenCommentID.type = "hidden";
    hiddenCommentID.name = "commentID";
    hiddenCommentID.id = "commentID";
    hiddenCommentID.value = comment.ID;

    var hiddenIsDislike = document.createElement("input");
    hiddenIsDislike.type = "hidden";
    hiddenIsDislike.name = "isLike";
    hiddenIsDislike.id = "isLike";
    hiddenIsDislike.value = "dislike";

    var dislikeButton = document.createElement("button");
    dislikeButton.type = "submit";
    dislikeButton.className = "btn btn-outline-danger btn-sm";
    if (comment.Like.IsDisliked) {
        dislikeButton.classList.add("custom-hover-dislike");
    }
    dislikeButton.onclick = function () { likeDislikeComment(comment.ID, 'dislike'); };
    dislikeButton.innerHTML = '<i class="bi bi-hand-thumbs-down"></i>';

    var hiddenIsLike = document.createElement("input");
    hiddenIsLike.type = "hidden";
    hiddenIsLike.name = "isLike";
    hiddenIsLike.value = "like";

    var likeButton = document.createElement("button");
    likeButton.type = "submit";
    likeButton.className = "btn btn-outline-success btn-sm";
    if (comment.Like.IsLiked) {
        likeButton.classList.add("custom-hover-like");
    }
    likeButton.onclick = function () { likeDislikeComment(comment.ID, 'like'); };
    likeButton.innerHTML = '<i class="bi bi-hand-thumbs-up"></i>';

    // Append hidden inputs and buttons to the colDiv
    colDiv.appendChild(hiddenCommentID);
    colDiv.appendChild(hiddenIsDislike);
    colDiv.appendChild(dislikeButton);
    colDiv.appendChild(hiddenIsLike);
    colDiv.appendChild(likeButton);

    // Append colDiv to card footer
    cardFooterDiv.appendChild(colDiv);

    // Append card body and card footer to new comment div
    newCommentDiv.appendChild(cardBodyDiv);
    newCommentDiv.appendChild(cardFooterDiv);

    // Append the new comment div to the container
    $("#commentsContainer").append(newCommentDiv);
}

// Example usage
var comment = {
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
    var data = JSON.stringify({
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
            debugger;
            $('#like-count-' + commentID).text(response.likes);
            $('#dislike-count-' + commentID).text(response.dislikes);
        },
        error: function (error) {
            alert(error.responseText);
        }
    });
}
