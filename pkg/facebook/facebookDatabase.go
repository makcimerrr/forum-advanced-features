package facebook

import "database/sql"

func userExistsInDatabase(db *sql.DB, user FacebookUser) (bool, error) {

	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM user WHERE id = ?", user.ID).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func insertUserIntoDatabase(db *sql.DB, user FacebookUser) error {
	_, err := db.Exec("INSERT INTO user (id, username, email, password) VALUES (?, ?, NULL, NULL)", user.ID, user.Name)
	if err != nil {
		return err
	}

	return nil
}
