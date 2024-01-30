package api

import "database/sql"


//function for table dislike Discussion

func GetDislikesFromOneDiscussion(db *sql.DB, discussionIDInt int, username int) (bool, error) {
	var disliked bool
	err := db.QueryRow("SELECT 1 FROM dislikeDiscussion WHERE discussion_id = ? AND user_id = ?", discussionIDInt, username).Scan(&disliked)
	return disliked, err
}

func SetDisLikesDiscussion(db *sql.DB, discussionIDInt int, username int) error {
	_, err := db.Exec("INSERT INTO dislikeDiscussion (discussion_id, user_id) VALUES (?, ?)", discussionIDInt, username)
	return err
}

func DeleteDislikeByUserIdForDiscussion(db *sql.DB, discussionIDInt int, username int) error {
	_, err := db.Exec("DELETE FROM dislikeDiscussion WHERE discussion_id = ? AND user_id = ?", discussionIDInt, username)
	return err
}

func CheckIfUserDislikedDiscussion(db *sql.DB, idUser int, discussionID int) (bool, error) {

	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM dislikeDiscussion WHERE discussion_id = ? AND user_id = ?)", discussionID, idUser).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func CheckNumberOfDislikesForDiscussion(db *sql.DB, discussionID int) (int, error) {

	var likeCount int
	err := db.QueryRow("SELECT COUNT(*) FROM dislikeDiscussion WHERE discussion_id = ?", discussionID).Scan(&likeCount)
	if err != nil {
		return 0, err
	}

	return likeCount, nil
}

func DeleteDislikeFromDiscussion(db *sql.DB, discussionIDInt int) error {
	_, err := db.Exec("DELETE FROM dislikeDiscussion WHERE discussion_id = ?", discussionIDInt)
	return err
}


//function for table like Comment

func GetDislikeFromOneComment(db *sql.DB, commentId int, username int)(bool, error){
	var liked bool
	err := db.QueryRow("SELECT 1 FROM dislikeComment WHERE comment_id = ? AND user_id = ?", commentId, username).Scan(&liked)
	return liked, err
}

func SetDislikeComment(db *sql.DB, discussionID int, commentId int, username int) (error){
	_, err := db.Exec("INSERT INTO dislikeComment (discussion_id, comment_id, user_id) VALUES (?, ?, ?)", discussionID, commentId, username)
	return err
}

func DeleteDislikeByUserIdForComment(db *sql.DB, commentId int, username int) (error) {
	_, err := db.Exec("DELETE FROM dislikeComment WHERE comment_id = ? AND user_id = ?", commentId, username)
	return err
}

func CheckIfUserDislikeComment(db *sql.DB, username int, commentId int)(bool, error){
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM dislikeComment WHERE comment_id = ? AND user_id = ?)", commentId, username).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func CheckNumberOfDislikeForComment(db *sql.DB, commentId int) (int, error){
	var likeCount int
	err := db.QueryRow("SELECT COUNT(*) FROM dislikeComment WHERE comment_id = ?", commentId).Scan(&likeCount)
	if err != nil {
		return 0, err
	}

	return likeCount, nil
}

func DeleteDislikeCommentFromID(db *sql.DB, commentId int) error{
	_, err := db.Exec("DELETE FROM likeComment WHERE comment_id = ?", commentId)
	return err
}

func DeleteDislikeCommentFromDiscussion(db *sql.DB, discussionIDInt int) error {
	_, err := db.Exec("DELETE FROM dislikeComment WHERE discussion_id = ?", discussionIDInt)
	return err
}