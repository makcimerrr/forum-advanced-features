package api

import (
	"database/sql"
)

//function for table like Discussion

func GetLikesFromOneDiscussion(db *sql.DB, discussionIDInt int, username int) (bool, error) {
	var liked bool
	err := db.QueryRow("SELECT 1 FROM likeDiscussion WHERE discussion_id = ? AND user_id = ?", discussionIDInt, username).Scan(&liked)
	return liked, err
}

func SetLikesDiscussion(db *sql.DB, discussionIDInt int, username int) error {
	_, err := db.Exec("INSERT INTO likeDiscussion (discussion_id, user_id) VALUES (?, ?)", discussionIDInt, username)
	return err
}

func DeleteLikeByUserIdForDiscussion(db *sql.DB, discussionIDInt int, username int) error {
	_, err := db.Exec("DELETE FROM likeDiscussion WHERE discussion_id = ? AND user_id = ?", discussionIDInt, username)
	return err
}

func CheckIfUserLikedDiscussion(db *sql.DB, idUser int, discussionID int) (bool, error) {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM likeDiscussion WHERE discussion_id = ? AND user_id = ?)", discussionID, idUser).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func CheckNumberOfLikesForDiscussion(db *sql.DB, discussionID int) (int, error) {

	var likeCount int
	err := db.QueryRow("SELECT COUNT(*) FROM likeDiscussion WHERE discussion_id = ?", discussionID).Scan(&likeCount)
	if err != nil {
		return 0, err
	}

	return likeCount, nil
}

func GetDiscussionIdByLikeForOneUser(db *sql.DB, idUser int) ([]int, error) {
	var tempId []int
	rows, err := db.Query("SELECT discussion_id FROM likeDiscussion WHERE user_id = ?", idUser)
	if err != nil {
        return nil, err
    }
    defer rows.Close()

	for rows.Next() {
        var id int
        if err := rows.Scan(&id); err != nil {
            return nil, err
        }
        tempId = append(tempId, id)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

	return  tempId, err
}

func DeleteLikeFromDiscussion(db *sql.DB, discussionIDInt int) error {
	_, err := db.Exec("DELETE FROM likeDiscussion WHERE discussion_id = ?", discussionIDInt)
	return err
}



//function for table like Comment

func GetLikesFromOneComment(db *sql.DB, commentId int, username int)(bool, error){
	var liked bool
	err := db.QueryRow("SELECT 1 FROM likeComment WHERE comment_id = ? AND user_id = ?", commentId, username).Scan(&liked)
	return liked, err
}

func SetLikesComment(db *sql.DB, discussionID int, commentId int, username int) (error){
	_, err := db.Exec("INSERT INTO likeComment (discussion_id, comment_id, user_id) VALUES (?, ?, ?)", discussionID, commentId, username)
	return err
}

func DeleteLikeByUserIdForComment(db *sql.DB, commentId int, username int) (error) {
	_, err := db.Exec("DELETE FROM likeComment WHERE comment_id = ? AND user_id = ?", commentId, username)
	return err
}

func CheckIfUserLikedComment(db *sql.DB, username int, commentid int)(bool, error){
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM likeComment WHERE comment_id = ? AND user_id = ?)", commentid, username).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func CheckNumberOfLikesForComment(db *sql.DB, commentId int) (int, error){
	var likeCount int
	err := db.QueryRow("SELECT COUNT(*) FROM likeComment WHERE comment_id = ?", commentId).Scan(&likeCount)
	if err != nil {
		return 0, err
	}

	return likeCount, nil
}

func GetCommentIdByLikeForOneUser(db *sql.DB, idUser int) ([]int, error) {
	var tempId []int
	rows, err := db.Query("SELECT comment_id FROM likeComment WHERE user_id = ?", idUser)
	if err != nil {
        return nil, err
    }
    defer rows.Close()

	for rows.Next() {
        var id int
        if err := rows.Scan(&id); err != nil {
            return nil, err
        }
        tempId = append(tempId, id)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

	return  tempId, err
}

func DeleteLikeCommentFromID(db *sql.DB, commentId int) error{
	_, err := db.Exec("DELETE FROM likeComment WHERE comment_id = ?", commentId)
	return err
}

func DeleteLikeCommentFromDiscussion(db *sql.DB, discussionIDInt int) error {
	_, err := db.Exec("DELETE FROM likeComment WHERE discussion_id = ?", discussionIDInt)
	return err
}