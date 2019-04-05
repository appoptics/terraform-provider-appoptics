package appoptics

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

type LegacyClient interface {
	Post(batch *MeasurementsBatch) error
}

type SimpleClient struct {
	httpClient *http.Client
	URL        string
	Token      string
}

func NewLegacyClient(url, token string) LegacyClient {
	return &SimpleClient{
		URL:   url,
		Token: token,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				MaxIdleConnsPerHost: 4,
				IdleConnTimeout:     30 * time.Second,
			},
		},
	}
}

// Used for the collector's internal metrics
func (c *SimpleClient) Post(batch *MeasurementsBatch) error {
	if c.Token == "" {
		return errors.New("AppOptics httpClient not authenticated")
	}
	return c.post(batch, c.Token)
}

func (c *SimpleClient) post(batch *MeasurementsBatch, token string) error {
	json, err := json.Marshal(batch)
	if err != nil {
		log.Error("Error marshaling AppOptics measurements", "err", err)
		return err
	}

	log.Debug("POSTing measurements to AppOptics", "body", string(json))
	req, err := http.NewRequest("POST", c.URL, bytes.NewBuffer(json))
	if err != nil {
		log.Error("Error POSTing measurements to AppOptics", "err", err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(token, "")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Error("Error reading response to AppOptics measurements request", "err", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Error("Error reading response body from AppOptics", "statusCode", resp.StatusCode, "err", err)
			return err
		}

		log.Error("Error POSTing measurements to AppOptics", "statusCode", resp.StatusCode, "respBody", string(body))
		return ErrBadStatus
	}

	log.Debug("Finished uploading AppOptics measurements")

	return nil
}
