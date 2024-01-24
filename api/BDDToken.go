package api

import "database/sql"

func GetTokenByUser(db *sql.DB, username string)(string, error){
	var dbSessionToken string
	err := db.QueryRow("SELECT sessionToken FROM session_user WHERE username = ?", username).Scan(&dbSessionToken)
	return dbSessionToken, err
}

func SetTokenSession(db *sql.DB, username string, sessionToken string) (error){
	_, err := db.Exec("INSERT INTO session_user (username, sessionToken) VALUES (?, ?)", username, sessionToken)
	return err
}

func DeleteTokenSession(db *sql.DB, username string) (error){
	_, err := db.Exec("DELETE FROM session_user WHERE username = ? ", username)
	return err
}