package api

import (
	"database/sql"
)

func GetAllDiscussions(db *sql.DB) ([]Discussion, error) {
	// Exécutez une requête SQL pour récupérer toutes les discussions
	rows, err := db.Query("SELECT discussion.id, user.username AS users, discussion.title, discussion.message FROM discussion JOIN user ON discussion.user_id = user.id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Créez une slice pour stocker les discussions
	var discussions []Discussion

	// Parcourez les résultats et stockez-les dans la slice
	for rows.Next() {
		var discussion Discussion
		err := rows.Scan(&discussion.ID, &discussion.Username, &discussion.Title, &discussion.Message)
		if err != nil {
			return nil, err
		}
		discussions = append(discussions, discussion)
	}
	return discussions, nil
}

func GetOneDiscussions(db *sql.DB, discussionIDInt int) (Discussion, error) {
	var discussion Discussion
	err := db.QueryRow("SELECT discussion.id, user.username AS users, discussion.title, discussion.message FROM discussion JOIN user ON discussion.user_id = user.id WHERE discussion.id = ?", discussionIDInt).Scan(&discussion.ID, &discussion.Username, &discussion.Title, &discussion.Message)
	return discussion, err
}

func SetDiscussion(db *sql.DB, idUser int, title string, message string) error {
	_, err := db.Exec("INSERT INTO discussion (user_id, title, message) VALUES (?, ?, ?)", idUser, title, message)
	return err
}

func DeleteDiscussion(db *sql.DB, discussionIDInt int) error {
	_, err := db.Exec("DELETE FROM discussion WHERE id = ?", discussionIDInt)
	return err
}

func GetAllDiscussionsForOneUser(db *sql.DB, idUser int) ([]Discussion, error) {
	// Exécutez une requête SQL pour récupérer toutes les discussions
	rows, err := db.Query("SELECT discussion.id, user.username AS users, discussion.title, discussion.message FROM discussion JOIN user ON discussion.user_id = user.id WHERE discussion.user_id = ?",idUser)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Créez une slice pour stocker les discussions
	var discussions []Discussion

	// Parcourez les résultats et stockez-les dans la slice
	for rows.Next() {
		var discussion Discussion
		err := rows.Scan(&discussion.ID, &discussion.Username, &discussion.Title, &discussion.Message)
		if err != nil {
			return nil, err
		}
		discussions = append(discussions, discussion)
	}
	return discussions, nil
}

