package utilsx

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
	"time"
)

type HttpMethod string

const (
	HTTP_METHOD_GET     HttpMethod = "GET"
	HTTP_METHOD_POST    HttpMethod = "POST"
	HTTP_METHOD_PUT     HttpMethod = "PUT"
	HTTP_METHOD_PATCH   HttpMethod = "PATCH"
	HTTP_METHOD_DELETE  HttpMethod = "DELETE"
	HTTP_METHOD_HEAD    HttpMethod = "HEAD"
	HTTP_METHOD_OPTIONS HttpMethod = "OPTIONS"
)

type ApiRequest struct {
	method         string // request method
	schema         string // request schema, default use https
	serviceAddress string // request url or reuqest host
	uri            string // request path

	headers map[string]string // request headers

	query url.Values             // request query parameters
	body  map[string]interface{} // request post body

	timeout time.Duration // request timeout

	apiResponse           *http.Response // request response
	apiResponseStatus     string         // request response status
	apiResponseStatusCode int            // request response status code
	apiResponseData       []byte         // request response data
	apiResponseError      error          // request response error
}

// NewHttpRequest creates a new HTTP request.
//
// It takes a string parameter `address` which represents the URL of the request.
// The function returns an `executableApiRequest` object.
func NewHttpRequest(address string) ExecutableApiRequest {
	Url, _ := url.Parse(address)
	rawQuery, _ := url.ParseQuery(Url.RawQuery)
	// default request scheme is https
	if Url.Scheme == "" {
		Url.Scheme = "https"
	}
	return &ApiRequest{
		schema:         Url.Scheme,
		serviceAddress: address,
		query:          rawQuery,
		body:           make(map[string]interface{}),
		headers:        make(map[string]string),
		timeout:        time.Second * 5,
	}
}

// Do sends an API request and returns the response.
//
// It takes a method of type httpMethod as a parameter and returns an apiResponse
// and an error.
func (r *ApiRequest) Do(method HttpMethod) (apiResponse, error) {
	defer r.printLog()
	postBody, _ := json.Marshal(r.body)
	req, err := r.newRequest(method, bytes.NewReader(postBody))
	if err != nil {
		return nil, err
	}

	httpClient := http.DefaultClient
	httpClient.Timeout = r.timeout

	r.apiResponse, err = httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	r.apiResponseStatus = r.apiResponse.Status
	r.apiResponseStatusCode = r.apiResponse.StatusCode
	r.apiResponseData, err = io.ReadAll(r.apiResponse.Body)
	if err != nil {
		log.Printf("response result read error: %s", err.Error())
		return nil, err
	}
	defer r.apiResponse.Body.Close()

	return r, err
}

func (r *ApiRequest) printLog() {
	jsonBody, _ := json.Marshal(r.body)
	jsonHeaders, _ := json.Marshal(r.headers)
	var logOut = fmt.Sprintf("%s %s Body: %s Headers: %s \n",
		r.method,
		r.GetUrl(),
		string(jsonBody),
		string(jsonHeaders),
	)
	if r.apiResponseStatusCode >= http.StatusOK && r.apiResponseStatusCode < http.StatusBadRequest {
		logOut += "Success ResponseStatus: %s ResponseData: %s"
		log.Printf(logOut, r.apiResponseStatus, string(r.apiResponseData))
		return
	}
	if r.apiResponseStatusCode >= http.StatusBadRequest && r.apiResponseError != nil {
		logOut += "Failed ResponseStatus: %s ResponseData: %s ResponseError: %s"
		log.Printf(logOut, r.apiResponseStatus, string(r.apiResponseData), r.apiResponseError.Error())
		return
	}
	log.Println(logOut)
}

// newRequest creates a new HTTP request with the given method and body.
//
// It takes in a method of type httpMethod and a body of type io.Reader.
// It returns a pointer to an http.Request and an error.
func (r *ApiRequest) newRequest(method HttpMethod, body io.Reader) (*http.Request, error) {
	r.method = string(method)
	req, err := http.NewRequest(r.method, r.GetUrl(), body)
	if err != nil {
		return nil, err
	}
	for key, value := range r.headers {
		req.Header.Set(key, value)
	}
	return req, nil
}

