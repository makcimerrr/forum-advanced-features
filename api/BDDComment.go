package api

import (
	"database/sql"
)

func SetComments(db *sql.DB, discussionIDInt int, message string, idUser int) error {
	// Insérez le nouveau message dans la base de données en incluant l'ID de discussion
	_, err := db.Exec("INSERT INTO comment (discussion_id, user_id, message) VALUES (?, ?, ?)", discussionIDInt, idUser, message)

	return err
}

func GetCommentsFromDiscussion(db *sql.DB, discussionIDInt int)([]Comment, error){
	rows, err := db.Query("SELECT comment.id, user.username AS users, comment.message FROM comment JOIN user ON comment.user_id = user.id WHERE discussion_id = ?", discussionIDInt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var comments []Comment

	// Parcourez les commentaires et ajoutez-les à la structure de données
	for rows.Next() {
		var comment Comment
		err := rows.Scan(&comment.ID, &comment.Username, &comment.Message) // Assurez-vous d'avoir une structure Comment avec des champs User et Message
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

func GetOneCommentById(db *sql.DB, id int)(Comment, error){
	var comment Comment
	err := db.QueryRow("SELECT comment.id, user.username AS users, comment.discussion_id, comment.message FROM comment JOIN user ON comment.user_id = user.id WHERE comment.id = ?", id).Scan(&comment.ID, &comment.Username, &comment.Discussion_id, &comment.Message)
	return comment, err
}

func CheckNumberOfCommentForDiscussion(db *sql.DB, discussionID int) (int, error){
	var commentCount int
	err := db.QueryRow("SELECT COUNT(*) FROM comment WHERE discussion_id = ?", discussionID).Scan(&commentCount)
	if err != nil {
		return 0, err
	}

	return commentCount, nil
}

func GetAllCommentForOneUser(db *sql.DB, idUser int) ([]Comment, error){
	rows, err := db.Query("SELECT comment.id, user.username AS users, comment.message, comment.discussion_id FROM comment JOIN user ON comment.user_id = user.id WHERE comment.user_id = ?", idUser)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []Comment

	// Parcourez les résultats et stockez-les dans la slice
	for rows.Next() {
		var comment Comment
		err := rows.Scan(&comment.ID, &comment.Username, &comment.Message, &comment.Discussion_id)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}



