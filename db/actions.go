package db

import (
	"database/sql"
	"fmt"
	"time"
)

func GetPosts(db *sql.DB, user int, posts *[]Post) {
	rows, err := db.Query("SELECT id, user_id, title, content, strftime('%Y-%m-%d %H:%M:%S', created_at) AS created_at FROM posts")
	if err != nil {
		fmt.Println(err) 
		return
	}
	defer rows.Close()

	for rows.Next() {
		post := &Post{}
		var userId int
		err = rows.Scan(&post.ID, &userId, &post.Title, &post.Content, &post.CreatedOn)
		if err != nil {
			fmt.Println(err)
			return
		}
		post.CreatedBy = GetUsername(db, userId)
		post.IsCreatedByUser = user == userId
		*posts = append(*posts, *post)
	}

	// get the categories of the post_categories
	for i := range *posts {
		rows, err := db.Query("SELECT c.name FROM categories c INNER JOIN post_categories pc ON c.id = pc.category_id WHERE pc.post_id = ?", (*posts)[i].ID)
		if err != nil {
			fmt.Println(err)
			return
		}

		var categories []string
		for rows.Next() {
			var category string
			err = rows.Scan(&category)
			if err != nil {
				fmt.Println(err)
				return
			}
			categories = append(categories, category)
		}
		if err := rows.Err(); err != nil {
			fmt.Println(err)
			return
		}
		rows.Close()

		(*posts)[i].Categories = categories
	}

	// get the likes and dislikes
	for i := range *posts {
		var likes, dislikes int
		err = db.QueryRow("SELECT COUNT(*) FROM likes WHERE post_id = ? AND is_like = ?", (*posts)[i].ID, true).Scan(&likes)
		if err != nil {
			fmt.Println(err)
			return
		}
		err = db.QueryRow("SELECT COUNT(*) FROM likes WHERE post_id = ? AND is_like = ?", (*posts)[i].ID, false).Scan(&dislikes)
		if err != nil {
			fmt.Println(err)
			return
		}
		(*posts)[i].Like.CountLikes = likes
		(*posts)[i].Like.CountDislikes = dislikes
	}

	if user != -1 {
		// get the likes and dislikes of the user
		for i := range *posts {
			var isLiked, isDisliked int
			err = db.QueryRow("SELECT COUNT(*) FROM likes WHERE post_id = ? AND user_id = ? AND is_like = ?", (*posts)[i].ID, user, true).Scan(&isLiked)
			if err != nil {
				fmt.Println(err)
				return
			}
			err = db.QueryRow("SELECT COUNT(*) FROM likes WHERE post_id = ? AND user_id = ? AND is_like = ?", (*posts)[i].ID, user, false).Scan(&isDisliked)
			if err != nil {
				fmt.Println(err)
				return
			}
			(*posts)[i].Like.IsLiked = isLiked > 0
			(*posts)[i].Like.IsDisliked = isDisliked > 0
		}
	}
	// comments in each post
	/*
		for i := range *posts {
			rows, err := db.Query("SELECT * FROM comments WHERE post_id = ?", (*posts)[i].ID)
			if err != nil {
				fmt.Println(err) return
			}
			defer rows.Close()
			for rows.Next() {
				comment := &Comments{}
				var userId int
				err = rows.Scan(&comment.ID, &comment.PostID, &userId, &comment.Content, &comment.CreatedOn)
				if err != nil {
					fmt.Println(err) return
				}
				comment.CreatedBy = GetUsername(db, userId)
				// convert user to
				if user == userId {
					comment.IsCreatedByUser = true
				}
				(*posts)[i].Comments = append((*posts)[i].Comments, *comment)
			}
			// get the likes and dislikes of the comments
			for j := range (*posts)[i].Comments {
				var likes, dislikes int
				err = db.QueryRow("SELECT COUNT(*) FROM likes WHERE comment_id = ? AND is_like = ?", (*posts)[i].Comments[j].ID, true).Scan(&likes)
				if err != nil {
					fmt.Println(err) return
				}
				err = db.QueryRow("SELECT COUNT(*) FROM likes WHERE comment_id = ? AND is_like = ?", (*posts)[i].Comments[j].ID, false).Scan(&dislikes)
				if err != nil {
					fmt.Println(err) return
				}
				(*posts)[i].Comments[j].Like.CountLikes = likes
				(*posts)[i].Comments[j].Like.CountDislikes = dislikes
			}
			if user != -1 {
				// get the likes and dislikes of the user
				for j := range (*posts)[i].Comments {
					var isLiked, isDisliked int
					err = db.QueryRow("SELECT COUNT(*) FROM likes WHERE comment_id = ? AND user_id = ? AND is_like = ?", (*posts)[i].Comments[j].ID, user, true).Scan(&isLiked)
					if err != nil {
						fmt.Println(err) return
					}
					err = db.QueryRow("SELECT COUNT(*) FROM likes WHERE comment_id = ? AND user_id = ? AND is_like = ?", (*posts)[i].Comments[j].ID, user, false).Scan(&isDisliked)
					if err != nil {
						fmt.Println(err) return
					}
					(*posts)[i].Comments[j].Like.IsLiked = isLiked > 0
					(*posts)[i].Comments[j].Like.IsDisliked = isDisliked > 0
				}
			}
		}
	*/
}

