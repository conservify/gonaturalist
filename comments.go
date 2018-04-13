package gonaturalist

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type CommentParentType string

const (
	AssessmentSection CommentParentType = "AssessmentSection"
	ListedTaxon       CommentParentType = "ListedTaxon"
	Observation       CommentParentType = "Observation"
	ObservationField  CommentParentType = "ObservationField"
	Post              CommentParentType = "Post"
	TaxonChange       CommentParentType = "TaxonChange"
)

type AddCommentOpt struct {
	ParentType CommentParentType `json:"parent_type"`
	ParentId   int64             `json:"parent_id"`
	Body       string            `json:"body"`
}

func (c *Client) AddComment(opt *AddCommentOpt) error {
	u := c.buildUrl("/comments.json")

	bodyJson, err := json.Marshal(opt)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", u, bytes.NewReader(bodyJson))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	var p interface{}
	err = c.execute(req, &p, http.StatusCreated)
	if err != nil {
		return err
	}

	return nil
}

type UpdateCommentOpt struct {
	ParentType CommentParentType `json:"parent_type"`
	ParentId   int64             `json:"parent_id"`
	Body       string            `json:"body"`
}

func (c *Client) UpdateComment(opt *UpdateCommentOpt) error {
	return nil
}

type DeleteCommentOpt struct {
}

func (c *Client) DeleteComment(opt *DeleteCommentOpt) error {
	return nil
}
