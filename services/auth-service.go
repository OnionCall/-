package services

import "net/http"

func authorize(req *http.Request, contentType string) (*http.Response, error) {
	client := &http.Client{}
	req.SetBasicAuth("admin", "hopefullywedontneedthispasswordlongterm")

	if req.Method == "POST" {
		req.Header.Set("Content-Type", contentType)	
	}
	resp, err := client.Do(req)

	return resp, err
}