func GetPost(db *sql.DB, user int, post *Post) {
	var userIdTemp int

	row := db.QueryRow("SELECT id, user_id, title, content, strftime('%Y-%m-%d %H:%M:%S', created_at) FROM posts WHERE id = ?", post.ID)
	if row == nil {
		fmt.Println(row)
		return
	}

	err := row.Scan(&post.ID, &userIdTemp, &post.Title, &post.Content, &post.CreatedOn)
	if err != nil {
		fmt.Println(row)
		return
	}

	post.CreatedBy = GetUsername(db, userIdTemp)

	// get the categories of the post_categories
	rows, err := db.Query("SELECT name FROM categories INNER JOIN post_categories ON categories.id = post_categories.category_id WHERE post_categories.post_id = ?", post.ID)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var category string
		err = rows.Scan(&category)
		if err != nil {
			fmt.Println(err)
			return
		}
		post.Categories = append(post.Categories, category)
	}
	// get the likes and dislikes
	var likes, dislikes int
	err = db.QueryRow("SELECT COUNT(*) FROM likes WHERE post_id = ? AND is_like = ?", post.ID, true).Scan(&likes)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = db.QueryRow("SELECT COUNT(*) FROM likes WHERE post_id = ? AND is_like = ?", post.ID, false).Scan(&dislikes)
	if err != nil {
		fmt.Println(err)
		return
	}
	post.Like.CountLikes = likes
	post.Like.CountDislikes = dislikes

	if user != -1 {
		// get the likes and dislikes of the user
		var isLiked, isDisliked int
		err = db.QueryRow("SELECT COUNT(*) FROM likes WHERE post_id = ? AND user_id = ? AND is_like = ?", post.ID, user, true).Scan(&isLiked)
		if err != nil {
			fmt.Println(err)
			return
		}
		err = db.QueryRow("SELECT COUNT(*) FROM likes WHERE post_id = ? AND user_id = ? AND is_like = ?", post.ID, user, false).Scan(&isDisliked)
		if err != nil {
			fmt.Println(err)
			return
		}
		post.Like.IsLiked = isLiked > 0
		post.Like.IsDisliked = isDisliked > 0
	}
	// comments in each post

	rows, err = db.Query("SELECT id, post_id, user_id, content, strftime('%Y-%m-%d %H:%M:%S', created_at) AS created_at FROM comments WHERE post_id = ?", post.ID)
	if err != nil {
		fmt.Println(err) 
		return
		
	}
	defer rows.Close()
	for rows.Next() {
		comment := &Comment{}
		var userId int
		err = rows.Scan(&comment.ID, &comment.PostID, &userId, &comment.Content, &comment.CreatedOn)
		if err != nil {
			fmt.Println(err) 
			return
			
		}
		comment.CreatedBy = GetUsername(db, userId)
		// convert user to
		if user == userId {
			comment.IsCreatedByUser = true
		}

		// get the likes and dislikes of the comments
		var likes, dislikes int
		err = db.QueryRow("SELECT COUNT(*) FROM likes WHERE comment_id = ? AND is_like = ?", comment.ID, true).Scan(&likes)
		if err != nil {
			fmt.Println(err) 
			return
		}
		err = db.QueryRow("SELECT COUNT(*) FROM likes WHERE comment_id = ? AND is_like = ?", comment.ID, false).Scan(&dislikes)
		if err != nil {
			fmt.Println(err) 
			return
		}
		comment.Like.CountLikes = likes
		comment.Like.CountDislikes = dislikes

		if user != -1 {
			// get the likes and dislikes of the user
			var isLiked, isDisliked int
			err = db.QueryRow("SELECT COUNT(*) FROM likes WHERE comment_id = ? AND user_id = ? AND is_like = ?", comment.ID, user, true).Scan(&isLiked)
			if err != nil {
				fmt.Println(err) 
				return
			}
			err = db.QueryRow("SELECT COUNT(*) FROM likes WHERE comment_id = ? AND user_id = ? AND is_like = ?", comment.ID, user, false).Scan(&isDisliked)
			if err != nil {
				fmt.Println(err) 
				return
			}
			comment.Like.IsLiked = isLiked > 0
			comment.Like.IsDisliked = isDisliked > 0
		}
		post.Comments = append(post.Comments, *comment)

	}
}

