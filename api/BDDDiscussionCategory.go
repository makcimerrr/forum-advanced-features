package api

import "database/sql"

func SetDiscussionCategory(db *sql.DB, discussionID int, categories int)(error){
	_, err := db.Exec("INSERT INTO discussion_category (discussion_id, category_id) VALUES (?, ?)", discussionID, categories)
	return err
}

func GetCategoryIDByDiscussionID(db *sql.DB, discussionID int) ([]int, error) {
    var categoryID []int

    rows, err := db.Query("SELECT category_id FROM discussion_category WHERE discussion_id = ?", discussionID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var id int
        if err := rows.Scan(&id); err != nil {
            return nil, err
        }
        categoryID = append(categoryID, id)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

    return categoryID, nil
}

func GetDiscussionIDByCategoryID(db *sql.DB, categoryID int)([]int, error){
	var discussionID []int

	rows, err := db.Query("SELECT discussion_id FROM discussion_category WHERE category_id = ?", categoryID)
	if err != nil {
        return nil, err
    }
    defer rows.Close()

	for rows.Next() {
        var id int
        if err := rows.Scan(&id); err != nil {
            return nil, err
        }
        discussionID = append(discussionID, id)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

	return  discussionID, err
}


func CheckIfDiscussionCategoryOk(db *sql.DB, category_id int, discussionID int) (bool, error) {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM discussion_category WHERE discussion_id = ? AND category_id = ?)", discussionID, category_id).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}