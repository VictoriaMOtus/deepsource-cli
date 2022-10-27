package report

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

// makeQuery makes a HTTP query with a specified body and returns the response
func makeQuery(url string, body []byte, bodyMimeType string, skipCertificateVerification bool) ([]byte, error) {
	var resBody []byte
	httpClient := &http.Client{
		Timeout: time.Second * 60,
	}

	if skipCertificateVerification {
		// Create a custom HTTP Transport for skipping verification of SSL certificates
		// if `--skip-verify` flag is passed.
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
		httpClient.Transport = tr
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", bodyMimeType)
	res, err := httpClient.Do(req)
	if err != nil {
		return resBody, err
	}
	defer res.Body.Close()

	resBody, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return resBody, err
	}

	if res.StatusCode >= http.StatusInternalServerError || res.StatusCode != 200 {
		if resBody != nil {
			return resBody, fmt.Errorf("Server responded with %s: %s", strconv.Itoa(res.StatusCode), string(resBody))
		}
		return resBody, fmt.Errorf("Server responded with %s", strconv.Itoa(res.StatusCode))
	}

	return resBody, nil
}
