package model

import "time"

type User struct {
	ID        uint       `db:"id"`
	Email     string     `db:"email"`
	Nickname  *string    `db:"nickname"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

type UserSocial struct {
	ID        uint       `db:"id"`
	UserID    uint       `db:"user_id"`
	Google    *string    `db:"google"`
	Kakao     *string    `db:"kakao"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

type UserAndSocial struct {
	User       User       `db:"user"`
	UserSocial UserSocial `db:"user_social"`
}
