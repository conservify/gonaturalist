package main

import (
	"fmt"
	"github.com/Conservify/gonaturalist"
	"golang.org/x/oauth2"
	"log"
	"net/http"
)

var authenticator = gonaturalist.NewAuthenticator(applicationId, secret, redirectUrl)

func handleNaturalistLogin(w http.ResponseWriter, r *http.Request) {
	url := authenticator.AuthUrl()
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func completeAuth(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")

	token, err := authenticator.Exchange(code)

	_ = err
	fmt.Println(token.AccessToken)
}

func main() {
	if accessToken == "" {
		http.HandleFunc("/login", handleNaturalistLogin)
		http.HandleFunc("/callback", completeAuth)
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			log.Println("Ignoring:", r.URL.String())
		})

		http.ListenAndServe(":8000", nil)
	}

	var oauthToken oauth2.Token
	oauthToken.AccessToken = accessToken
	c := authenticator.NewClient(&oauthToken)

	log.Printf("GetCurrentUser:")

	user, err := c.GetCurrentUser()
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	fmt.Printf("%v\n\n", user)

	log.Printf("GetObservations:")

	observations, err := c.GetObservations(nil)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	fmt.Printf("%v\n\n", observations)

	log.Printf("GetObservation:")

	o, err := c.GetObservation(observations.Observations[0].Id)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	fmt.Printf("%v\n\n", o)

	log.Printf("GetProjects:")

	projects, err := c.GetProjects(nil)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	fmt.Printf("%v\n\n", projects)

	log.Printf("GetProject:")

	project, err := c.GetProject("the-sonoran-desert")
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	fmt.Printf("%v\n\n", project)
}
