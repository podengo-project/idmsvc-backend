package middleware

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/redhatinsights/platform-go-middlewares/identity"
	"github.com/rs/zerolog/log"
)

type XRhIValidator interface {
	ValidateXRhIdentity(xrhi *identity.Identity) error
}

func GetXRhIdentityFromRequest(req *http.Request) (*identity.Identity, error) {
	var (
		b64Identity   string
		bytesIdentity []byte
		err           error
		data          *identity.Identity
	)
	if b64Identity = req.Header.Get(headerXRhIdentity); b64Identity == "" {
		return nil, fmt.Errorf(headerXRhIdentity + " header is missing")
	}

	if bytesIdentity, err = base64.StdEncoding.DecodeString(b64Identity); err != nil {
		return nil, fmt.Errorf(headerXRhIdentity + " header is malformed")
	}

	data = &identity.Identity{}
	if err = json.Unmarshal(bytesIdentity, &data); err != nil {
		log.Error().Err(err)
		return nil, fmt.Errorf(headerXRhIdentity + " header is malformed")
	}

	return data, nil
}

func NewAuthenticator(v XRhIValidator) openapi3filter.AuthenticationFunc {
	return func(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
		return Authenticate(v, ctx, input)
	}
}

func Authenticate(v XRhIValidator, ctx context.Context, input *openapi3filter.AuthenticationInput) error {
	var (
		err  error
		data *identity.Identity
	)
	if input.SecuritySchemeName != "ApiKeyAuth" {
		return fmt.Errorf("security scheme %s != 'ApiKeyAuth", input.SecuritySchemeName)
	}
	if data, err = GetXRhIdentityFromRequest(input.RequestValidationInput.Request); err != nil {
		return fmt.Errorf("Retrieving "+headerXRhIdentity+": %s", err.Error())
	}

	if err = v.ValidateXRhIdentity(data); err != nil {
		return fmt.Errorf("No valid " + headerXRhIdentity)
	}

	return nil
}
