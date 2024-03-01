package database

// Take the list of users that liked a post
func (db *appdbimpl) GetLikesList(id uint64) ([]Like, error) {

	rows, err := db.c.Query("SELECT * FROM LIKES WHERE post_id = ?", id)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var likes []Like
	for rows.Next() {
		var l Like
		err = rows.Scan(&l.ID, &l.Post_ID, &l.User_ID, &l.Username)
		if err != nil {
			return nil, err
		}

		likes = append(likes, l)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return likes, nil
}

// Add a like of a user to a photo
func (db *appdbimpl) LikePost(Post_ID uint64, u User) error {

	_, err := db.c.Exec("INSERT INTO LIKES (post_id, user_id, username) VALUES (?, ?, ?)", Post_ID, u.ID, u.Username)
	if err != nil {
		return err
	}

	return nil
}

// Removes a like of a user from a photo
func (db *appdbimpl) UnlikePost(Post_ID uint64, u User) error {

	_, err := db.c.Exec("DELETE FROM LIKES WHERE(post_id = ? AND user_id = ?)", Post_ID, u.ID)
	if err != nil {
		return err
	}

	return nil
}
