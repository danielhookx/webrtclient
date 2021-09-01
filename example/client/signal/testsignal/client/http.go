package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func Put(url string, name string, sdp string) error {
	c := http.Client{
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       0,
	}
	var params = struct {
		Name string `json:"name"`
		SDP  string `json:"sdp"`
	}{
		Name: name,
		SDP:  sdp,
	}
	data, err := json.Marshal(params)
	if err != nil {
		return err
	}
	resp, err := c.Post(url, "application/json", bytes.NewReader(data))
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error http code: %d", resp.StatusCode)
	}
	return nil
}
