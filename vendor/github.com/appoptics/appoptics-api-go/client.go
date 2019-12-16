package appoptics

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"regexp"

	"fmt"

	"time"

	"io/ioutil"

	log "github.com/sirupsen/logrus"
)


const (
	// MeasurementPostMaxBatchSize defines the max number of Measurements to send to the API at once
	MeasurementPostMaxBatchSize = 1000
	// DefaultPersistenceErrorLimit sets the number of errors that will be allowed before persistence shuts down
	DefaultPersistenceErrorLimit = 5
	defaultBaseURL               = "https://api.appoptics.com/v1/"
	defaultMediaType             = "application/json"
	clientIdentifier             = "appoptics-api-go"
)

var (
	regexpIllegalNameChars = regexp.MustCompile("[^A-Za-z0-9.:_-]") // from https://www.AppOptics.com/docs/api/#measurements
	// ErrBadStatus is returned if the AppOptics API returns a non-200 error code.
	ErrBadStatus = errors.New("Received non-OK status from AppOptics POST")
	client = &http.Client{
		Timeout: 30 * time.Second,
	}
)

// ServiceAccessor defines an interface for talking to via domain-specific service constructs
type ServiceAccessor interface {
	AlertsService() AlertsCommunicator
	AnnotationsService() AnnotationsCommunicator
	ApiTokensService() ApiTokensCommunicator
	ChartsService() ChartsCommunicator
	JobsService() JobsCommunicator
	MeasurementsService() MeasurementsCommunicator
	MetricsService() MetricsCommunicator
	ServicesService() ServicesCommunicator
	SnapshotsService() SnapshotsCommunicator
	SpacesService() SpacesCommunicator
}

// ErrorResponse represents the response body returned when the API reports an error
type ErrorResponse struct {
	// Errors holds the error information from the API
	Errors   interface{} `json:"errors"`
	Response *http.Response
}

// QueryInfo holds pagination information coming from list actions
type QueryInfo struct {
	Found  int `json:"found,omitempty"`
	Length int `json:"length,omitempty"`
	Offset int `json:"offset,omitempty"`
	Total  int `json:"total,omitempty"`
}

// RequestErrorMessage represents the error schema for request errors
// TODO: add API reference URLs here
type RequestErrorMessage map[string][]string

// ParamErrorMessage represents the error schema for param errors
// TODO: add API reference URLs here
type ParamErrorMessage []map[string]string

// Client implements ServiceAccessor
type Client struct {
	baseURL                 *url.URL
	httpClient              httpClient
	token                   string
	alertsService           AlertsCommunicator
	annotationsService      AnnotationsCommunicator
	apiTokensService        ApiTokensCommunicator
	chartsService           ChartsCommunicator
	jobsService             JobsCommunicator
	measurementsService     MeasurementsCommunicator
	metricsService          MetricsCommunicator
	spacesService           SpacesCommunicator
	snapshotsService        SnapshotsCommunicator
	servicesService         ServicesCommunicator
	callerUserAgentFragment string
	debugMode               bool
}

// httpClient defines the http.Client method used by Client.
type httpClient interface {
	Do(*http.Request) (*http.Response, error)
}

// ClientOption provides functional option-setting behavior
type ClientOption func(*Client) error

