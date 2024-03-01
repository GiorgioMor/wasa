/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"
)

var ErrUserDoesNotExist = errors.New("user does not exist")

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	// Get a single user filtered by the username
	GetUser(username string) (User, error)

	// Create a new user
	CreateUser(username string) (uint64, error)

	// Get a single user filtered by the id
	GetUserByID(id uint64) (User, error)

	// Delete a single user filtered by the id
	DeleteUser(id uint64) error

	// Create a new post
	CreatePost(post Post) error

	// Get list of follower
	GetFollower(id uint64) ([]User, error)

	// Get list of follower
	GetFollowing(id uint64) ([]User, error)

	// Get user's posts
	GetPosts(id uint64) ([]Post, error)

	// Put a follow
	FollowUser(User_ID uint64, Follower_ID uint64) error

	// Remove a follow
	UnfollowUser(User_ID uint64, Follower_ID uint64) error

	// Put a ban
	BanUser(User_ID uint64, Ban_ID uint64) error

	// Remove a ban
	UnbanUser(User_ID uint64, Ban_ID uint64) error

	// Search User
	SearchUser(Username string, user_id uint64) ([]User, error)

	// Change Username
	ChangeUsername(User User) error

	// Remove Post
	RemovePost(Post_ID uint64, User_ID uint64) error

	// Get a single post by ID
	GetPost(Post_ID uint64) (Post, error)

	// Check if a user is banned
	IsBanned(UserAsking uint64, UserToCheck uint64) (bool, error)

	// Like a post
	LikePost(Post_ID uint64, u User) error

	// Remove a like from a post
	UnlikePost(Post_ID uint64, u User) error

	// Add a comment on a post
	AddComment(id uint64, u User, comment string) error

	// Remove a comment from a post
	RemoveComment(id uint64) error

	// Get a single comment by ID
	GetCommentByID(id uint64) (Comment, error)

	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	// Check if table exists. If not, the database is empty, and we need to create the structure
	var tableName string
	var err error

	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='USERS';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE USERS (ID INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, username TEXT NOT NULL);`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure (table users): %w", err)
		}
	}

	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='FOLLOWER';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE FOLLOWER (following_user_id INTEGER, follower_user_id INTEGER, FOREIGN KEY (following_user_id) REFERENCES USERS(ID), FOREIGN KEY (follower_user_id) REFERENCES USERS(ID));`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure (table follower): %w", err)
		}
	}

	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='BANNED';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE BANNED (banned_user_id INTEGER, user_id INTEGER, FOREIGN KEY (banned_user_id) REFERENCES USERS(ID), FOREIGN KEY (user_id) REFERENCES USERS(ID));`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure (table banned): %w", err)
		}
	}

	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='POST';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE POST (ID INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, path_photo TEXT NOT NULL, created_at TEXT NOT NULL, caption TEXT, user_id INTEGER, FOREIGN KEY (user_id) REFERENCES USERS(ID));`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure (table post): %w", err)
		}
	}

	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='LIKES';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE LIKES (ID INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, post_id INTEGER, user_id INTEGER, username TEXT, FOREIGN KEY (user_id) REFERENCES USERS(ID), FOREIGN KEY (post_id) REFERENCES POST(ID));`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure (table like): %w", err)
		}
	}

	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='COMMENTS';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE COMMENTS (ID INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, post_id INTEGER, user_id INTEGER, username TEXT, comment TEXT, FOREIGN KEY (user_id) REFERENCES USERS(ID), FOREIGN KEY (post_id) REFERENCES POST(ID));`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure (table comments): %w", err)
		}
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
