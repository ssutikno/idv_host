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

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// get the token from the Authorization header
		tokenString := c.GetHeader("Authorization")
		// print the token
		// log.Println("Auth Midleware Token:", tokenString)

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
		// log the token
		// log.Println("Auth Middleware - Token Exiting :")

		// Store the intended URL (if needed)
		c.Set("intendedUrl", c.Request.URL.Path)
		c.Request.URL.Path = c.FullPath()
		c.Next()
	}
}

// authenicate the user
func Auth(c *gin.Context) {
	var user User
	// log.Println("Auth 1", c)

	// log.Println("User / password", c.PostForm("username"), " / ", c.PostForm("password"))
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
	// log.Println("LoginHandler")

	// Check if the template file exists
	if _, err := os.Stat("templates/login.html"); os.IsNotExist(err) {
		log.Println("Template file not found:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Render the login page
	c.HTML(http.StatusOK, "login.html", nil)
}

// function for user check exist on the users.json
func userExists(username, password string) bool {
	// return false if the username or password is empty
	if username == "" || password == "" {
		return false
	}

	for _, u := range users {
		if u.Username == username && u.Password == password {
			return true
		}
	}
	return false
}

// function login for api to get token
func LoginApi(c *gin.Context) {
	log.Println("Login API Handler")
	// get the username and password from form data
	username := c.PostForm("username")
	password := c.PostForm("password")
	// log the form data
	log.Println("LoginHandler - username:", username)
	log.Println("LoginHandler - password:", password)

	// check if the user exists

	if userExists(username, password) {
		// return the token
		token, err := generateToken(User{Username: username, Password: password})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"token": token})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

}

// make another authMiddleware function for the homehandler
func homeauthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// get the username and password from form data

		// get form data of the username and password
		username := c.PostForm("username")
		password := c.PostForm("password")
		// log the form data
		// log.Println("HomeAuthMiddleware - username:", username)
		// log.Println("HomeAuthMiddleware - password:", password)

		// check if the username and password are empty
		if username == "" || password == "" {
			c.Redirect(http.StatusFound, "/login")
			// log.Println("LoginHandler - empty username or password")
			c.Abort()
			return
		}

		// check if the username and password are valid, then redirect to the homehandler
		found := false
		for _, u := range users {
			if u.Username == username && u.Password == password {
				found = true
				break
			}
		}

		if found {
			// add authorization token header and Bearer to the response
			token, err := generateToken(User{Username: username, Password: password})
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
				return
			}
			// add Bearer to the token
			c.Request.Header.Set("Authorization", "Bearer "+token)
			// add token to the data
			c.Set("token", token)

			// call the next handler
			c.Next()
		} else {
			c.Redirect(http.StatusFound, "/login")
			log.Println("HomeAuthMiddleware - invalid username or password")
		}

	}
}

// main function
func main() {
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	// add the handlers for the API, with the authMiddleware. make sure the user is authenticated
	api := router.Group("/api/vms", authMiddleware())
	{

		api.GET("/list", handlers.ListVMs)
		api.POST("/start", handlers.StartVM)
		api.POST("/reboot", handlers.RebootVM)
		api.POST("/reset", handlers.ResetVM)
		api.POST("/shutdown", handlers.ShutdownVM)
		api.POST("/poweroff", handlers.PowerOffVM)

	}

	hostapi := router.Group("/api/host", authMiddleware())
	{
		hostapi.GET("/createvm", handlers.CreateVM)
		hostapi.GET("/restart", handlers.RestartHost)
		hostapi.GET("/reset", handlers.ResetHost)
		hostapi.GET("/network", handlers.GetNetworkData)
	}

	// handle login form
	router.GET("/login", LoginHandler)
	router.GET("/createvmform", handlers.CreateVMForm)

	// handle authentication from login form, and redirect to home page
	router.POST("/login", homeauthMiddleware(), handlers.HomeHandler)

	// handle the login api
	router.POST("/loginapi", LoginApi)

	// handle for the homehandler page, if not athenticated, then redirect to login page
	router.GET("/", homeauthMiddleware(), handlers.HomeHandler)

	// Start the server

	// check if the port is already used
	// if the port is already used, then exit the program and print the error

	log.Println("Server starting on port 8080...")

	if err := router.Run(":8080"); err != nil {
		log.Fatal("Server starting error:", err)
	}

}
