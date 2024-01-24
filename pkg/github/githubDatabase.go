package github

import "database/sql"

// Function to check if the user already exists in the database
func userExistsInDatabase(db *sql.DB, user GitHubUser) (bool, error) {

	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM user WHERE username = ?", user.Login).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// Function to insert a new user into the database
func insertUserIntoDatabase(db *sql.DB, user GitHubUser) error {

	_, err := db.Exec("INSERT INTO user (username, email, password) VALUES (?, NULL, NULL)",user.Login)
	if err != nil {
		return err
	}

	return nil
}
