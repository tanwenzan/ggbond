package service

import (
	"github.com/gin-gonic/gin"
	"github.com/tanwenzan/blog-api/db"
	"github.com/tanwenzan/blog-api/jwt"
	"github.com/tanwenzan/blog-api/model"
	"golang.org/x/crypto/bcrypt"
	"log"
)

// 用户注册

func Register(c *gin.Context) {
	datasource := db.GetDatasource()
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": "无效请求"})
		return
	}

	// 密码加密
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)

	if result := datasource.Create(&user); result.Error != nil {
		c.JSON(500, gin.H{"error": "创建用户失败"})
		return
	}

	c.JSON(201, gin.H{"id": user.ID})
}

// 用户登录

func Login(c *gin.Context) {
	datasource := db.GetDatasource()
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": "参数有误"})
		return
	}
	dbUser := model.User{}
	err := datasource.First(&dbUser, user.ID).Error
	if err != nil || len(dbUser.Username) == 0 {
		c.JSON(400, gin.H{"error": "用户不存在"})
	}
	if user.Password != dbUser.Password {
		c.JSON(500, gin.H{"error": "密码错误"})
		return
	}
	token, err := jwt.GenerateToken(dbUser.ID)
	if err != nil {
		log.Print("[error] GenerateToken err: ", err)
		c.JSON(500, gin.H{"error": "系统错误"})
		return
	}
	c.JSON(200, gin.H{"token": token})
}
