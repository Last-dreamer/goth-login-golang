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

func main() {

	key := RandomInteger(7)
	store := sessions.NewCookieStore([]byte(strconv.Itoa(key)))
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

func AuthCallback(cntx *gin.Context) {

	// Pull the provider from the path parameters
	query := cntx.Request.URL.Query()
	query.Add("provider", cntx.Param("provider"))
	cntx.Request.URL.RawQuery = query.Encode()

	gothic.SetState(cntx.Request)

	// Call CompleteUserAuth to have it process the response
	gothUser, err := gothic.CompleteUserAuth(cntx.Writer, cntx.Request)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("testing goth User ", gothUser.ExpiresAt)
	// Save the needed information from the response
	// cntx.HTML(http.StatusOK, "success", nil)
	cntx.Redirect(http.StatusPermanentRedirect, "http://20.203.31.58/dashboard/default")
}
