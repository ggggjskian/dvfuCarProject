package handlers

import (
	"dvfucar/internal/config"
	"dvfucar/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type RegisterInput struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
	Name     string `json:"name" binding:"required"`
}

type LoginInput struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Все поля обязательны"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка шифрования пароля"})
		return
	}

	hashStr := string(hashedPassword)
	user := models.User{
		Phone:        &input.Phone,
		PasswordHash: &hashStr,
		Username:     input.Name,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Этот номер телефона уже зарегистрирован"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Регистрация успешна", "user": user})
}

func Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Укажите телефон и пароль"})
		return
	}

	var user models.User
	if err := config.DB.Where("phone = ?", input.Phone).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный телефон или пароль"})
		return
	}

	if user.PasswordHash == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Этот аккаунт зарегистрирован через Telegram"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*user.PasswordHash), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный телефон или пароль"})
		return
	}

	token := strconv.FormatUint(uint64(user.ID), 10)

	c.JSON(http.StatusOK, gin.H{
		"message": "Вход выполнен успешно",
		"token":   token,
		"user":    user,
	})
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Токен авторизации отсутствует"})
			return
		}

		const prefix = "Bearer "
		if len(authHeader) < len(prefix) || authHeader[:len(prefix)] != prefix {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Неверный формат авторизации"})
			return
		}

		token := authHeader[len(prefix):]

		userID, err := strconv.ParseUint(token, 10, 64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Недействительный токен"})
			return
		}

		var user models.User
		if err := config.DB.First(&user, userID).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не найден"})
			return
		}
		c.Set("userID", uint(userID))
		c.Next()
	}
}
