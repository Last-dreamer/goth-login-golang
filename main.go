package main

import (
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)
 
var store = sessions.NewCookieStore([]byte(strconv.Itoa(8888)))

func main() {
	// key := RandomInteger(7)

	store.MaxAge(86400 * 30) // 30 days
	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = false

	gothic.Store = store

	os.Setenv("GOOGLE_KEY", "")
	os.Setenv("GOOGLE_SECRET", "")

	goth.UseProviders(
		google.New(os.Getenv("GOOGLE_KEY"), os.Getenv("GOOGLE_SECRET"), "http://localhost:8080/auth/google/callback"),
	)

	// /Request details: redirect_uri=http://localhost:8080/auth/google/callback

	r := gin.Default()

	r.GET("/auth/google", Auth)

	r.GET("/auth/google/callback", AuthCallback)
	r.GET("/getUserDetail", GetUser)

	r.Run(":8080")
}

func RandomInteger(n int) int {
	var sb strings.Builder

	ran := rand.Int63n(999999999999999999)

	strng := strconv.FormatInt(ran, 10)
	k := len(strng)

	for i := 0; i < n; i++ {
		c := strng[rand.Intn(k)]
		sb.WriteByte(c)
	}
	finalInt, _ := strconv.Atoi(sb.String())
	return finalInt
}

func Auth(cntx *gin.Context) {
	gothic.BeginAuthHandler(cntx.Writer, cntx.Request)
}

func GetUser(c *gin.Context) {

	session, _ := store.Get(c.Request, "user_session") // Ensure you have initialized 'store'
	userValue, ok := session.Values["user"]
	if !ok || userValue == nil {
		// Handle the case where there is no user in the session
		// For example, redirect to the login page or return an error response
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not logged in"})
		return
	}
	user, ok := userValue.(goth.User)
	if !ok {
		// Handle the case where the session value is not of type goth.User
		// This should not happen if your session management is correct
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	// Now you can safely use the 'user' object
	// For example, return the user's information as a JSON response
	c.JSON(http.StatusOK, gin.H{"user": user})
}

// testing for flutter
func AuthCallback(cntx *gin.Context) {

	// Pull the provider from the path parameters
	query := cntx.Request.URL.Query()
	query.Add("provider", "google")
	cntx.Request.URL.RawQuery = query.Encode()

	gothic.SetState(cntx.Request)

	// Call CompleteUserAuth to have it process the response
	gothUser, err := gothic.CompleteUserAuth(cntx.Writer, cntx.Request)
	if err != nil {
		log.Println("error while getting goth User", err)
		return
	}

	log.Println("testing gothUser ", gothUser)
	session, _ := store.Get(cntx.Request, "user_session")
	session.Values["user"] = gothUser
	session.Save(cntx.Request, cntx.Writer)

	cntx.JSON(200, gothUser)
	// cntx.Redirect(http.StatusPermanentRedirect, "http://20.203.31.58/dashboard/default")

}
