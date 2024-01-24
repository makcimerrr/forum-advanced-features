package api

import (
	"database/sql"
	"fmt"
	"os"
)

func CreateBDD() {
	// Chemin du fichier à vérifier
	nomFichier := "./website/data.db"

	// Vérifier si le fichier existe
	if _, err := os.Stat(nomFichier); os.IsNotExist(err) {
		// Le fichier n'existe pas, créons-le
		fichier, err := os.Create(nomFichier)
		if err != nil {
			fmt.Println("Erreur lors de la création du fichier :", err)
			return
		}
		defer fichier.Close()
		fmt.Println("Le fichier a été créé avec succès.")

	} else {
		// Le fichier existe déjà
		fmt.Println("Le fichier existe déjà.")
	}

	db, err := OpenBDD()
	if err != nil {
		return
	}

	// Création de la table s'il n'existe pas
	createTable := `
				CREATE TABLE IF NOT EXISTS "user" (
					"id"	INTEGER NOT NULL UNIQUE,
					"username"	TEXT NOT NULL UNIQUE,
					"email"	TEXT UNIQUE,
					"password"	TEXT,
					PRIMARY KEY("id" AUTOINCREMENT)
				)
				`
	_, err = db.Exec(createTable)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Création de la table s'il n'existe pas
	createTable = `
				CREATE TABLE IF NOT EXISTS "categories" (
					"id"	INTEGER NOT NULL UNIQUE,
					"category"	TEXT NOT NULL UNIQUE,
					PRIMARY KEY("id" AUTOINCREMENT)
					)
				`
	_, err = db.Exec(createTable)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Création de la table s'il n'existe pas
	createTable = `
					CREATE TABLE IF NOT EXISTS "discussion" (
						"id"	INTEGER NOT NULL UNIQUE,
						"user_id"	INTEGER NOT NULL,
						"title"	TEXT NOT NULL,
						"message"	TEXT NOT NULL,
						"bool_edit"	INTEGER NOT NULL DEFAULT 0,
						FOREIGN KEY("user_id") REFERENCES "user"("id"),
						PRIMARY KEY("id" AUTOINCREMENT)
					);
				`
	_, err = db.Exec(createTable)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Création de la table s'il n'existe pas
	createTable = `
				CREATE TABLE IF NOT EXISTS "discussion_category" (
					"id"	INTEGER NOT NULL UNIQUE,
					"discussion_id"	INTEGER NOT NULL,
					"category_id"	INTEGER NOT NULL,
					FOREIGN KEY("discussion_id") REFERENCES "discussion"("id"),
					PRIMARY KEY("id" AUTOINCREMENT),
					FOREIGN KEY("category_id") REFERENCES "categories"("id")
				)
				`
	_, err = db.Exec(createTable)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Création de la table s'il n'existe pas
	createTable = `
	CREATE TABLE IF NOT EXISTS "comment" (
		"id"	INTEGER NOT NULL UNIQUE,
		"user_id"	INTEGER NOT NULL,
		"discussion_id"	INTEGER NOT NULL,
		"message"	TEXT NOT NULL,
		"bool_edit"	INTEGER NOT NULL DEFAULT 0,
		FOREIGN KEY("user_id") REFERENCES "user"("id"),
		FOREIGN KEY("discussion_id") REFERENCES "discussion"("id"),
		PRIMARY KEY("id" AUTOINCREMENT)
	);
				`
	_, err = db.Exec(createTable)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Création de la table s'il n'existe pas
	createTable = `
				CREATE TABLE IF NOT EXISTS "likeDiscussion" (
					"id"	INTEGER NOT NULL UNIQUE,
					"discussion_id"	INTEGER NOT NULL,
					"user_id"	INTEGER NOT NULL,
					PRIMARY KEY("id" AUTOINCREMENT),
					FOREIGN KEY("discussion_id") REFERENCES "discussion"("id"),
					FOREIGN KEY("user_id") REFERENCES "user"("id")
				)
				`
	_, err = db.Exec(createTable)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Création de la table s'il n'existe pas
	createTable = `
				CREATE TABLE IF NOT EXISTS "dislikeDiscussion" (
					"id"	INTEGER NOT NULL UNIQUE,
					"discussion_id"	INTEGER NOT NULL,
					"user_id"	INTEGER NOT NULL,
					PRIMARY KEY("id" AUTOINCREMENT),
					FOREIGN KEY("discussion_id") REFERENCES "discussion"("id"),
					FOREIGN KEY("user_id") REFERENCES "user"("id")
				)
				`
	_, err = db.Exec(createTable)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Création de la table s'il n'existe pas
	createTable = `
				CREATE TABLE IF NOT EXISTS  "likeComment" (
					"id"	INTEGER NOT NULL UNIQUE,
					"comment_id"	INTEGER NOT NULL,
					"user_id"	INTEGER NOT NULL,
					PRIMARY KEY("id" AUTOINCREMENT),
					FOREIGN KEY("user_id") REFERENCES "user"("id"),
					FOREIGN KEY("comment_id") REFERENCES "comment"("id")
				)
				`
	_, err = db.Exec(createTable)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Création de la table s'il n'existe pas
	createTable = `
				CREATE TABLE IF NOT EXISTS  "dislikeComment" (
					"id"	INTEGER NOT NULL UNIQUE,
					"comment_id"	INTEGER NOT NULL,
					"user_id"	INTEGER NOT NULL,
					PRIMARY KEY("id" AUTOINCREMENT),
					FOREIGN KEY("user_id") REFERENCES "user"("id"),
					FOREIGN KEY("comment_id") REFERENCES "comment"("id")
				)
				`
	_, err = db.Exec(createTable)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Création de la table s'il n'existe pas
	createTable = `
				CREATE TABLE IF NOT EXISTS "session_user" (
					"id"	INTEGER NOT NULL UNIQUE,
					"username"	TEXT,
					"sessionToken"	TEXT,
					PRIMARY KEY("id" AUTOINCREMENT)
				)
				`
	_, err = db.Exec(createTable)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func OpenBDD() (*sql.DB, error) {
	// Ouvrez la connexion à la base de données
	db, err := sql.Open("sqlite", "./website/data.db")
	if err != nil {
		return nil, err
	}
	return db, err
}

func GetCategoryByName(db *sql.DB, category string) (int, error) {
	var id int
	err := db.QueryRow("SELECT id FROM categories WHERE category = ?", category).Scan(&id)
	return id, err
}

func GetCommentFromDiscussion(db *sql.DB) {

}
