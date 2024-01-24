package api

import "database/sql"

func GetUserByUsername(db *sql.DB, name string) (int, error){
	var id int
	err := db.QueryRow("SELECT id FROM user WHERE username = ?", name).Scan(&id)
	return id, err
} 

func SetUser(db *sql.DB, user string, mail string, password uint32) error{
	_, err := db.Exec("INSERT INTO user (username, email, password) VALUES (?, ?, ?)", user, mail, password)
	return err
}