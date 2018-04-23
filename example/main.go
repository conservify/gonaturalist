package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Conservify/gonaturalist"
)

var authenticator = gonaturalist.NewAuthenticatorAtCustomRoot(applicationId, secret, redirectUrl, "http://127.0.0.1:3000")

func handleNaturalistLogin(w http.ResponseWriter, r *http.Request) {
	url := authenticator.AuthUrl()

	log.Printf("Redirecting: %s", url)

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func completeAuth(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")

	token, err := authenticator.Exchange(code)
	if err != nil {
		log.Printf("Error: %v", err)
	}

	log.Printf("AccessToken: %s", token.AccessToken)
}

func main() {
	if accessToken == "" {
		log.Printf("No access token, staring web server.")

		http.HandleFunc("/login", handleNaturalistLogin)
		http.HandleFunc("/callback", completeAuth)
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			log.Println("Ignoring:", r.URL.String())
		})

		log.Printf("Open http://127.0.0.1:8000/login")

		http.ListenAndServe(":8000", nil)
	}

	c := authenticator.NewClientWithAccessToken(accessToken, &gonaturalist.NoopCallbacks{})

	log.Printf("GetCurrentUser:")
	user, err := c.GetCurrentUser()
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	fmt.Printf("%v\n\n", user)

	log.Printf("GetObservationsByUsername:")
	myObservations, err := c.GetObservationsByUsername(user.Login)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	hasObservations := len(myObservations.Observations) > 0
	for _, o := range myObservations.Observations {
		if true {
			updateObservation := gonaturalist.UpdateObservationOpt{
				Id:          o.Id,
				Description: "Updated!",
			}
			err := c.UpdateObservation(&updateObservation)
			if err != nil {
				log.Fatalf("Error: %v", err)
			}
		}
		fmt.Printf("%v\n", o)
	}
	fmt.Printf("\n")

	if !hasObservations {
		addObservation := gonaturalist.AddObservationOpt{
			SpeciesGuess:       "Duck",
			ObservedOnString:   time.Now(),
			Description:        "Look what I found!",
			Latitude:           41.27872259999999,
			Longitude:          -72.5276073,
			PositionalAccuracy: 1,
		}
		err := c.AddObservation(&addObservation)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
		fmt.Printf("\n")
	}

	if false {
		log.Printf("GetProject:")
		project, err := c.GetProject("the-sonoran-desert")
		if err != nil {
			fmt.Printf("%v\n", err)
		}
		fmt.Printf("%v\n\n", project)
		fmt.Printf("\n")
	}

	if false {
		log.Printf("GetPlaces:")
		places, err := c.GetPlaces(nil)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
		for _, place := range places.Places {
			fmt.Printf("%v\n", place)
		}
		fmt.Printf("\n")
	}

	if true {
		lon := -118.25
		lat := 34.05
		log.Printf("GetPlaces(%v, %v):", lon, lat)
		places, err := c.GetPlaces(&gonaturalist.GetPlacesOpt{Longitude: &lon, Latitude: &lat})
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
		for _, place := range places.Places {
			fmt.Printf("%v\n", place)
		}
		fmt.Printf("\n")
	}

	if false {
		log.Printf("GetObservations:")
		observations, err := c.GetObservations(nil)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
		for _, observation := range observations.Observations {
			fmt.Printf("%v\n", observation)
		}
		fmt.Printf("\n")
	}

	if false {
		log.Printf("GetObservation(%d):", 100)
		o, err := c.GetObservation(100)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
		fmt.Printf("%v\n", o)
		fmt.Printf("\n")

		comms, err := c.GetObservationComments(100)
		for _, c := range comms {
			fmt.Printf("%v\n", c)
		}
		fmt.Printf("\n")

		if len(comms) == 0 {
			addComment := gonaturalist.AddCommentOpt{
				ParentType: gonaturalist.Observation,
				ParentId:   100,
				Body:       "Hello, world!",
			}
			err = c.AddComment(&addComment)
			if err != nil {
				log.Fatalf("Error: %v", err)
			}
		} else {
			err = c.UpdateCommentBody(comms[0].Id, "Goodbye!")
			if err != nil {
				log.Fatalf("Error: %v", err)
			}
		}
	}

	if false {
		on, err := time.Parse("2006-01-02", "2011-11-15")
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
		log.Printf("GetObservations(%v):", on)
		observations, err := c.GetObservations(&gonaturalist.GetObservationsOpt{On: &on})
		if err != nil {
			fmt.Printf("%v\n", err)
		}
		for _, observation := range observations.Observations {
			fmt.Printf("%v\n", observation)
		}
		fmt.Printf("\n")
	}
}
