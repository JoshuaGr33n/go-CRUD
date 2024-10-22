package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
)

func main() {
	db = getMySQLDB()
	migrateDB()

	router := gin.Default()

	key, err := generateSecretKey()
	if err != nil {
        fmt.Println("Error:", err)
        return
    }
	store := cookie.NewStore([]byte(key))
	router.Use(sessions.Sessions("session-name", store))


	router.Static("/assets", "./assets")
	router.LoadHTMLGlob("templates/*.html")

	router.GET("/", homePage)
	router.POST("/", homePage)
	router.GET("/records/:id", fetchRecord)
	router.POST("/records/update", updateRecord)
	router.DELETE("/records/:id", deleteRecord)
	router.GET("/register", register)
	router.POST("/register-process", registerHandler)
	router.GET("/login", login)
	router.POST("/login-process", loginHandler)
	router.GET("/dashboard", AuthRequired(), dashboard)
	router.GET("/logout", logout)
	router.Run(":8999")
}
