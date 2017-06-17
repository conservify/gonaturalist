package gonaturalist

import (
	_ "bytes"
	_ "encoding/json"
	_ "net/http"
)

type AddIdentificationOpt struct {
}

func (c *Client) AddIdentification(opt *AddIdentificationOpt) error {
	return nil
}

type UpdateIdentificationOpt struct {
}

func (c *Client) UpdateIdentification(opt *UpdateIdentificationOpt) error {
	return nil
}

type DeleteIdentificationOpt struct {
}

func (c *Client) DeleteIdentification(opt *DeleteIdentificationOpt) error {
	return nil
}
