// Code generated by goa v2.0.0-wip, DO NOT EDIT.
//
// authorize endpoints
//
// Command:
// $ goa gen github.com/proepkes/speeddate/authsvc/design

package authorize

import (
	"context"

	goa "goa.design/goa"
	"goa.design/goa/security"
)

// Endpoints wraps the "authorize" service endpoints.
type Endpoints struct {
	Login goa.Endpoint
}

// NewEndpoints wraps the methods of the "authorize" service with endpoints.
func NewEndpoints(s Service) *Endpoints {
	// Casting service to Auther interface
	a := s.(Auther)
	return &Endpoints{
		Login: NewLoginEndpoint(s, a.BasicAuth),
	}
}

// Use applies the given middleware to all the "authorize" service endpoints.
func (e *Endpoints) Use(m func(goa.Endpoint) goa.Endpoint) {
	e.Login = m(e.Login)
}

// NewLoginEndpoint returns an endpoint function that calls the method "login"
// of service "authorize".
func NewLoginEndpoint(s Service, authBasicFn security.AuthBasicFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*LoginPayload)
		var err error
		sc := security.BasicScheme{
			Name: "basic",
		}
		ctx, err = authBasicFn(ctx, p.Username, p.Password, &sc)
		if err != nil {
			return nil, err
		}
		return s.Login(ctx, p)
	}
}