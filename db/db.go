package db

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Database struct {
	Address  string
	Username string
	Password string
}

type search struct {
	GT    string `json:"$gt,omitempty"`
	LT    string `json:"$lt,omitempty"`
	Regex string `json:"$regex,omitempty"`
}

type query struct {
	Selector struct {
		ID search `json:"_id"`
	} `json:"selector"`
	Limit int `json:"limit"`
}

type response struct {
	Docs     []device `json:"docs"`
	Bookmark string   `json:"bookmark"`
	Warning  string   `json:"warning"`
}

type device struct {
	ID      string `json:"_id"`
	Address string `json:"address"`
	Roles   []role `json:"roles"`
}

type role struct {
	ID string `json:"_id"`
}

func (d *Database) GetReceiverAddress(roomID string) (string, error) {
	//make request body
	var query query
	query.Limit = 1000
	query.Selector.ID = search{
		Regex: fmt.Sprintf("%s-", roomID),
	}

	body, err := json.Marshal(query)
	if err != nil {
		return "", err
	}

	//make request
	var resp response
	err = d.makeRequest("/devices/_find", "POST", "application/json", body, &resp)
	if err != nil {
		return "", err
	}

	//find receiver address
	for _, dev := range resp.Docs {
		if hasRole("Receiver", dev) {
			return dev.Address, nil
		}
	}
	return "", nil
}

func hasRole(id string, dev device) bool {
	for _, role := range dev.Roles {
		if role.ID == id {
			return true
		}
	}
	return false
}

func (d *Database) makeRequest(path, method, content string, body []byte, respBody interface{}) error {
	url := fmt.Sprintf("%s/%s", d.Address, path)

	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		return err
	}

	if len(d.Username) > 0 {
		req.SetBasicAuth(d.Username, d.Password)
	}
	if len(content) > 0 {
		req.Header.Add("Content-type", content)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("bad status code")
	}

	if respBody != nil {
		err = json.Unmarshal(b, respBody)
		if err != nil {
			return err
		}
	}

	return nil
}
