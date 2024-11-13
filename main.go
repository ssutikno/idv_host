package main

import (
	"encoding/json"
	"idv_host/handlers"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// User represents the user information
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// make default user file
// the default user file is used for the user authentication

var users []User
var jwtSecret = []byte("strong-secret-password") // Replace with a strong secret key

func init() {
	var data []byte
	var err error
	data, err = ioutil.ReadFile("users.json")
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(data, &users)
	if err != nil {
		log.Fatal(err)
	}

}

// generate a JWT token
// the token is used for the user authentication
// the token is generated based on the user information
// the token is signed with a secret key
func generateToken(user User) (string, error) {
	// create a new token
	token := jwt.New(jwt.SigningMethodHS256)
	// set the claims
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = user.Username
	claims["password"] = user.Password
	// set the expiration time
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	// sign the token with the secret key
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// authMiddleware is a middleware that checks if the user is authenticated
// If the user is authenticated, it will call the next handler
// If the user is not authenticated, it will return 401
// if there is token on header, then check token validity
// if the token is valid, then call the next handler
// if the token is not valid, then return 401

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// get the token from the Authorization header
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// parse the token split the token
		// remove the "Bearer " prefix from the token string
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// check if the token is valid
		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// call the next handler
		c.Next()
	}
}

// authenicate the user
func Auth(c *gin.Context) {
	var user User
	log.Println("Auth 1", c)

	log.Println("User / password", c.PostForm("username"), " / ", c.PostForm("password"))
	// bind the user information from the form data
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	for _, u := range users {
		if u.Username == user.Username && u.Password == user.Password {
			token, err := generateToken(u)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
				return
			}

			// Save the token to the Authorization header
			c.Header("Authorization", token)
			c.JSON(http.StatusOK, gin.H{"token": token})
			return
		}
	}

	c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
}

// handle the login handler
func LoginHandler(c *gin.Context) {
	log.Println("LoginHandler")

	// Check if the template file exists
	if _, err := os.Stat("templates/login.html"); os.IsNotExist(err) {
		log.Println("Template file not found:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Render the login page
	c.HTML(http.StatusOK, "login.html", nil)
}

// main function
func main() {

	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	// add the handlers for the API, with the authMiddleware. make sure the user is authenticated
	api := router.Group("/api", authMiddleware())
	{

		api.GET("/vms", handlers.ListVMs)
		api.POST("/vms/start", handlers.StartVM)
		api.POST("/vms/reboot", handlers.RebootVM)
		api.POST("/vms/reset", handlers.ResetVM)
		api.POST("/vms/shutdown", handlers.ShutdownVM)
		api.POST("/vms/poweroff", handlers.PowerOffVM)
		api.POST("/vms/create", handlers.CreateVM)

		api.GET("/host/network", handlers.GetNetworkData)

		api.POST("/host/restart", handlers.RestartHost)
		api.POST("/host/reset", handlers.ResetHost)
	}

	// handle login form
	router.GET("/login", LoginHandler)

	// handle authentication
	router.POST("/login", Auth)

	// handle for the homehandler page
	router.GET("/", handlers.HomeHandler)

	// Start the server

	log.Println("Server starting on port 8080...")
	router.Run(":8080")
}
