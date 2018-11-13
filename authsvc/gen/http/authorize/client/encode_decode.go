// Code generated by goa v2.0.0-wip, DO NOT EDIT.
//
// authorize HTTP client encoders and decoders
//
// Command:
// $ goa gen github.com/proepkes/speeddate/authsvc/design

package client

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/url"

	authorize "github.com/proepkes/speeddate/authsvc/gen/authorize"
	goa "goa.design/goa"
	goahttp "goa.design/goa/http"
)

// BuildLoginRequest instantiates a HTTP request object with method and path
// set to call the "authorize" service "login" endpoint
func (c *Client) BuildLoginRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: LoginAuthorizePath()}
	req, err := http.NewRequest("POST", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("authorize", "login", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// EncodeLoginRequest returns an encoder for requests sent to the authorize
// login server.
func EncodeLoginRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	return func(req *http.Request, v interface{}) error {
		p, ok := v.(*authorize.LoginPayload)
		if !ok {
			return goahttp.ErrInvalidType("authorize", "login", "*authorize.LoginPayload", v)
		}
		req.SetBasicAuth(p.Username, p.Password)
		return nil
	}
}

// DecodeLoginResponse returns a decoder for responses returned by the
// authorize login endpoint. restoreBody controls whether the response body
// should be restored after having been read.
func DecodeLoginResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (interface{}, error) {
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
		case http.StatusNoContent:
			var (
				token string
				err   error
			)
			tokenRaw := resp.Header.Get("Authorization")
			if tokenRaw == "" {
				err = goa.MergeErrors(err, goa.MissingFieldError("Authorization", "header"))
			}
			token = tokenRaw
			if err != nil {
				return nil, goahttp.ErrValidationError("authorize", "login", err)
			}
			res := NewLoginResultNoContent(token)
			return res, nil
		default:
			body, _ := ioutil.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("authorize", "login", resp.StatusCode, string(body))
		}
	}
}