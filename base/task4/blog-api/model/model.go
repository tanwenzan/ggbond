package model

import "gorm.io/gorm"

// 用户模型

type User struct {
	gorm.Model
	Username string `gorm:"unique;not null"`
	Password string `gorm:"not null"` // 存储Bcrypt哈希值
	Email    string `gorm:"unique"`
	Posts    []Post `gorm:"foreignKey:UserID"`
}

// 文章模型

type Post struct {
	gorm.Model
	Title    string    `gorm:"size:255;not null"`
	Content  string    `gorm:"type:text;not null"`
	UserID   uint      `gorm:"index;not null"`
	Comments []Comment `gorm:"foreignKey:PostID"`
}

// 评论模型

type Comment struct {
	gorm.Model
	Content string `gorm:"type:text;not null"`
	UserID  uint   `gorm:"index;not null"`
	PostID  uint   `gorm:"index;not null"`
}
