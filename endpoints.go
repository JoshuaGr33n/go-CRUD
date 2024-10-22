package main

import (
	// "fmt"
	"github.com/gin-gonic/gin"
	// "log"/
	"encoding/gob"
	"github.com/gin-contrib/sessions"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
	"regexp"
)

type Crud struct {
	gorm.Model
	// ID       uint   `gorm:"primaryKey"`  (to skip timestamp columns)
	Name    string
	Phone   string
	Email   string
	Message string
}

func init() {
	// Register User struct with gob
	gob.Register(User{})
}

// var users = map[string]User{}

func homePage(c *gin.Context) {
	if c.Request.Method == http.MethodPost {

		var message string

		text := Crud{
			Name:    c.PostForm("name"),
			Phone:   c.PostForm("phone"),
			Email:   c.PostForm("email"),
			Message: c.PostForm("message"),
		}

		switch c.PostForm("submit") {
		case "Insert":
			message = "Insert"
			err := insertRecord(text)
			if err != nil {
				message = "Insert failed: " + err.Error()
			} else {
				message = "Insert successful"
			}
		case "Read":
			message = "Read"
		case "Update":
			message = "Update"
		case "Delete":
			message = "Delete"
		}
		c.HTML(http.StatusOK, "index.html", gin.H{
			"Success": true,
			"Message": message,
		})
		return // Return from the function after sending the response
	}
	cruds := selectRecords()

	// Render the template with data
	c.HTML(http.StatusOK, "index.html", gin.H{
		"Cruds": cruds,
	})
}

func insertRecord(info Crud) error {
	// _, err := db.Exec("INSERT INTO crud (name, phone, email, message) VALUES (?, ?, ?, ?)", info.Name, info.Phone, info.Course, info.Message)
	// return err
	result := db.Create(&info)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func selectRecords() []Crud {
	var cruds []Crud
	db.Find(&cruds)
	return cruds
}

func fetchRecord(c *gin.Context) {
	// Extract record ID from URL parameter
	id := c.Param("id")

	// Query the database for the record with the specified ID
	var record Crud
	result := db.First(&record, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch record"})
		return
	}

	// Return the record as JSON response
	c.JSON(http.StatusOK, record)
}

func updateRecord(c *gin.Context) {
	// Get the record ID from the URL parameter

	// Parse the request form data to get updated record information
	var updatedRecord Crud
	if err := c.ShouldBind(&updatedRecord); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id := c.PostForm("id")

	// Find the record in the database by ID
	var existingRecord Crud
	if err := db.First(&existingRecord, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}

	// Update the record with the new information
	// existingRecord.Name = updatedRecord.Name
	// existingRecord.Phone = updatedRecord.Phone
	// existingRecord.Email = updatedRecord.Email
	// existingRecord.Message = updatedRecord.Message
	existingRecord.Name = c.PostForm("name")
	existingRecord.Phone = c.PostForm("phone")
	existingRecord.Email = c.PostForm("email")
	existingRecord.Message = c.PostForm("message")

	// Save the updated record back to the database
	if err := db.Save(&existingRecord).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update record"})
		return
	}

	// Return a success response
	c.JSON(http.StatusOK, gin.H{"message": "Record updated successfully"})
}

func deleteRecord(c *gin.Context) {
	id := c.Param("id")
	var record Crud
	result := db.Delete(&record, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	// Return a success response
	c.JSON(http.StatusOK, gin.H{"message": "Record deleted successfully"})
}

func register(c *gin.Context) {
	// Render the template with data
	c.HTML(http.StatusOK, "register.html", gin.H{})
}

func registerHandler(c *gin.Context) {
	// var user User

	user := User{
		Username: c.PostForm("username"),
		Email:    c.PostForm("email"),
		Password: c.PostForm("password"),
	}

	// Bind request data to User struct
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate username
	if len(user.Username) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username must be at least 2 characters"})
		return
	}

	// Validate email format
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(user.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email address"})
		return
	}

	// Validate password complexity
	if len(user.Password) < 8 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 8 characters"})
		return
	}
	if !containsUpperLowerDigitSpecial(user.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must contain at least one uppercase letter, one lowercase letter, one digit, and one special character"})
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user.Password = string(hashedPassword)

	// Create the user record
	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

func login(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

func loginHandler(c *gin.Context) {

	// var form User
	form := User{
		Username: c.PostForm("username"),
		Password: c.PostForm("password"),
	}

	if err := c.ShouldBind(&form); err != nil {
		// c.HTML(http.StatusBadRequest, "login.html", gin.H{"error": err.Error()})
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user User
	if err := db.Where("username = ?", form.Username).First(&user).Error; err != nil {
		// c.HTML(http.StatusUnauthorized, "login.html", gin.H{"error": "Invalid username or password"})
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username"})
		return
	}

	// Compare the hashed password from the database with the provided password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	// Authentication successful
	// c.HTML(http.StatusOK, "login.html", gin.H{"message": fmt.Sprintf("Welcome, %s!", form.Username)})
	// c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Welcome, %s!", form.Username)})

	//Store the user in the session
	session := sessions.Default(c)
	session.Set("user", user)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}
	// session := sessions.Default(c)
	// session.Set("user", user)
	// fmt.Println("Session before save:", session)
	// if err := session.Save(); err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session: " + err.Error()})
	// 	return
	// }

	// Redirect to the dashboard
	c.Redirect(http.StatusFound, "/dashboard")
}

func dashboard(c *gin.Context) {
	user := getUser(c)
	c.HTML(http.StatusOK, "dashboard.html", gin.H{"user": user})
}

func logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear session"})
		return
	}
	c.Redirect(http.StatusFound, "/login") // Redirect to the login page after logout
}