// NewClient returns a new AppOptics API client. Optional arguments UserAgentClientOption and BaseURLClientOption can be provided.
func NewClient(token string, opts ...func(*Client) error) *Client {
	baseURL, _ := url.Parse(defaultBaseURL)
	c := &Client{
		token:   token,
		baseURL: baseURL,
		httpClient: client,
	}

	c.alertsService = NewAlertsService(c)
	c.annotationsService = NewAnnotationsService(c)
	c.apiTokensService = NewApiTokensService(c)
	c.chartsService = NewChartsService(c)
	c.jobsService = NewJobsService(c)
	c.measurementsService = NewMeasurementsService(c)
	c.metricsService = NewMetricsService(c)
	c.servicesService = NewServiceService(c)
	c.snapshotsService = NewSnapshotsService(c)
	c.spacesService = NewSpacesService(c)

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// NewRequest standardizes the request being sent
func (c *Client) NewRequest(method, path string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	requestURL := c.baseURL.ResolveReference(rel)

	var buffer io.ReadWriter

	if body != nil {
		buffer = &bytes.Buffer{}
		encodeErr := json.NewEncoder(buffer).Encode(body)
		if encodeErr != nil {
			log.Println(encodeErr)
		}
	}
	req, err := http.NewRequest(method, requestURL.String(), buffer)

	if err != nil {
		return nil, err
	}

	req.SetBasicAuth("token", c.token)
	req.Header.Set("Accept", defaultMediaType)
	req.Header.Set("Content-Type", defaultMediaType)
	req.Header.Set("User-Agent", c.completeUserAgentString())

	return req, nil
}

// UserAgentClientOption is a config function allowing setting of the User-Agent header in requests
func UserAgentClientOption(userAgentString string) ClientOption {
	return func(c *Client) error {
		c.callerUserAgentFragment = userAgentString
		return nil
	}
}

// BaseURLClientOption is a config function allowing setting of the base URL the API is on
func BaseURLClientOption(urlString string) ClientOption {
	return func(c *Client) error {
		var altURL *url.URL
		var err error
		if altURL, err = url.Parse(urlString); err != nil {
			return err
		}
		c.baseURL = altURL
		return nil
	}
}

// SetDebugMode sets the debugMode struct member to true
func SetDebugMode() ClientOption {
	return func(c *Client) error {
		c.debugMode = true
		return nil
	}
}

// SetHTTPClient allows the user to provide a custom http.Client configuration
func SetHTTPClient(client *http.Client) ClientOption {
	return func(c *Client) error {
		c.httpClient = client
		return nil
	}
}

// AlertsService represents the subset of the API that deals with Alerts
func (c *Client) AlertsService() AlertsCommunicator {
	return c.alertsService
}

// AnnotationsService represents the subset of the API that deals with Annotations
func (c *Client) AnnotationsService() AnnotationsCommunicator {
	return c.annotationsService
}

func (c *Client) ApiTokensService() ApiTokensCommunicator {
	return c.apiTokensService
}

// ChartsService represents the subset of the API that deals with Charts
func (c *Client) ChartsService() ChartsCommunicator {
	return c.chartsService
}

// JobsService represents the subset of the API that deals with Jobs
func (c *Client) JobsService() JobsCommunicator {
	return c.jobsService
}

// MeasurementsService represents the subset of the API that deals with Measurements
func (c *Client) MeasurementsService() MeasurementsCommunicator {
	return c.measurementsService
}

// MetricsService represents the subset of the API that deals with Metrics
func (c *Client) MetricsService() MetricsCommunicator {
	return c.metricsService
}

// ServicesService represents the subset of the API that deals with Services
func (c *Client) ServicesService() ServicesCommunicator {
	return c.servicesService
}

// SnapshotsService represents the subset of the API that deals with Snapshots
func (c *Client) SnapshotsService() SnapshotsCommunicator {
	return c.snapshotsService
}

// SpacesService represents the subset of the API that deals with Spaces
func (c *Client) SpacesService() SpacesCommunicator {
	return c.spacesService
}

// Error makes ErrorResponse satisfy the error interface and can be used to serialize error responses back to the httpClient
func (e *ErrorResponse) Error() string {
	errorData, _ := json.Marshal(e)
	return string(errorData)
}

// DefaultPaginationParameters provides a *PaginationParameters with minimum required fields
func (c *Client) DefaultPaginationParameters(length int) *PaginationParameters {
	return &PaginationParameters{
		Length:  length,
		Sort:    "asc",
		Orderby: "name",
	}
}

// Do performs the HTTP request on the wire, taking an optional second parameter for containing a response
func (c *Client) Do(req *http.Request, respData interface{}) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)

	// error in performing request
	if err != nil {
		return resp, err
	}

	if c.debugMode {
		dumpResponse(resp)
	}
	// request response contains an error
	if err = checkError(resp); err != nil {
		return resp, err
	}

	defer resp.Body.Close()
	if respData != nil {
		err = json.NewDecoder(resp.Body).Decode(respData)
	}

	return resp, err
}

// completeUserAgentString returns the string that will be placed in the User-Agent header.
// It ensures that any caller-set string has the client name and version appended to it.
func (c *Client) completeUserAgentString() string {
	if c.callerUserAgentFragment == "" {
		return clientVersionString()
	}
	return fmt.Sprintf("%s:%s", c.callerUserAgentFragment, clientVersionString())
}

// clientVersionString returns the canonical name-and-version string
func clientVersionString() string {
	return fmt.Sprintf("%s", clientIdentifier)
}

// checkError creates an ErrorResponse from the http.Response.Body, if there is one
func checkError(resp *http.Response) error {
	errResponse := &ErrorResponse{}
	if resp.StatusCode >= 400 {
		errResponse.Response = resp
		if resp.ContentLength != 0 {
			decoder := json.NewDecoder(resp.Body)
			err := decoder.Decode(errResponse)

			if err != nil {
				return err
			}
			log.Debugf("error: %+v\n", errResponse)
			return errResponse
		}
		return errResponse
	}
	return nil
}

// dumpResponse is a debugging function which dumps the HTTP response to stdout
func dumpResponse(resp *http.Response) {
	fmt.Printf("response status: %s\n", resp.Status)
	if resp.ContentLength != 0 {
		if respBytes, err := ioutil.ReadAll(resp.Body); err != nil {
			log.Printf("Error reading body: %s", err)
			return
		} else {
			resp.Body.Close()
			resp.Body = ioutil.NopCloser(bytes.NewBuffer(respBytes))
			fmt.Printf("response body: %s\n\n", string(respBytes))
		}
	}
}
