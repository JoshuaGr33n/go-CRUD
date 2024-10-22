
package main
import (
	"github.com/gin-gonic/gin"
	"unicode" 
	"net/http"
	"github.com/gin-contrib/sessions"
	"crypto/rand"
    "encoding/base64"
)

// var store = cookie.NewStore([]byte("secret"))
// router.Use(sessions.Sessions("mysession", store))

func containsUpperLowerDigitSpecial(s string) bool {
	var (
		hasUpper, hasLower, hasDigit, hasSpecial bool
	)
	for _, char := range s {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	return hasUpper && hasLower && hasDigit && hasSpecial
}


// GetUserFromContext retrieves the current logged-in user from the Gin context
func GetUserFromContext(c *gin.Context) *User {
	user, exists := c.Get("user")
	if !exists {
		return nil
	}
	return user.(*User)
}

// Middleware to populate the current user in the Gin context
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve user from session or JWT token and set it in the Gin context
		user := getUserFromSession(c) // Logic to retrieve user from session or JWT token
		c.Set("user", user)
		c.Next()
	}
}

// Handler to get the current user
func CurrentUserHandler(c *gin.Context) {
	user := GetUserFromContext(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user})
}

func getUserFromSession(c *gin.Context) *User {
	// Retrieve the session
	session := sessions.Default(c)

	// Retrieve the user from the session
	user := session.Get("user")
	if user != nil {
		// If user is found in the session, assert the type to *User and return
		if u, ok := user.(*User); ok {
			return u
		}
	}

	// Return nil if user is not found or the type assertion fails
	return nil
}


func generateSecretKey() (string, error) {
    // Generate a random byte slice with sufficient length
    key := make([]byte, 32)
    if _, err := rand.Read(key); err != nil {
        return "", err
    }

    // Encode the random bytes to base64 to get a string representation
    encodedKey := base64.URLEncoding.EncodeToString(key)

    return encodedKey, nil
}

func getUser(c *gin.Context) User {
	session := sessions.Default(c)
	user := session.Get("user").(User)
	return user
}


func AuthRequired() gin.HandlerFunc {
    return func(c *gin.Context) {
        session := sessions.Default(c)
        user := session.Get("user")
        if user == nil {
            // Redirect to login page if user is not authenticated
            c.Redirect(http.StatusFound, "/login")
            c.Abort()
            return
        }
        c.Next()
    }
}

