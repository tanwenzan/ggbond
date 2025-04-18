package main

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/tanwenzan/blog-api/db"
	"github.com/tanwenzan/blog-api/jwt"
	"github.com/tanwenzan/blog-api/service"
	"gorm.io/gorm"
	"log"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			switch {
			case errors.Is(err, gorm.ErrRecordNotFound):
				c.JSON(404, gin.H{"error": "资源不存在"})
			default:
				c.JSON(500, gin.H{"error": "服务器错误"})
			}
		}
	}
}

func main() {
	// 初始化数据库连接
	dbErr := db.InitDatasource()
	if dbErr != nil {
		log.Fatal(dbErr)
		return
	}
	r := gin.Default()
	// 公共路由
	r.POST("/register", service.Register)
	r.POST("/login", service.Login)
	// 需要认证的路由组
	auth := r.Group("/")
	auth.Use(jwt.AuthMiddleware())
	auth.Use(ErrorHandler())
	{
		auth.POST("/posts", service.CreatePost)
		auth.PUT("/posts/:id", service.PostOwnerMiddleware(), service.UpdatePost)
		auth.DELETE("/posts/:id", service.PostOwnerMiddleware(), service.DeletePost)
		auth.POST("/posts/:postID/comments", service.AddComment)
	}

	// 公共访问路由
	r.GET("/posts", service.GetPosts)
	r.GET("/posts/:id", service.GetPost)
	r.GET("/posts/:postID/comments", service.GetComments)
	err := r.Run(":8080")
	if err != nil {
		return
	}
}
