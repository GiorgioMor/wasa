package database

import "database/sql"

// SetName is an example that shows you how to execute insert/update
func (db *appdbimpl) CreateUser(username string) (uint64, error) {
	_, err := db.c.Exec("INSERT INTO USERS (username) VALUES (?)", username)

	if err != nil {
		return 0, err
	} else {
		var ret User

		err := db.c.QueryRow("SELECT ID, username  FROM USERS WHERE username = ?", username).Scan(&ret.ID, &ret.Username)
		if err != nil {
			return 0, err
		}

		return ret.ID, nil
	}
}

func (db *appdbimpl) DeleteUser(id uint64) error {
	rslt, err := db.c.Exec("DELETE FROM USERS WHERE ID = ?", id)
	if err != nil {
		return err
	}

	_, errr := rslt.RowsAffected()

	if errr != nil {
		return err
	}

	return nil
}

func (db *appdbimpl) GetUserByID(id uint64) (User, error) {
	var ret User

	err := db.c.QueryRow("SELECT ID, username FROM USERS WHERE id = ?", id).Scan(&ret.ID, &ret.Username)
	if err != nil {
		if err != sql.ErrNoRows {
			return User{}, err
		}
	}

	return ret, nil
}

func (db *appdbimpl) GetUser(username string) (User, error) {
	var ret User

	err := db.c.QueryRow("SELECT ID, username FROM USERS WHERE username = ?", username).Scan(&ret.ID, &ret.Username)
	if err != nil {
		if err != sql.ErrNoRows {
			return User{}, err
		}
	}

	return ret, nil
}

/*
func (db *appdbimpl) ListUsers() ([]User, error) {
	var ret []User

	rows, err := db.c.Query("SELECT ID, username  FROM USERS")
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		var u User
		err = rows.Scan(&u.ID, &u.Username)
		if err != nil {
			return nil, err
		}

		ret = append(ret, u)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return ret, nil
}
*/

func (db *appdbimpl) ChangeUsername(u User) error {
	var err error

	_, err = db.c.Exec(`UPDATE USERS SET username = ? WHERE ID = ?`, u.Username, u.ID)
	if err != nil {
		return err
	}

	_, err = db.c.Exec(`UPDATE COMMENTS SET username = ? WHERE user_id = ?`, u.Username, u.ID)
	if err != nil {
		return err
	}

	_, err = db.c.Exec(`UPDATE LIKES SET username = ? WHERE user_id = ?`, u.Username, u.ID)
	if err != nil {
		return err
	}

	return nil
}

func (db *appdbimpl) SearchUser(username string, user_id uint64) ([]User, error) {
	var ret []User

	rows, err := db.c.Query("SELECT * FROM USERS WHERE ID != ? AND username LIKE ? AND ID NOT IN (SELECT user_id FROM BANNED WHERE banned_user_id = ?)", user_id, username+"%", user_id)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		var user User
		err = rows.Scan(&user.ID, &user.Username)
		if err != nil {
			return nil, err
		}
		ret = append(ret, user)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return ret, nil
}
