package gonaturalist

import (
	"crypto/tls"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"net/url"
	"strings"
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

