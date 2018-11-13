// Code generated by goa v2.0.0-wip, DO NOT EDIT.
//
// health HTTP client encoders and decoders
//
// Command:
// $ goa gen github.com/proepkes/speeddate/usersvc/design

package client

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/url"

	goahttp "goa.design/goa/http"
)

// BuildCheckHealthRequest instantiates a HTTP request object with method and
// path set to call the "health" service "checkHealth" endpoint
func (c *Client) BuildCheckHealthRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: CheckHealthHealthPath()}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("health", "checkHealth", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// DecodeCheckHealthResponse returns a decoder for responses returned by the
// health checkHealth endpoint. restoreBody controls whether the response body
// should be restored after having been read.
func DecodeCheckHealthResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (interface{}, error) {
	return func(resp *http.Response) (interface{}, error) {
		if restoreBody {
			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			resp.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			defer func() {
				resp.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			}()
		} else {
			defer resp.Body.Close()
		}
		switch resp.StatusCode {
		case http.StatusOK:
			var (
				body []byte
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("health", "checkHealth", err)
			}
			return body, nil
		default:
			body, _ := ioutil.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("health", "checkHealth", resp.StatusCode, string(body))
		}
	}
}