func GetUsername(db *sql.DB, id int) string {
	// get username from db
	var username string
	err := db.QueryRow("SELECT username FROM users WHERE id = ?", id).Scan(&username)
	if err != nil {
		fmt.Println(err) 
		return ""
	}
	return username
}

func LikeDislikePost(db *sql.DB, userID int, postID string, isLike bool) bool {
	// select and checked the saved like
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM likes WHERE user_id = ? AND post_id = ? LIMIT 1)", userID, postID).Scan(&exists)
	if err != nil {
		fmt.Println(err)
		return false
	}
	if exists {
		// reverse the value of isliked
		var oldVal bool
		err = db.QueryRow("SELECT is_like FROM likes WHERE user_id = ? AND post_id = ?", userID, postID).Scan(&oldVal)
		if err != nil {
			fmt.Println(err) 
			return false
		}
		// delete if the value is the same
		if oldVal == isLike {
			_, err = db.Exec("DELETE FROM likes WHERE user_id = ? AND post_id = ?", userID, postID)
			if err != nil {
				fmt.Println(err) 
				return false
			}
			return true
		}
		_, err = db.Exec("UPDATE likes SET is_like = ? WHERE user_id = ? AND post_id = ?", !oldVal, userID, postID)
		if err != nil {
			fmt.Println(err) 
			return false
		}
		return true
	}
	_, err = db.Exec("INSERT INTO likes (user_id, post_id, is_like, created_at) VALUES(?,?,?,?)", userID, postID, isLike, time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		fmt.Println(err) 
		return false
	}
	return true
}

func KnowPostLike(db *sql.DB, userID int, postID string) string {
	var isLiked bool
	err := db.QueryRow("SELECT is_like FROM likes WHERE user_id = ? AND post_id = ?", userID, postID).Scan(&isLiked)
	if err != nil {
		return "none"
	}
	if isLiked {
		return "liked"
	}
	return "disliked"
}

func CreatePost(db *sql.DB, userID int, post PostJson) int {
	// create the post
	_, err := db.Exec("INSERT INTO posts (user_id, title, content) VALUES (?, ?, ?)", userID, post.Title, post.Content)
	if err != nil {
		return -1
	}

	var postID int
	err = db.QueryRow("SELECT id FROM posts WHERE user_id = ? AND title = ? ORDER BY created_at DESC LIMIT 1", userID, post.Title).Scan(&postID)
	if err != nil {
		return -1
	}

	// add the category
	for _, ct := range post.Categories {
		if ct != "" { // check if not empty
			var categoryID int
			err = db.QueryRow("SELECT id FROM categories WHERE name = ?", ct).Scan(&categoryID)
			if err != nil {
				return -1
			}

			// Check if the combination of post_id and category_id already exists
			var exists bool
			err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM post_categories WHERE post_id = ? AND category_id = ?)", postID, categoryID).Scan(&exists)
			if err != nil {
				// handle error
				return -1
			}

			if !exists {
				_, err = db.Exec("INSERT INTO post_categories (post_id, category_id) VALUES (?, ?)", postID, categoryID)
				if err != nil {
					// handle error
					return -1
				}
			} else {
				continue
			}
		}
	}
	return postID
}

