package auth

import (
	"encoding/json"
	"fmt"
	dataLayer "github.com/zacharyworks/huddle-gateway/data-layer"
	types "github.com/zacharyworks/huddle-shared/data"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"net/http"
)

// Authoriser contains auth stuff
type Authoriser struct {
	config *oauth2.Config
}

// NewAuth creates an oAuth provided the secrets
func NewAuth(clientID string, clientSecret string) *Authoriser {
	conf := &oauth2.Config{
		RedirectURL:  "http://localhost:8080/callback",
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes: []string{"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint: google.Endpoint}
	authoriser := Authoriser{config: conf}

	return &authoriser
}

// HandleLogin directs to oAuth login page
func (a *Authoriser) HandleLogin(w http.ResponseWriter, r *http.Request) {
	// Retrieve session and generate random state
	session := r.FormValue("session")
	state, err := GetRandomString(8)
	if err != nil {
		fmt.Println("Error", err)
		return
	}

	dataLayer.SaveNewSession(session, state)

	// Redirect to be authorised
	url := a.config.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// HandleCallback handles response from oAuth2 api
func (a *Authoriser) HandleCallback(w http.ResponseWriter, r *http.Request) {
	// Attempt to find an existing session for the state
	providedState := r.FormValue("state")
	session, err := dataLayer.RetrieveSessionByState(providedState)
	println(session.State)
	if err != nil {
		fmt.Printf("Error finding session for state: %s\n", err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	token, err := a.config.Exchange(oauth2.NoContext, r.FormValue("code"))
	if err != nil {
		fmt.Printf("Could not get token %s\n", err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// this is our token to access the user info from google, we've already authorised user at this point
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		fmt.Printf("Could not get request %s\n", err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// put response into a structure
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	var response types.Response
	err = json.Unmarshal(content, &response)
	if err != nil {
		fmt.Printf("Error decoding response %s\n", err.Error())
	}

	// does the user exist?
	userOauthID := dataLayer.GetUser(response.ID).OauthID
	if userOauthID == "" {
		postOauthUser(response)
		// now it's been posted, go get it again
		userOauthID = dataLayer.GetUser(response.ID).OauthID
	} else {
		// update the user with any changes
		putOauthUser(response)
	}

	dataLayer.UpdateSession(session, userOauthID)

	// add user ID to session, as it's probably the user
	if err != nil {
		fmt.Printf("could not parse response %s\n", err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	http.Redirect(w, r, "http://localhost:3000", http.StatusPermanentRedirect)
}
