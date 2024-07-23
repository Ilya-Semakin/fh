package controllers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/Ilya-Semakin/fh/config"
	middelwares "github.com/Ilya-Semakin/fh/middlewares"
	"github.com/Ilya-Semakin/fh/models"
	"github.com/Ilya-Semakin/fh/services"
	"github.com/Ilya-Semakin/fh/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

func Register(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user := models.User{Username: input.Username, Email: input.Email, Password: hashedPassword}
	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Registration successful"})
}

func Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	if !utils.CheckPasswordHash(input.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	token, err := middelwares.GenerateJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func OAuthVKLogin(c *gin.Context) {
	url := services.VKConfig.AuthCodeURL("state")
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func OAuthVKCallback(c *gin.Context) {
	code := c.Query("code")
	token, err := services.VKConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Println("Code exchange failed:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Code exchange failed"})
		return
	}

	response, err := http.Get("https://api.vk.com/method/users.get?access_token=" + token.AccessToken + "&v=5.130&fields=uid,first_name,last_name")
	if err != nil {
		log.Println("Failed getting user info:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed getting user info"})
		return
	}
	defer response.Body.Close()

	data, _ := ioutil.ReadAll(response.Body)
	var vkResponse map[string]interface{}
	json.Unmarshal(data, &vkResponse)

	userData := vkResponse["response"].([]interface{})[0].(map[string]interface{})
	vkID := userData["id"].(float64)
	email := userData["email"].(string) // VK API может не вернуть email, его нужно запросить дополнительно

	var user models.User
	if err := config.DB.Where("vk_id = ?", vkID).First(&user).Error; err != nil {
		user = models.User{
			VKID:     vkID,
			Email:    email,
			Username: userData["first_name"].(string) + " " + userData["last_name"].(string),
		}
		config.DB.Create(&user)
	}

	tokenString, err := middelwares.GenerateJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func OAuthGoogleLogin(c *gin.Context) {
	url := services.GoogleConfig.AuthCodeURL("state")
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func OAuthGoogleCallback(c *gin.Context) {
	code := c.Query("code")
	token, err := services.GoogleConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Println("Code exchange failed:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Code exchange failed"})
		return
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		log.Println("Failed getting user info:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed getting user info"})
		return
	}
	defer response.Body.Close()

	data, _ := ioutil.ReadAll(response.Body)
	var userInfo map[string]interface{}
	json.Unmarshal(data, &userInfo)

	googleID := userInfo["id"].(string)
	email := userInfo["email"].(string)

	var user models.User
	if err := config.DB.Where("google_id = ?", googleID).First(&user).Error; err != nil {
		user = models.User{
			GoogleID: googleID,
			Email:    email,
			Username: userInfo["name"].(string),
		}
		config.DB.Create(&user)
	}

	tokenString, err := middelwares.GenerateJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
