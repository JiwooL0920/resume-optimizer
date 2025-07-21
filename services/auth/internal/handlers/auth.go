package handlers

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/resume-optimizer/auth-service/internal/services"
	"github.com/resume-optimizer/auth-service/internal/utils"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	userService = services.NewUserService()
	jwtSecret   = os.Getenv("JWT_SECRET")
	oauth2Config = &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Endpoint:     google.Endpoint,
		RedirectURL:  os.Getenv("REDIRECT_URL"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
	}
)

// GoogleAuth initiates the OAuth2 flow for Google
func GoogleAuth(c *gin.Context) {
	url := oauth2Config.AuthCodeURL("state", oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// GoogleCallback handles the OAuth2 callback from Google
func GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	token, err := oauth2Config.Exchange(c, code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error exchanging code: " + err.Error()})
		return
	}

	client := oauth2Config.Client(c, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting user info: " + err.Error()})
		return
	}
	defer resp.Body.Close()

	var userInfo map[string]interface{}
	if err = json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error parsing user info: " + err.Error()})
		return
	}

	email := userInfo["email"].(string)
	name := userInfo["name"].(string)
	picture := userInfo["picture"].(string)
	googleID := userInfo["id"].(string)

	user, err := userService.GetUserByEmail(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error accessing user service: " + err.Error()})
		return
	}

	if user == nil {
		user, err = userService.CreateUser(email, name, &googleID, &picture)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user: " + err.Error()})
			return
		}
	}

	jwtToken, err := utils.GenerateJWT(user.ID, user.Email, jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": jwtToken})
}

// Logout of the user session
func Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}

// GetProfile returns the profile of the user
func GetProfile(c *gin.Context) {
	userToken := c.GetHeader("Authorization")
	claims, err := utils.ValidateJWT(userToken, jwtSecret)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	user, err := userService.GetUserByID(claims.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching user data: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}
