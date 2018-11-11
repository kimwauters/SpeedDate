// Code generated by goa v2.0.0-wip, DO NOT EDIT.
//
// repository HTTP server encoders and decoders
//
// Command:
// $ goa gen github.com/proepkes/speeddate/usersvc/design

package server

import (
	"context"
	"io"
	"net/http"

	repository "github.com/proepkes/speeddate/usersvc/gen/repository"
	repositoryviews "github.com/proepkes/speeddate/usersvc/gen/repository/views"
	goa "goa.design/goa"
	goahttp "goa.design/goa/http"
)

// EncodeInsertResponse returns an encoder for responses returned by the
// repository insert endpoint.
func EncodeInsertResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, interface{}) error {
	return func(ctx context.Context, w http.ResponseWriter, v interface{}) error {
		res := v.(string)
		enc := encoder(ctx, w)
		body := res
		w.WriteHeader(http.StatusCreated)
		return enc.Encode(body)
	}
}

// DecodeInsertRequest returns a decoder for requests sent to the repository
// insert endpoint.
func DecodeInsertRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	return func(r *http.Request) (interface{}, error) {
		var (
			body InsertRequestBody
			err  error
		)
		err = decoder(r).Decode(&body)
		if err != nil {
			if err == io.EOF {
				return nil, goa.MissingPayloadError()
			}
			return nil, goa.DecodePayloadError(err.Error())
		}
		err = body.Validate()
		if err != nil {
			return nil, err
		}
		payload := NewInsertUser(&body)

		return payload, nil
	}
}

// EncodeDeleteResponse returns an encoder for responses returned by the
// repository delete endpoint.
func EncodeDeleteResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, interface{}) error {
	return func(ctx context.Context, w http.ResponseWriter, v interface{}) error {
		w.WriteHeader(http.StatusNoContent)
		return nil
	}
}

// DecodeDeleteRequest returns a decoder for requests sent to the repository
// delete endpoint.
func DecodeDeleteRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	return func(r *http.Request) (interface{}, error) {
		var (
			id string

			params = mux.Vars(r)
		)
		id = params["id"]
		payload := NewDeletePayload(id)

		return payload, nil
	}
}

// EncodeGetResponse returns an encoder for responses returned by the
// repository get endpoint.
func EncodeGetResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, interface{}) error {
	return func(ctx context.Context, w http.ResponseWriter, v interface{}) error {
		res := v.(*repositoryviews.StoredUser)
		w.Header().Set("goa-view", res.View)
		enc := encoder(ctx, w)
		var body interface{}
		switch res.View {
		case "default", "":
			body = NewGetResponseBody(res.Projected)
		case "tiny":
			body = NewGetResponseBodyTiny(res.Projected)
		}
		w.WriteHeader(http.StatusOK)
		return enc.Encode(body)
	}
}

// DecodeGetRequest returns a decoder for requests sent to the repository get
// endpoint.
func DecodeGetRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	return func(r *http.Request) (interface{}, error) {
		var (
			id   string
			view *string
			err  error

			params = mux.Vars(r)
		)
		id = params["id"]
		viewRaw := r.URL.Query().Get("view")
		if viewRaw != "" {
			view = &viewRaw
		}
		if view != nil {
			if !(*view == "default" || *view == "tiny") {
				err = goa.MergeErrors(err, goa.InvalidEnumValueError("view", *view, []interface{}{"default", "tiny"}))
			}
		}
		if err != nil {
			return nil, err
		}
		payload := NewGetPayload(id, view)

		return payload, nil
	}
}

// EncodeGetError returns an encoder for errors returned by the get repository
// endpoint.
func EncodeGetError(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, error) error {
	encodeError := goahttp.ErrorEncoder(encoder)
	return func(ctx context.Context, w http.ResponseWriter, v error) error {
		en, ok := v.(ErrorNamer)
		if !ok {
			return encodeError(ctx, w, v)
		}
		switch en.ErrorName() {
		case "not_found":
			res := v.(*repository.NotFound)
			enc := encoder(ctx, w)
			body := NewGetNotFoundResponseBody(res)
			w.Header().Set("goa-error", "not_found")
			w.WriteHeader(http.StatusNotFound)
			return enc.Encode(body)
		default:
			return encodeError(ctx, w, v)
		}
	}
}
