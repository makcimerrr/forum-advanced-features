package api

import "database/sql"


//func SetCategory()

func GetCategory(db *sql.DB) ([]Categories, error) {

	// Créez une slice pour stocker les discussions
	rows, err := db.Query("SELECT id, category FROM categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []Categories

	// Parcourez les résultats et stockez-les dans la slice
	for rows.Next() {
		var category Categories
		err := rows.Scan(&category.ID, &category.Category)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func GetCategoryByID(db *sql.DB, id int)(string, error){
	var category string
	err := db.QueryRow("SELECT category FROM categories WHERE id = ?", id).Scan(&category)
	return  category, err
}

func GetIDCategoryByCategory(db *sql.DB, category string)(int, error){
	var id int
	err := db.QueryRow("SELECT id FROM categories WHERE category = ?", category).Scan(&id)
	return  id, err
}