func GetCategories(db *sql.DB) []Category {
	var categories []Category
	rows, err := db.Query("SELECT id, name FROM categories")
	if err != nil {
		fmt.Println(err)
		return []Category{}
	}
	defer rows.Close()

	for rows.Next() {
		category := Category{}
		err = rows.Scan(&category.ID, &category.Name)
		if err != nil {
			fmt.Println(err)
			return []Category{}
		}
		categories = append(categories, category)
	}
	return categories
}

func CreateComment(db *sql.DB, userID int, comment CommentJson) (bool, *Comment) {
	// create the comment
	_, err := db.Exec("INSERT INTO comments (post_id, user_id, content, created_at) VALUES (?, ?, ?, ?)", comment.PostID, userID, comment.Comment, time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		return false, nil
	}
	commentDB := &Comment{}
	err = db.QueryRow("SELECT id, post_id, user_id, content, created_at FROM comments WHERE user_id = ? AND content = ? ORDER BY created_at DESC LIMIT 1", userID, comment.Comment).Scan(&commentDB.ID, &commentDB.PostID, &commentDB.CreatedBy, &commentDB.Content, &commentDB.CreatedOn)

	return true, commentDB
}

func LikeDislikeComment(db *sql.DB, userID int, commentID string, isLike bool) bool {
	// select and checked the saved like
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM likes WHERE user_id = ? AND comment_id = ? LIMIT 1)", userID, commentID).Scan(&exists)
	if err != nil {
		fmt.Println(err)
		return false
	}
	if exists {
		// reverse the value of isliked
		var oldVal bool
		err := db.QueryRow("SELECT is_like FROM likes WHERE user_id = ? AND comment_id = ?", userID, commentID).Scan(&oldVal)
		if err != nil {
			fmt.Println(err)
			return false
		}
		// delete if the value is the same
		if oldVal == isLike {
			_, err = db.Exec("DELETE FROM likes WHERE user_id = ? AND comment_id = ?", userID, commentID)
			if err != nil {
				fmt.Println(err)
				return false
			}
			return true
		}
		_, err = db.Exec("UPDATE likes SET is_like = ? WHERE user_id = ? AND comment_id = ?", !oldVal, userID, commentID)
		if err != nil {
			fmt.Println(err)
			return false
		}
		return true
	}
	_, err = db.Exec("INSERT INTO likes (user_id, comment_id, is_like, created_at) VALUES(?,?,?,?)", userID, commentID, isLike, time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func KnowCommentLike(db *sql.DB, userID int, commnetID string) string {
	var isLiked bool
	err := db.QueryRow("SELECT is_like FROM likes WHERE user_id = ? AND comment_id = ?", userID, commnetID).Scan(&isLiked)
	if err != nil {
		return "none"
	}
	if isLiked {
		return "liked"
	}
	return "disliked"
}

func GetPostLikesAndDislikesCount(db *sql.DB, postID string) (int, int) {
	var likes, dislikes int
	err := db.QueryRow("SELECT COUNT(*) FROM likes WHERE post_id = ? AND is_like = ?", postID, true).Scan(&likes)
	if err != nil {
		fmt.Println(err)
		return 0, 0
	}
	err = db.QueryRow("SELECT COUNT(*) FROM likes WHERE post_id = ? AND is_like = ?", postID, false).Scan(&dislikes)
	if err != nil {
		fmt.Println(err)
		return 0, 0
	}
	return likes, dislikes
}

func GetCommentLikesAndDislikesCount(db *sql.DB, commentID string) (int, int) {
	var likes, dislikes int
	err := db.QueryRow("SELECT COUNT(*) FROM likes WHERE comment_id = ? AND is_like = ?", commentID, true).Scan(&likes)
	if err != nil {
		fmt.Println(err)
		return 0, 0
	}
	err = db.QueryRow("SELECT COUNT(*) FROM likes WHERE comment_id = ? AND is_like = ?", commentID, false).Scan(&dislikes)
	if err != nil {
		fmt.Println(err)
		return 0, 0
	}
	return likes, dislikes
}
