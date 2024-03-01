package database

import "database/sql"

// Take the list of comments of a photo
func (db *appdbimpl) GetCommentsList(id uint64) ([]Comment, error) {

	rows, err := db.c.Query("SELECT * FROM COMMENTS WHERE post_id = ?", id)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var comments []Comment
	for rows.Next() {
		var comment Comment
		err = rows.Scan(&comment.ID, &comment.Post_ID, &comment.User_ID, &comment.Username, &comment.Text)
		if err != nil {
			return nil, err
		}

		comments = append(comments, comment)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return comments, nil
}

// Add a comment to a photo
func (db *appdbimpl) AddComment(id uint64, u User, comment string) error {

	_, err := db.c.Exec("INSERT INTO COMMENTS (post_id, user_id, username, comment) VALUES (?, ?, ?, ?)", id, u.ID, u.Username, comment)
	if err != nil {
		return err
	}

	return nil
}

// Database function that removes a comment of a user from a photo
func (db *appdbimpl) RemoveComment(id uint64) error {

	_, err := db.c.Exec("DELETE FROM COMMENTS WHERE ID = ?", id)
	if err != nil {
		return err
	}

	return nil
}

func (db *appdbimpl) GetCommentByID(id uint64) (Comment, error) {
	var ret Comment

	err := db.c.QueryRow("SELECT * FROM COMMENTS WHERE ID = ?", id).Scan(&ret.ID, &ret.Post_ID, &ret.User_ID, &ret.Username, &ret.Text)
	if err != nil {
		if err != sql.ErrNoRows {
			return Comment{}, err
		}
	}

	return ret, nil
}
