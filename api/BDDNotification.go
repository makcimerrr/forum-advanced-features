package api

import (
	"database/sql"
	"fmt"
)

func GetNotificationByIdUserAndVu(db *sql.DB, userIDCreateur int)([]Notification, error){
	// Exécutez une requête SQL pour récupérer toutes les discussions
	rows, err := db.Query("SELECT * FROM notification WHERE userIDCreateur = ? AND user_id != ? AND vu = 0", userIDCreateur, userIDCreateur)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Créez une slice pour stocker les discussions
	var notifications []Notification

	// Parcourez les résultats et stockez-les dans la slice
	for rows.Next() {
		var notification Notification
		err := rows.Scan(&notification.ID, &notification.UserIDCreateur, &notification.User_id, &notification.Discussion_id, &notification.Message, &notification.vu)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, notification)
	}
	return notifications, nil
}

func GetNumberNotificationById(db *sql.DB, userIDCreateur int) (int, error) {
	var notificationCount int
	err := db.QueryRow("SELECT COUNT(*) FROM notification WHERE userIDCreateur = ? And user_id != ? AND vu = 0", userIDCreateur, userIDCreateur).Scan(&notificationCount)
	if err != nil {
		return 0, err
	}

	return notificationCount, nil
}

func SetNotification(db *sql.DB, userIDCreateur int, userId int, discussionId int, message string) error {
	_, err := db.Exec("INSERT INTO notification (userIDCreateur, discussion_id, user_id, message) VALUES (?, ?, ?, ?)",userIDCreateur, discussionId, userId, message)

	return err
}

func UpdateNotificationTrue(db *sql.DB, id int) error {
	// Update statement
	updateQuery := "UPDATE notification SET vu = 1 WHERE id = ?"

	// Prepare the statement
	stmt, err := db.Prepare(updateQuery)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the update
	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}
	return err
}

func CheckIfNotificationNotDouble(db *sql.DB, userId int, discussionId int, message string) (bool, error) {
    var notif bool
    err := db.QueryRow("SELECT 1 FROM notification WHERE user_id = ? AND discussion_id = ? AND message = ?", userId, discussionId, message).Scan(&notif)
	fmt.Println("test", err)
    return notif, err
}

func DeleteNotificationFromDiscussion(db *sql.DB, discussionId int) error{
	_, err := db.Exec("DELETE FROM notification WHERE discussion_id = ?", discussionId)
	return err
}