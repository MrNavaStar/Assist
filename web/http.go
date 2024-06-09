package web

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// GetJson makes a request to the given url and tries to convert the json response to the given json struct
func GetJson(url string, j interface{}) (err error) {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed request to: %s. Bad status: %s", url, resp.Status)
	}
	if err = json.NewDecoder(resp.Body).Decode(&j); err != nil {
		return err
	}
	return
}
