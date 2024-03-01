package database

type User struct {
	ID       uint64
	Username string
}

type Post struct {
	ID               uint64
	Created_Datetime string
	Media_file       string
	User_ID          uint64
	Caption          string
	Comments         []Comment
	Likes            []Like
}

type PostID struct {
	Post_ID uint64 `json:"post_id"`
}

type Comment struct {
	ID       uint64
	Text     string
	Username string
	User_ID  uint64
	Post_ID  uint64
}

type Like struct {
	ID       uint64
	User_ID  uint64
	Username string
	Post_ID  uint64
}
