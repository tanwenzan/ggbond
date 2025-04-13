package main

import (
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ## 进阶gorm
// ### 题目1：模型定义
// - 假设你要开发一个博客系统，有以下几个实体： User （用户）、 Post （文章）、 Comment （评论）。
//   - 要求 ：
//   - 使用Gorm定义 User 、 Post 和 Comment 模型，其中 User 与 Post 是一对多关系（一个用户可以发布多篇文章）， Post 与 Comment 也是一对多关系（一篇文章可以有多个评论）。
//   - 编写Go代码，使用Gorm创建这些模型对应的数据库表。

// 用户模型

type User struct {
	gorm.Model        // 嵌入标准字段（ID, CreatedAt, UpdatedAt, DeletedAt）
	Name       string `gorm:"type:varchar(100);not null"`
	Email      string `gorm:"type:varchar(100);uniqueIndex"`
	Posts      []Post `gorm:"foreignKey:UserID"` // 一对多关系配置
}

// 文章模型

type Post struct {
	gorm.Model
	Title    string    `gorm:"type:varchar(200);index"`
	Content  string    `gorm:"type:text"`
	UserID   uint      // 外键字段（对应User模型的ID）
	Comments []Comment `gorm:"foreignKey:PostID"`   // 一对多关系配置
	Status   string    `gorm:"default:'有评论';index"` // 新增状态字段
}

// 评论模型

type Comment struct {
	gorm.Model
	Content string `gorm:"type:text"`
	PostID  uint   // 外键字段（对应Post模型的ID）
	UserID  uint   `gorm:"index"` // 可选的评论作者关联
}

//### 题目2：关联查询
//- 基于上述博客系统的模型定义。
//  - 要求 ：
//    - 编写Go代码，使用Gorm查询某个用户发布的所有文章及其对应的评论信息。
//    - 编写Go代码，使用Gorm查询评论数量最多的文章信息。

func GetUserPostsWithComments(userID uint) (User, error) {
	dsn := "user:password@tcp(localhost:3306)/blog?charset=utf8mb4&parseTime=True"
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	var user User

	// 嵌套预加载实现三级关联
	err := db.Preload("Posts.Comments").
		Where("id = ?", userID).
		First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return User{}, fmt.Errorf("用户不存在")
	} else if err != nil {
		return User{}, fmt.Errorf("数据库查询失败: %v", err)
	}

	return user, nil
}

func GetMostCommentedPost() (Post, error) {
	dsn := "user:password@tcp(localhost:3306)/blog?charset=utf8mb4&parseTime=True"
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	var post Post

	// 子查询统计评论数并排序
	err := db.Model(&Post{}).
		Select("posts.*, COUNT(comments.id) as comment_count").
		Joins("LEFT JOIN comments ON posts.id = comments.post_id").
		Group("posts.id").
		Order("comment_count DESC").
		Limit(1).
		Preload("Comments"). // 可选预加载评论详情
		First(&post).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return Post{}, fmt.Errorf("暂无文章数据")
	}
	if err != nil {
		return Post{}, fmt.Errorf("查询失败: %v", err)
	}

	return post, nil
}

//### 题目3：钩子函数
//- 继续使用博客系统的模型。
//  - 要求 ：
//    - 为 Post 模型添加一个钩子函数，在文章创建时自动更新用户的文章数量统计字段。
//    - 为 Comment 模型添加一个钩子函数，在评论删除时检查文章的评论数量，如果评论数量为 0，则更新文章的评论状态为 "无评论"。
// 为 Post 添加 AfterCreate 钩子

func (post *Post) AfterCreate(tx *gorm.DB) (err error) {
	// 原子操作更新用户文章计数
	result := tx.Model(&User{}).
		Where("id = ?", post.UserID).
		Update("posts_count", gorm.Expr("posts_count + 1"))

	if result.Error != nil {
		return fmt.Errorf("更新用户统计失败: %v", result.Error)
	}
	return nil
}

// 为 Comment 添加 AfterDelete 钩子

func (c *Comment) AfterDelete(tx *gorm.DB) (err error) {
	// 检查当前文章的评论总数
	var count int64
	tx.Model(&Comment{}).
		Where("post_id = ?", c.PostID).
		Count(&count)

	// 更新文章状态
	if count == 0 {
		tx.Model(&Post{}).
			Where("id = ?", c.PostID).
			Update("status", "无评论")
	}
	return nil
}
