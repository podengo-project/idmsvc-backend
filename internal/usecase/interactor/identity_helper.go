package interactor

import (
	"fmt"

	b64 "encoding/base64"
	"encoding/json"

	"github.com/redhatinsights/platform-go-middlewares/identity"
)

func DecodeIdentity(data string) (*identity.Identity, error) {
	if data == "" {
		return nil, fmt.Errorf("X-Rh-Identity content cannot be an empty string")
	}
	bytes, err := b64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}
	identity := &identity.Identity{}
	err = json.Unmarshal(bytes, identity)
	if err != nil {
		return nil, err
	}
	return identity, nil
}
