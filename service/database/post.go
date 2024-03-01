package database

import "database/sql"

func (db *appdbimpl) CreatePost(post Post) error {
	path := post.Media_file

	_, err := db.c.Exec("INSERT INTO POST (path_photo, created_at, caption, user_id) VALUES (?, ?, ?, ?)", path, post.Created_Datetime, post.Caption, post.User_ID)

	if err != nil {
		return err
	} else {
		return nil
	}
}

func (db *appdbimpl) GetPost(id uint64) (Post, error) {
	var ret Post

	err := db.c.QueryRow("SELECT * FROM POST WHERE ID = ?", id).Scan(&ret.ID, &ret.Media_file, &ret.Created_Datetime, &ret.Caption, &ret.User_ID)
	if err != nil {
		if err != sql.ErrNoRows {
			return Post{}, err
		}
	}

	return ret, nil
}

func (db *appdbimpl) GetPosts(id uint64) ([]Post, error) {
	var ret []Post

	rows, err := db.c.Query("SELECT * FROM POST WHERE user_id = ? ORDER BY created_at DESC", id)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		var p Post
		err = rows.Scan(&p.ID, &p.Media_file, &p.Created_Datetime, &p.Caption, &p.User_ID)
		if err != nil {
			return nil, err
		}

		comments, err := db.GetCommentsList(p.ID)
		if err != nil {
			return nil, err
		}
		p.Comments = comments

		likes, err := db.GetLikesList(p.ID)
		if err != nil {
			return nil, err
		}
		p.Likes = likes

		ret = append(ret, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return ret, nil
}

func (db *appdbimpl) RemovePost(User_ID uint64, Post_ID uint64) error {

	var _ sql.Result
	var err error

	_, err = db.c.Exec("DELETE FROM COMMENTS WHERE post_id = ?", Post_ID)
	if err != nil {
		return err
	}

	_, err = db.c.Exec("DELETE FROM LIKES WHERE post_id = ?", Post_ID)
	if err != nil {
		return err
	}

	_, err = db.c.Exec("DELETE FROM POST WHERE ID = ? AND user_id = ?", Post_ID, User_ID)
	if err != nil {
		return err
	}

	return nil
}
