package models

import (
	//Inbuild package(s)
	"log"
	"time"

	//Third-party package(s)
	"gorm.io/gorm"
)

// Login credentials
type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Signup credentials
type SignupReq struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

// Post Request
type PostReq struct {
	PostTitle   string `json:"post_title"`
	PostContent string `json:"post_content"`
	Catagory    string `json:"catagory"`
}

// User details
type User struct {
	UserId   uint   `json:"-" gorm:"primarykey"`
	Username string `json:"username" binding:"required" gorm:"column:username;type:varchar(100)"`
	Email    string `json:"email" binding:"required" gorm:"column:email;type:varchar(100) unique"`
	Password string `json:"password" binding:"required" gorm:"column:password;type:varchar(100)"`
	Role     string `json:"role" binding:"required" gorm:"-:all"`
	RoleId   uint   `json:"-" gorm:"column:role_id;type:bigint references Roles(role_id)"`
}

// Roles table
type Roles struct {
	RoleId uint   `gorm:"column:role_id;type:bigint primary key"`
	Role   string `gorm:"column:role;type:varchar(50)"`
}

// Catagory Table
type Catagory struct {
	CatagoryId uint   `gorm:"column:catagory_id;type:bigint primary key"`
	Catagory   string `gorm:"column:catagory;type:varchar(50)"`
}

// Token values for each user-id
type Authentication struct {
	UserId uint   `json:"user_id" gorm:"column:user_id;type:bigint primary Key"`
	Token  string `json:"token" gorm:"column:token;type:varchar(200)"`
}

// Post details
type Post struct {
	PostId      uint           `json:"-" gorm:"primarykey"`
	PostTitle   string         `json:"post_title,omitempty" binding:"required" gorm:"column:post_title;type:varchar(100) unique"`
	PostContent string         `json:"post_content,omitempty" binding:"required" gorm:"column:post_content;type:varchar(500)"`
	Catagory    string         `json:"catagory,omitempty" binding:"required" gorm:"-"`
	CatagoryId  uint           `json:"-" gorm:"column:catagory_id;type:bigint references Catagories(catagory_id)"`
	UserId      uint           `json:"-" gorm:"column:user_id;type:bigint references Users(user_id)"`
	CreatedAt   time.Time      `json:"-" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"-" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// Comments
type Comments struct {
	CommentId uint           `json:"-" gorm:"primarykey"`
	Comment   string         `json:"comment" binding:"required" gorm:"column:comment;type:varchar(200)"`
	UserId    uint           `json:"-" gorm:"column:user_id;type:bigint references Users(user_id)"`
	PostId    uint           `json:"-" gorm:"column:post_id;type:bigint references Posts(post_id)"`
	PostTitle string         `json:"post_title,omitempty" binding:"required" gorm:"-"`
	CreatedAt time.Time      `json:"-" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"-" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// Comments request
type CommentReq struct {
	PostTitle string `json:"post_title"`
	Comment   string `json:"comment"`
}

// Custom Log
type Logs struct {
	Info  *log.Logger
	Error *log.Logger
}
