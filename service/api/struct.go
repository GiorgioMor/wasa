package api

import "gitlab.com/Protinus/homework/service/database"

type User struct {
	ID       uint64 `json:"id"`
	Username string `json:"username"`
}

type Post struct {
	ID               uint64 `json:"id"`
	Created_Datetime string `json:"create_datetime"`
	Media_file       string `json:"media_file"`
	User_ID          uint64 `json:"user_id"`
	Caption          string `json:"caption"`
}

type Likes struct {
	ID      uint64 `json:"id"`
	User_ID uint64 `json:"user_id"`
	Post_ID uint64 `json:"post_id"`
}

type Comment struct {
	ID      uint64 `json:"id"`
	Text    string `json:"text"`
	User_ID uint64 `json:"user_id"`
	Post_ID uint64 `json:"post_id"`
}

type Follower struct {
	Username         string `json:"follower_username"`
	Follower_User_ID uint64 `json:"follower_user_id"`
}

type Following struct {
	Username          string `json:"following_username"`
	Following_User_ID uint64 `json:"following_user_id"`
}

type Banned struct {
	User_ID        uint64 `json:"user_id"`
	Banned_User_ID uint64 `json:"banned_user_id"`
}

type JSONErrorMsg struct {
	Message string
}

type Username struct {
	Username string `json:"username"`
}

type CommentText struct {
	Comment string `json:"comment"`
}

type SendPost struct {
	File     string `json:"file"`
	Caption  string `json:"caption"`
	FileName string `json:"fileName"`
}

type UserProfile struct {
	User_ID   uint64          `json:"user_id"`
	Username  string          `json:"username"`
	Followers []database.User `json:"followers"`
	Following []database.User `json:"following"`
	Posts     []database.Post `json:"posts"`
	IsBanned  bool            `json:"isBanned"`
}

type UserHome struct {
	Followers []database.User `json:"followers"`
	Posts     []database.Post `json:"posts"`
}

/*
func (u *User) FromDatabase(user database.User) {
	u.ID = user.ID
	u.Username = user.Username
}
*/

// ToDatabase returns the user in a database-compatible representation
func (u User) ToDatabase() database.User {
	return database.User{
		ID:       u.ID,
		Username: u.Username,
	}
}
