package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

func AdminAuthMiddleware(redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
			c.Abort()
			return
		}

		tokenParts := strings.Split(tokenString, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		tokenString = tokenParts[1]

		ctx := context.Background()
		isBlacklisted, err := redisClient.Exists(ctx, "blacklist:"+tokenString).Result()
		if err != nil {
			fmt.Println("Error checking Redis:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error while checking token"})
			c.Abort()
			return
		}
		if isBlacklisted > 0 {
			fmt.Println("Token is blacklisted:", tokenString)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Session expired. Please log in again."})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("invalid signing method")
			}
			return jwtSecret, nil
		})
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		role, ok := claims["role"].(string)
		if !ok || role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied: Only admins are allowed"})
			c.Abort()
			return
		}

		c.Set("admin_id", claims["user_id"])
		c.Next()
	}
}
