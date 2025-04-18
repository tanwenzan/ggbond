package service

import (
	"github.com/gin-gonic/gin"
	"github.com/tanwenzan/blog-api/db"
	"github.com/tanwenzan/blog-api/jwt"
	"github.com/tanwenzan/blog-api/model"
	"gorm.io/gorm"
	"log"
)

func PostOwnerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		datasource := db.GetDatasource()
		userID := c.MustGet("userID").(uint)
		postID := c.Param("id")

		var post model.Post
		if result := datasource.Select("user_id").First(&post, postID); result.Error != nil {
			c.AbortWithStatusJSON(404, gin.H{"error": "文章不存在"})
			return
		}

		if post.UserID != userID {
			c.AbortWithStatusJSON(403, gin.H{"error": "无操作权限"})
			return
		}
		c.Next()
	}
}

func CreatePost(c *gin.Context) {
	datasource := db.GetDatasource()
	userID := c.MustGet("userID").(uint)

	var req struct {
		Title   string `json:"title" binding:"required,min=3"`
		Content string `json:"content" binding:"required,min=10"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "参数格式错误"})
		return
	}

	post := model.Post{
		Title:   req.Title,
		Content: req.Content,
		UserID:  userID,
	}

	if result := datasource.Create(&post); result.Error != nil {
		c.JSON(500, gin.H{"error": "创建失败"})
		return
	}

	c.JSON(201, gin.H{
		"id":      post.ID,
		"title":   post.Title,
		"content": post.Content,
		"created": post.CreatedAt,
	})
}

func UpdatePost(c *gin.Context) {
	datasource := db.GetDatasource()
	var req struct {
		Title   string `json:"title" binding:"omitempty,min=3"`
		Content string `json:"content" binding:"omitempty,min=10"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "参数格式错误"})
		return
	}

	postID := c.Param("id")
	updates := make(map[string]interface{})
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Content != "" {
		updates["content"] = req.Content
	}

	if result := datasource.Model(&model.Post{}).Where("id = ?", postID).Updates(updates); result.Error != nil {
		c.JSON(500, gin.H{"error": "更新失败"})
		return
	}

	c.JSON(200, gin.H{"message": "更新成功"})
}

func DeletePost(c *gin.Context) {
	datasource := db.GetDatasource()
	postID := c.Param("id")

	// 使用事务确保数据一致性
	err := datasource.Transaction(func(tx *gorm.DB) error {
		// 删除关联评论
		if err := tx.Where("post_id = ?", postID).Delete(&model.Comment{}).Error; err != nil {
			return err
		}
		// 删除文章
		return tx.Delete(&model.Post{}, postID).Error
	})

	if err != nil {
		c.JSON(500, gin.H{"error": "删除失败"})
		return
	}

	c.JSON(200, gin.H{"message": "删除成功"})
}

func GetPosts(c *gin.Context) {
	datasource := db.GetDatasource()
	user := jwt.GetCurrentUser(c)
	if user == nil {
		c.JSON(500, gin.H{"error": "当前用户为空，请登录"})
		return
	}
	var posts []model.Post
	datasource.Where("user_id = ?", user.ID).Find(&posts)
	// 转换响应结构
	var response []gin.H
	for _, post := range posts {
		response = append(response, gin.H{
			"id":            post.ID,
			"title":         post.Title,
			"user_id":       post.UserID,
			"created":       post.CreatedAt.Format("2006-01-02"),
			"comment_count": len(post.Comments),
		})
	}
	c.JSON(200, response)
}

func GetPost(c *gin.Context) {
	datasource := db.GetDatasource()
	postID := c.MustGet("id").(uint)
	var post model.Post
	result := datasource.Where("post_id = ?", postID).First(&post).Error
	if result.Error != nil {
		log.Println("[error] GetPost err:", result)
		c.JSON(500, gin.H{"error": "获取文章失败"})
	}
	c.JSON(200, gin.H{"post": post})
}
