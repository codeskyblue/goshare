package gotang

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func HttpGetAsJSON(url string, obj interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(obj)
}

func HttpGetAsString(url string) (content string, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return content, err
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return content, err
	}
	return string(bytes), nil
}
