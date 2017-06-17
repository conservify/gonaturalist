package gonaturalist

import (
	"time"
	"net/url"
	"strconv"
)

type GetProjectsOpt struct {
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
