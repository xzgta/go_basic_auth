package main

import (
	"fmt"
	"log"
	"net/http"
	"pratice/config"
	"pratice/structs"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type Credentials struct {
	Username string `form:"username" json:"username"`
	Password string `form:"password" json:"password"`
}

func main() {

	// inDB := &controllers.InDB{DB: db}

	router := gin.Default()
	router.POST("/login", loginHandler)
	router.GET("/home", auth, home)
	router.Run(":3000")
}
func home(c *gin.Context) {
	log.Println("sukses")
}
func loginHandler(c *gin.Context) {
	db := config.Dbconn()
	defer db.Close()

	// r.ParseForm()
	username := c.PostForm("username")
	password := c.PostForm("password")

	u := structs.User{}
	result := db.QueryRow("SELECT username, password FROM tb_user WHERE username= ?;", username).Scan(&u.Username, &u.Password)
	log.Println(result)
	log.Println(u.Username)
	log.Println(u.Password)
	hash := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if username != "" && password != "" {
		if u.Username == username && hash == nil {
			sign := jwt.New(jwt.GetSigningMethod("HS256"))
			token, err := sign.SignedString([]byte("secret"))
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": err.Error(),
				})
				c.Abort()
			}
			c.JSON(http.StatusOK, gin.H{
				"token": token,
			})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  http.StatusUnauthorized,
				"message": "wrong username or password",
			})
		}
	}

}

func auth(c *gin.Context) {
	tokenString := c.Request.Header.Get("Authorization")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod("HS256") != token.Method {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte("secret"), nil
	})

	// if token.Valid && err == nil {
	if token != nil && err == nil {
		fmt.Println("token verified")
	} else {
		result := gin.H{
			"message": "not authorized",
			"error":   err.Error(),
		}
		c.JSON(http.StatusUnauthorized, result)
		c.Abort()
	}
}
