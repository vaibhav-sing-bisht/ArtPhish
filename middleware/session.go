package middleware

import (
	"encoding/gob"
	"net/http"

	"github.com/gophish/gophish/models"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

// init registers the necessary models to be saved in the session later
func init() {
	gob.Register(&models.User{})
	gob.Register(&models.Flash{})
	Store.Options.HttpOnly = true
	//Secure cookies if request is over HTTPS
	Store.Options.Secure = true                    // Ensure cookies are marked as Secure
	Store.Options.SameSite = http.SameSiteNoneMode // Explicitly set SameSite=None
	// This sets the maxAge to 5 days for all cookies
	// This sets the maxAge to 5 days for all cookies
	Store.MaxAge(86400 * 5)

}

// Store contains the session information for the request
var Store = sessions.NewCookieStore(
	[]byte(securecookie.GenerateRandomKey(64)), //Signing key
	[]byte(securecookie.GenerateRandomKey(32)))
