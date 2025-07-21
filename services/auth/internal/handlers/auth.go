package handlers

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/resume-optimizer/auth-service/internal/models"
	"github.com/resume-optimizer/auth-service/internal/services"
	"github.com/resume-optimizer/auth-service/internal/utils"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	userService *services.UserService
	jwtSecret   = os.Getenv("JWT_SECRET")
	oauth2Config = &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Endpoint:     google.Endpoint,
		RedirectURL:  os.Getenv("REDIRECT_URL"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
	}
)

// InitHandlers initializes the handlers after database is ready
func InitHandlers() {
	log.Println("Initializing handlers with database connection")
	userService = services.NewUserService()
	if userService == nil {
		log.Fatal("Failed to create user service")
	}
	log.Println("Handlers initialized successfully")
}

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

	// Redirect to frontend with token
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:3000"
	}
	
	redirectURL := frontendURL + "/auth/callback?token=" + jwtToken
	c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

// Logout of the user session
func Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}

// GetProfile returns the profile of the user
func GetProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	user, err := userService.GetUserByID(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching user data: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// GetUserAPIKeys returns user's API keys
func GetUserAPIKeys(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	apiKeys, err := userService.GetUserAPIKeys(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching API keys: " + err.Error()})
		return
	}

	// Return masked keys for security
	var responseKeys []gin.H
	for _, key := range apiKeys {
		responseKeys = append(responseKeys, gin.H{
			"id":         key.ID,
			"provider":   key.Provider,
			"masked_key": maskAPIKey(key.EncryptedKey),
			"created_at": key.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{"api_keys": responseKeys})
}

// CreateUserAPIKey creates a new API key for the user
func CreateUserAPIKey(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req struct {
		Provider string `json:"provider" binding:"required"`
		APIKey   string `json:"api_key" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// Encrypt the API key
	encryptedKey, err := encryptAPIKey(req.APIKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt API key"})
		return
	}

	userIDStr := userID.(string)
	log.Printf("Creating API key for user: %s, provider: %s", userIDStr, req.Provider)
	
	apiKey := &models.UserAPIKey{
		ID:           uuid.New().String(),
		UserID:       userIDStr,
		Provider:     req.Provider,
		EncryptedKey: encryptedKey,
	}

	log.Printf("API key struct: ID=%s, UserID=%s, Provider=%s", apiKey.ID, apiKey.UserID, apiKey.Provider)

	if err := userService.CreateUserAPIKey(apiKey); err != nil {
		log.Printf("Failed to create API key: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save API key: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":         apiKey.ID,
		"provider":   apiKey.Provider,
		"masked_key": maskAPIKey(encryptedKey),
		"created_at": apiKey.CreatedAt,
	})
}

// DeleteUserAPIKey deletes a user's API key
func DeleteUserAPIKey(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	keyID := c.Param("id")
	if err := userService.DeleteUserAPIKey(userID.(string), keyID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete API key: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "API key deleted successfully"})
}

// Helper functions for encryption/decryption
func encryptAPIKey(apiKey string) (string, error) {
	key := []byte(os.Getenv("ENCRYPTION_KEY"))
	if len(key) != 32 {
		// Use a default key if not set (not recommended for production)
		key = []byte("your-32-byte-encryption-key-here!!")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	plaintext := []byte(apiKey)
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return hex.EncodeToString(ciphertext), nil
}

func decryptAPIKey(encryptedKey string) (string, error) {
	key := []byte(os.Getenv("ENCRYPTION_KEY"))
	if len(key) != 32 {
		key = []byte("your-32-byte-encryption-key-here!!")
	}

	ciphertext, err := hex.DecodeString(encryptedKey)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		return "", err
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext), nil
}

func maskAPIKey(encryptedKey string) string {
	if len(encryptedKey) < 8 {
		return "****"
	}
	return encryptedKey[:4] + "****" + encryptedKey[len(encryptedKey)-4:]
}
