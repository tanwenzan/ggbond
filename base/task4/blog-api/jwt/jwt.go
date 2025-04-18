package jwt

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/tanwenzan/blog-api/db"
	"github.com/tanwenzan/blog-api/model"
	"time"
)

const secretKey = "zeroable"

// 生成JWT Token

func GenerateToken(userID uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})
	return token.SignedString([]byte(secretKey))
}

// JWT验证中间件

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "需要认证"})
			return
		}

		token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("userID", claims["user_id"])
			c.Next()
		} else {
			c.AbortWithStatusJSON(401, gin.H{"error": "无效Token"})
		}
	}
}

// 获取当前用户

func GetCurrentUser(c *gin.Context) *model.User {
	if value, exists := c.Get("userID"); exists {
		datasource := db.GetDatasource()
		var user model.User
		if err := datasource.Where("id = ?", value).First(&user).Error; err == nil {
			return &user
		}
	}
	return nil
}
