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

func main() {
	router := gin.Default()
	router.POST("/login", loginHandler)
	router.POST("/register", registHandler)
	router.GET("/me", auth, me)
	router.Run(":3000")
}
func me(c *gin.Context) {
	db := config.Dbconn()
	defer db.Close()
	token := c.Request.Header.Get("Authorization")
	u := structs.User{}
	db.QueryRow("SELECT username FROM tb_user WHERE token = ?;", token).Scan(&u.Username)
	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"username": &u.Username,
	})
}

func registHandler(c *gin.Context) {
	db := config.Dbconn()
	defer db.Close()
	username := c.PostForm("username")
	password := c.PostForm("password")
	if username != "" && password != "" {
		err := c.Request.ParseForm()
		if err != nil {
			panic(err)
		}

		hash, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
		_, err = db.Exec("INSERT INTO tb_user (username, password) VALUES (?,?)",
			username,
			hash,
		)

		if err != nil {
			log.Print(err)
		}
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "Success Regist :)",
		})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  http.StatusUnauthorized,
			"message": "Failed, Check again :(",
		})
	}

}
func loginHandler(c *gin.Context) {
	db := config.Dbconn()
	defer db.Close()

	username := c.PostForm("username")
	password := c.PostForm("password")

	u := structs.User{}
	db.QueryRow("SELECT username, password FROM tb_user WHERE username= ?;", username).Scan(&u.Username, &u.Password)
	hash := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if username != "" && password != "" {
		if u.Username == username && hash == nil {
			sign := jwt.New(jwt.GetSigningMethod("HS256"))
			claims := sign.Claims.(jwt.MapClaims)
			claims["user"] = u.Username
			token, err := sign.SignedString([]byte("Note2019_"))
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": err.Error(),
				})
				c.Abort()
			}
			db.Exec("UPDATE tb_user SET token =  ? WHERE username = ? ", token, username)
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

		return []byte("Note2019_"), nil
	})
	log.Println(token)
	if token != nil && err == nil {
		fmt.Println("token verified")
	} else {
		result := gin.H{
			"status":  http.StatusOK,
			"message": "not authorized",
			"error":   err.Error(),
		}
		c.JSON(http.StatusUnauthorized, result)
		c.Abort()
	}
}