// GetUrl returns the URL string for the API request.
//
// It does not take any parameters.
// It returns a string.
func (r *ApiRequest) GetUrl() string {
	// build url path
	Url, _ := url.Parse(r.serviceAddress)
	Url.Scheme = r.schema
	if len(r.uri) > 0 {
		Url.Path = path.Join(Url.Path, r.uri)
	}
	// query parameters encode
	Url.RawQuery = r.query.Encode()
	return Url.String()
}

type ExecutableApiRequest interface {
	SetUri(uri string) ExecutableApiRequest
	SetHeader(key, value string) ExecutableApiRequest
	SetBody(key string, value interface{}) ExecutableApiRequest
	SetQueryParam(key string, value string) ExecutableApiRequest
	SetTimeout(time.Duration) ExecutableApiRequest
	Do(method HttpMethod) (apiResponse, error)
}

// SetUri sets the URI for the apiRequest.
//
// Parameters:
//
//   - uri: The URI to be set.
//
// Returns:
//   - executableApiRequest: The modified apiRequest struct.
func (r *ApiRequest) SetUri(uri string) ExecutableApiRequest {
	r.uri = uri
	return r
}

// SetHeader sets a header in the apiRequest.
//
// Parameters:
//
//   - key: The key of the header.
//
//   - value: The value of the header.
//
// Returns:
//   - executableApiRequest: The modified apiRequest struct.
func (r *ApiRequest) SetHeader(key, value string) ExecutableApiRequest {
	r.headers[key] = value
	return r
}

// SetBody sets a key-value pair in the postBody field of the apiRequest struct.
//
// Parameters:
//   - key: The key to set in the postBody map.
//   - value: The value to set for the given key in the postBody map.
//
// Returns:
//   - executableApiRequest: The modified apiRequest struct.
func (r *ApiRequest) SetBody(key string, value interface{}) ExecutableApiRequest {
	r.body[key] = value
	return r
}

// SetQueryParam sets a query parameter in the apiRequest.
//
// Parameters:
//   - key: the key of the query parameter.
//   - value: the value of the query parameter.
//
// Returns:
//   - executableApiRequest: The modified apiRequest struct.
func (r *ApiRequest) SetQueryParam(key, value string) ExecutableApiRequest {
	r.query.Set(key, value)
	return r
}

// SetTimeout sets the timeout duration for the API request.
//
// Parameters:
//   - timeout: the duration for the request timeout.
//
// Returns:
//   - executableApiRequest: The modified apiRequest struct.
func (r *ApiRequest) SetTimeout(timeout time.Duration) ExecutableApiRequest {
	r.timeout = timeout
	return r
}

type apiResponse interface {
	Result() ([]byte, error)
	Success() bool
	SuccessResult() ([]byte, error)
}

// Result returns the result of the HTTP response.
//
// It reads the response body and returns it as a byte slice along with any
// error encountered during the process.
//
// Returns:
//   - []byte: The response body as a byte slice.
//   - error: An error if any occurred during the process.
func (r *ApiRequest) Result() ([]byte, error) {
	return r.apiResponseData, r.apiResponseError
}

// Success checks if the HTTP response status code is less than 400.
//
// It returns a boolean value indicating whether the response is considered a success.
func (r *ApiRequest) Success() bool {
	return r.apiResponseStatusCode < http.StatusBadRequest
}

// SuccessResult returns the response body and error from the HTTP request.
//
// It calls the Result() method to get the response body and checks if the
// response status code is greater than or equal to http.StatusBadRequest.
// If it is, it logs the request URL, status code, and response body, and
// returns an error with the response body. Otherwise, it returns the response
// body and nil error.
//
// Returns:
//   - []byte: The response body.
//   - error: An error if the response status code is greater than or equal to
//     http.StatusBadRequest, otherwise nil.
func (r *ApiRequest) SuccessResult() ([]byte, error) {
	if r.apiResponseStatusCode >= http.StatusBadRequest {
		err := errors.New(string(r.apiResponseData))
		return nil, err
	}
	return r.apiResponseData, r.apiResponseError
}
