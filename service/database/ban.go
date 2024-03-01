package database

func (db *appdbimpl) BanUser(Ban_ID uint64, User_ID uint64) error {
	var err error
	var posts []Post

	_, err = db.c.Exec("INSERT INTO BANNED (banned_user_id, user_id) VALUES (?, ?)", Ban_ID, User_ID)
	if err != nil {
		return err
	}

	user_id, err := db.GetUserByID(User_ID)
	if err != nil {
		return err
	}
	banned_user_id, err := db.GetUserByID(Ban_ID)
	if err != nil {
		return err
	}

	// Prendi i tuoi post ed elimina i mi piace e i commenti che l'utente che hai bannato aveva fatto
	posts, err = db.GetPosts(user_id.ID)
	if err != nil {
		return err
	}
	for _, p := range posts {
		for _, c := range p.Comments {
			if c.User_ID == Ban_ID {
				err = db.RemoveComment(c.ID)
				if err != nil {
					return err
				}
			}
		}
		for _, l := range p.Likes {
			if l.User_ID == Ban_ID {
				err = db.UnlikePost(p.ID, banned_user_id)
				if err != nil {
					return err
				}
			}
		}
	}

	// Prendi i post dell'utente che stai bannando ed elimina i mi piace e i commenti dell'utente che lo ha bannato
	posts, err = db.GetPosts(banned_user_id.ID)
	if err != nil {
		return err
	}
	for _, p := range posts {
		for _, c := range p.Comments {
			if c.User_ID == User_ID {
				err = db.RemoveComment(c.ID)
				if err != nil {
					return err
				}
			}
		}
		for _, l := range p.Likes {
			if l.User_ID == User_ID {
				err = db.UnlikePost(p.ID, user_id)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (db *appdbimpl) UnbanUser(Ban_ID uint64, User_ID uint64) error {

	_, err := db.c.Exec("DELETE FROM BANNED WHERE(banned_user_id = ? AND user_id = ?)", Ban_ID, User_ID)
	if err != nil {
		return err
	}

	return nil
}

func (db *appdbimpl) IsBanned(UserAsking uint64, UserToCheck uint64) (bool, error) {

	var ret int
	err := db.c.QueryRow("SELECT COUNT(*) FROM BANNED WHERE user_id = ? AND banned_user_id = ?", UserToCheck, UserAsking).Scan(&ret)
	if err != nil {
		return false, err
	}

	if ret > 0 {
		return true, nil
	} else {
		return false, nil
	}
}
