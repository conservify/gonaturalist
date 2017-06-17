package gonaturalist

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	http          *http.Client
	autoRetry     bool
	retryDuration time.Duration
}

func isFailure(code int, validCodes []int) bool {
	for _, item := range validCodes {
		if item == code {
			return false
		}
	}
	return true
}

func shouldRetry(status int) bool {
	return status == http.StatusAccepted || status == http.StatusTooManyRequests
}

func (c *Client) execute(req *http.Request, result interface{}, needsStatus ...int) error {
	for {
		resp, err := c.http.Do(req)
		if err != nil {
			return err
		}

		defer resp.Body.Close()

		if c.autoRetry && shouldRetry(resp.StatusCode) {
			time.Sleep(c.retryDuration)
			continue
		} else if resp.StatusCode != http.StatusOK && isFailure(resp.StatusCode, needsStatus) {
			errorMessage := c.decodeError(resp)
			return errorMessage
		}

		if result != nil {
			if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
				return err
			}
		}

		break
	}

	return nil
}

type pageHeaders struct {
	totalEntries int
	perPage      int
	page         int
}

func (c *Client) get(url string, result interface{}) (paging *pageHeaders, err error) {
	for {
		resp, err := c.http.Get(url)
		if err != nil {
			return nil, err
		}

		defer resp.Body.Close()

		total := resp.Header["X-Total-Entries"]
		if len(total) > 0 {
			perPage := resp.Header["X-Per-Page"][0]
			page := resp.Header["X-Page"][0]
			t, _ := strconv.Atoi(string(total[0]))
			p, _ := strconv.Atoi(string(page[0]))
			pp, _ := strconv.Atoi(string(perPage[0]))
			paging = &pageHeaders{
				totalEntries: t,
				page:         p,
				perPage:      pp,
			}
		}
		if resp.StatusCode == rateLimitExceededStatusCode && c.autoRetry {
			time.Sleep(c.retryDuration)
			continue
		} else if resp.StatusCode != http.StatusOK {
			errorMessage := c.decodeError(resp)
			return nil, errorMessage
		}

		err = json.NewDecoder(resp.Body).Decode(result)
		if err != nil {
			return nil, err
		}

		break
	}

	return paging, nil
}

func (c *Client) decodeError(resp *http.Response) error {
	return nil
}

const (
	rateLimitExceededStatusCode = 429
)

type Authenticator struct {
	config  *oauth2.Config
	context context.Context
}

func NewAuthenticator(clientId string, clientSecret string, redirectUrl string) Authenticator {
	endpoint := oauth2.Endpoint{
		AuthURL:  "https://www.inaturalist.org/oauth/authorize",
		TokenURL: "https://www.inaturalist.org/oauth/token",
	}

	cfg := &oauth2.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		RedirectURL:  redirectUrl,
		Scopes:       []string{},
		Endpoint:     endpoint,
	}

	tr := &http.Transport{
		TLSNextProto: map[string]func(authority string, c *tls.Conn) http.RoundTripper{},
	}
	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, &http.Client{Transport: tr})
	return Authenticator{
		config:  cfg,
		context: ctx,
	}
}

func (a Authenticator) AuthUrl() string {
	authUrl, err := url.Parse(a.config.Endpoint.AuthURL)
	if err != nil {
		log.Fatal("Parse: ", err)
	}
	parameters := url.Values{}
	parameters.Add("client_id", a.config.ClientID)
	parameters.Add("scope", strings.Join(a.config.Scopes, " "))
	parameters.Add("redirect_uri", a.config.RedirectURL)
	parameters.Add("response_type", "code")
	authUrl.RawQuery = parameters.Encode()
	return authUrl.String()
}

func (a Authenticator) Exchange(code string) (*oauth2.Token, error) {
	return a.config.Exchange(a.context, code)
}

func (a *Authenticator) NewClient(token *oauth2.Token) Client {
	client := a.config.Client(a.context, token)
	return Client{
		http: client,
	}
}

type PrivateUser struct {
	Name      string
	Email     string
	Id        int64
	Login     string
	Uri       string
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	ObservationsCount int32 `json:"observations_count"`

	LifeListId        int32  `json:"life_list_id"`
	LifeListTaxaCount int32  `json:"life_list_taxa_count"`
	TimeZone          string `json:"time_zone"`

	IconUrl         string `json:"icon_url"`
	IconContentType string `json:"icon_content_type"`
	IconFileName    string `json:"icon_file_name"`
	IconFileSize    int32  `json:"icon_file_size"`
}

