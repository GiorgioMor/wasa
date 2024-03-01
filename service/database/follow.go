package database

func (db *appdbimpl) FollowUser(User_ID uint64, Follower_ID uint64) error {

	_, err := db.c.Exec("INSERT INTO FOLLOWER (following_user_id, follower_user_id) VALUES (?, ?)", User_ID, Follower_ID)
	if err != nil {
		return err
	}

	return nil
}

func (db *appdbimpl) UnfollowUser(User_ID uint64, Follower_ID uint64) error {

	_, err := db.c.Exec("DELETE FROM FOLLOWER WHERE(following_user_id = ? AND follower_user_id = ?)", User_ID, Follower_ID)
	if err != nil {
		return err
	}

	return nil
}

func (db *appdbimpl) GetFollower(id uint64) ([]User, error) {
	var ret []User

	rows, err := db.c.Query("SELECT USERS.ID, USERS.username FROM USERS, FOLLOWER WHERE FOLLOWER.following_user_id = ? AND USERS.ID = FOLLOWER.follower_user_id", id)
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

func (db *appdbimpl) GetFollowing(id uint64) ([]User, error) {
	var ret []User

	rows, err := db.c.Query("SELECT USERS.ID, USERS.username FROM USERS, FOLLOWER WHERE FOLLOWER.follower_user_id = ? AND USERS.ID = FOLLOWER.following_user_id", id)
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
