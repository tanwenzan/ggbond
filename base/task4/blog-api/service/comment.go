package service

import (
	"github.com/gin-gonic/gin"
	"github.com/tanwenzan/blog-api/db"
	"github.com/tanwenzan/blog-api/model"
	"log"
	"strconv"
)

// 添加评论

func AddComment(c *gin.Context) {
	datasource := db.GetDatasource()
	userID := c.MustGet("userID").(uint)
	postIdStr := c.PostForm("post_id")
	postID, err := strconv.ParseUint(postIdStr, 10, 64)
	if err != nil {
		c.JSON(500, gin.H{"error": "非法的评论id：" + postIdStr})
		return
	}
	var comment struct {
		Content string `json:"content"`
	}

	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(400, gin.H{"error": "无效请求"})
		return
	}

	newComment := model.Comment{
		Content: comment.Content,
		UserID:  userID,
		PostID:  uint(postID),
	}

	if result := datasource.Create(&newComment); result.Error != nil {
		c.JSON(500, gin.H{"error": "评论失败"})
		return
	}

	c.JSON(201, newComment)
}

// 获取对应文章的评论

func GetComments(c *gin.Context) {
	datasource := db.GetDatasource()
	var comments []*model.Comment
	postID := c.MustGet("postID").(uint)
	err := datasource.Model(&comments).Where("post_id = ?", postID).Find(&comments).Error
	if err != nil {
		log.Println("[error] GetComments err:", err)
		c.JSON(500, gin.H{"error": "服务器错误"})
	}
	c.JSON(200, gin.H{"comments": comments})
}