func (c *Client) CurrentUser() (*PrivateUser, error) {
	var result PrivateUser

	_, err := c.get("https://www.inaturalist.org/users/edit.json", &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

type SimpleObservation struct {
	UserLogin        string    `json:"user_login"`
	PlaceGuess       string    `json:"place_guess"`
	SpeciesGuess     string    `json:"species_guess"`
	Latitude         string    `json:"latitude"`
	Longitude        string    `json:"longitude"`
	CreatedAt        time.Time `json:"created_at"`
	ObservedOnString string    `json:"observed_on_string"`
	UpdatedAt        time.Time `json:"updated_at"`
	TaxonId          int32     `json:"taxon_id"`
	Id               int64     `json:"id"`
	UserId           int64     `json:"user_id"`
	TimeZone         string    `json:"time_zone"`
}

type ObservationsPage struct {
	paging       *pageHeaders
	Observations []SimpleObservation
}

type ProjectObservation struct {
}

type SimpleUser struct {
	Name  string `json:"name"`
	Id    int64  `json:"id"`
	Login string `json:"login"`
}

type Comment struct {
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Id        int64     `json:"id"`
	ParentId  int64     `json:"parent_id"`
	UserId    int64     `json:"user_id"`
	User      SimpleUser
}

type SimplePhoto struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Id        int64     `json:"id"`
	LargeUrl  string    `json:"large_url"`
	MediumUrl string    `json:"medium_url"`
	SmallUrl  string    `json:"small_url"`
	SquareUrl string    `json:"square_url"`
}

type ObservationPhoto struct {
	Id            int64     `json:"id"`
	PhotoId       int64     `json:"photo_id"`
	ObservationId int64     `json:"observation_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Photo         SimplePhoto
}

type FullObservation struct {
	Id               int64                `json:"id"`
	CreatedAt        time.Time            `json:"created_at"`
	UpdatedAt        time.Time            `json:"updated_at"`
	Latitude         string               `json:"latitude"`
	Longitude        string               `json:"longitude"`
	ObservedOnString string               `json:"observed_on_string"`
	Photos           []ObservationPhoto   `json:"observation_photos"`
	Comments         []Comment            `json:"comments"`
	Projects         []ProjectObservation `json:"project_observations"`
}

type GetProjectsOpt struct {
	Page *int
}

type GetObservationsOpt struct {
	Page *int
}

type SimpleProject struct {
	Id               int64                `json:"id"`
	CreatedAt        time.Time            `json:"created_at"`
	UpdatedAt        time.Time            `json:"updated_at"`
	Terms string `json:"terms"`
	Description string `json:"description"`
	Title string `json:"title"`
}

type ProjectsPage struct {
	paging       *pageHeaders
	Projects []SimpleProject
}

func (c *Client) GetProjects(opt *GetProjectsOpt) (*ProjectsPage, error) {
	var result []SimpleProject

	u := "https://www.inaturalist.org/projects.json"
	if opt != nil {
		v := url.Values{}
		if opt.Page != nil {
			v.Set("page", strconv.Itoa(*opt.Page))
		}
		if params := v.Encode(); params != "" {
			u += "?" + params
		}
	}
	p, err := c.get(u, &result)
	if err != nil {
		return nil, err
	}

	return &ProjectsPage{
		Projects: result,
		paging:       p,
	}, nil
}

func (c *Client) GetObservations(opt *GetObservationsOpt) (*ObservationsPage, error) {
	var result []SimpleObservation

	u := "https://www.inaturalist.org/observations.json"
	if opt != nil {
		v := url.Values{}
		if opt.Page != nil {
			v.Set("page", strconv.Itoa(*opt.Page))
		}
		if params := v.Encode(); params != "" {
			u += "?" + params
		}
	}
	p, err := c.get(u, &result)
	if err != nil {
		return nil, err
	}

	return &ObservationsPage{
		Observations: result,
		paging:       p,
	}, nil
}

func (c *Client) GetObservation(id int64) (*FullObservation, error) {
	var result FullObservation

	u := fmt.Sprintf("https://www.inaturalist.org/observations/%d.json", id)
	_, err := c.get(u, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
