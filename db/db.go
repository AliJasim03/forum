package db

import (
	"database/sql"
	_ "github.com/gofrs/uuid"
	_ "github.com/mattn/go-sqlite3"
	_ "golang.org/x/crypto/bcrypt"
	"log"
	"os"
	"time"
)

func OpenConnection() *sql.DB {
	db, err := sql.Open("sqlite3", "db/forum.db")
	if err != nil {
		log.Fatal(err)
	}

	return initDB(db)
}

func CloseDB(db *sql.DB) {
	db.Close()
}

func initDB(db *sql.DB) *sql.DB {
	sqlFile, err := os.ReadFile("db/init.sql")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(string(sqlFile))
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func GetPosts(db *sql.DB, user int, posts *[]Post, filter string) {

	rows, err := db.Query("SELECT * FROM posts")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var post = &Post{}
		var userId int
		err = rows.Scan(&post.ID, &userId, &post.Title, &post.Content, &post.CreatedOn)
		if err != nil {
			log.Fatal(err)
		}
		post.CreatedBy = GetUsername(db, userId)
		//convert user to
		if user == userId {
			post.IsCreatedByUser = true
		}
		*posts = append(*posts, *post)
	}

	//get the categories of the post_categories
	for i := range *posts {
		rows, err := db.Query("SELECT name FROM categories INNER JOIN post_categories ON categories.id = post_categories.category_id WHERE post_categories.post_id = ?", (*posts)[i].ID)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		for rows.Next() {
			var category string
			err = rows.Scan(&category)
			if err != nil {
				log.Fatal(err)
			}
			(*posts)[i].Categories = append((*posts)[i].Categories, category)
		}
	}
	//get the likes and dislikes
	for i := range *posts {
		var likes, dislikes int
		err = db.QueryRow("SELECT COUNT(*) FROM likes WHERE post_id = ? AND is_like = ?", (*posts)[i].ID, true).Scan(&likes)
		if err != nil {
			log.Fatal(err)
		}
		err = db.QueryRow("SELECT COUNT(*) FROM likes WHERE post_id = ? AND is_like = ?", (*posts)[i].ID, false).Scan(&dislikes)
		if err != nil {
			log.Fatal(err)
		}
		(*posts)[i].Like.CountLikes = likes
		(*posts)[i].Like.CountDislikes = dislikes
	}

	if user != -1 {
		//get the likes and dislikes of the user
		for i := range *posts {
			var isLiked, isDisliked int
			err = db.QueryRow("SELECT COUNT(*) FROM likes WHERE post_id = ? AND user_id = ? AND is_like = ?", (*posts)[i].ID, user, true).Scan(&isLiked)
			if err != nil {
				log.Fatal(err)
			}
			err = db.QueryRow("SELECT COUNT(*) FROM likes WHERE post_id = ? AND user_id = ? AND is_like = ?", (*posts)[i].ID, user, false).Scan(&isDisliked)
			if err != nil {
				log.Fatal(err)
			}
			(*posts)[i].Like.IsLiked = isLiked > 0
			(*posts)[i].Like.IsDisliked = isDisliked > 0
		}
	}
	//comments in each post
	for i := range *posts {
		rows, err := db.Query("SELECT * FROM comments WHERE post_id = ?", (*posts)[i].ID)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		for rows.Next() {
			var comment = &Comments{}
			var userId int
			err = rows.Scan(&comment.ID, &comment.PostID, &userId, &comment.Content, &comment.CreatedOn)
			if err != nil {
				log.Fatal(err)
			}
			comment.CreatedBy = GetUsername(db, userId)
			//convert user to
			if user == userId {
				comment.IsCreatedByUser = true
			}
			(*posts)[i].Comments = append((*posts)[i].Comments, *comment)
		}
		//get the likes and dislikes of the comments
		for j := range (*posts)[i].Comments {
			var likes, dislikes int
			err = db.QueryRow("SELECT COUNT(*) FROM likes WHERE comment_id = ? AND is_like = ?", (*posts)[i].Comments[j].ID, true).Scan(&likes)
			if err != nil {
				log.Fatal(err)
			}
			err = db.QueryRow("SELECT COUNT(*) FROM likes WHERE comment_id = ? AND is_like = ?", (*posts)[i].Comments[j].ID, false).Scan(&dislikes)
			if err != nil {
				log.Fatal(err)
			}
			(*posts)[i].Comments[j].Like.CountLikes = likes
			(*posts)[i].Comments[j].Like.CountDislikes = dislikes
		}
		if user != -1 {
			//get the likes and dislikes of the user
			for j := range (*posts)[i].Comments {
				var isLiked, isDisliked int
				err = db.QueryRow("SELECT COUNT(*) FROM likes WHERE comment_id = ? AND user_id = ? AND is_like = ?", (*posts)[i].Comments[j].ID, user, true).Scan(&isLiked)
				if err != nil {
					log.Fatal(err)
				}
				err = db.QueryRow("SELECT COUNT(*) FROM likes WHERE comment_id = ? AND user_id = ? AND is_like = ?", (*posts)[i].Comments[j].ID, user, false).Scan(&isDisliked)
				if err != nil {
					log.Fatal(err)
				}
				(*posts)[i].Comments[j].Like.IsLiked = isLiked > 0
				(*posts)[i].Comments[j].Like.IsDisliked = isDisliked > 0
			}
		}
	}

}

func GetUsername(db *sql.DB, id int) string {
	//get username from db
	var username string
	err := db.QueryRow("SELECT username FROM users WHERE id = ?", id).Scan(&username)
	if err != nil {
		log.Fatal(err)
	}
	return username
}

func LikeDislikePost(db *sql.DB, userID int, postID string, isLike bool) bool {
	//select and checked the saved like
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM likes WHERE user_id = ? AND post_id = ? LIMIT 1)", userID, postID).Scan(&exists)
	if err != nil {
		log.Fatal(err)
		return false
	}
	if exists {
		//reverse the value of isliked
		var oldVal bool
		err := db.QueryRow("SELECT is_like FROM likes WHERE user_id = ? AND post_id = ?", userID, postID).Scan(&oldVal)
		if err != nil {
			log.Fatal(err)
			return false
		}
		//delete if the value is the same
		if oldVal == isLike {
			_, err = db.Exec("DELETE FROM likes WHERE user_id = ? AND post_id = ?", userID, postID)
			if err != nil {
				log.Fatal(err)
				return false
			}
			return true
		}
		_, err = db.Exec("UPDATE likes SET is_like = ? WHERE user_id = ? AND post_id = ?", !oldVal, userID, postID)
		if err != nil {
			log.Fatal(err)
			return false
		}
		return true
	}
	_, err = db.Exec("INSERT INTO likes (user_id, post_id, is_like, created_at) VALUES(?,?,?,?)", userID, postID, isLike, time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		log.Fatal(err)
		return false
	}
	return true
}
