package gonaturalist

import (
	"net/url"
	"strconv"
	"time"
)

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
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
	ObservationId           int64     `json:"observation_id"`
	Id                      int64     `json:"id"`
	TrackingCode            string    `json:"tracking_code"`
	CuratorIdentificationId int64     `json:"curator_identification_id"`
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

type GetObservationsOpt struct {
	Page *int
}

func (c *Client) GetObservations(opt *GetObservationsOpt) (*ObservationsPage, error) {
	var result []SimpleObservation

	u := c.buildUrl("/observations.json")
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

type AddObservationOpt struct {
}

func (c *Client) AddObservation(opt *AddObservationOpt) error {
	return nil
}

func (c *Client) GetObservation(id int64) (*FullObservation, error) {
	var result FullObservation

	u := c.buildUrl("/observations/%d.json", id)
	_, err := c.get(u, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

type UpdateObservationOpt struct {
}

func (c *Client) UpdateObservation(opt *UpdateObservationOpt) error {
	return nil
}

type DeleteObservationOpt struct {
}

func (c *Client) DeleteObservation(opt *DeleteObservationOpt) error {
	return nil
}

func (c *Client) GetObservationsByUsername(username string) (*ObservationsPage, error) {
	var result []SimpleObservation

	u := c.buildUrl("/observations/%s.json", username)
	p, err := c.get(u, &result)
	if err != nil {
		return nil, err
	}

	return &ObservationsPage{
		Observations: result,
		paging:       p,
	}, nil
}
