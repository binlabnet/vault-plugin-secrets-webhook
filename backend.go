package webhook

import (
	"context"
	"strings"

	"sync"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

// New returns a new backend as an interface. This func
// is only necessary for builtin backend plugins.
func New() (interface{}, error) {
	return Backend(), nil
}

// Factory returns a new backend as logical.Backend.
func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	b := Backend()
	if err := b.Setup(ctx, conf); err != nil {
		return nil, err
	}
	return b, nil
}

// FactoryType is a wrapper func that allows the Factory func to specify
// the backend type for the mock backend plugin instance.
func FactoryType(backendType logical.BackendType) logical.Factory {
	return func(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
		b := Backend()
		b.BackendType = backendType
		if err := b.Setup(ctx, conf); err != nil {
			return nil, err
		}
		return b, nil
	}
}

// Backend returns a private embedded struct of framework.Backend.
func Backend() *backend {
	var b backend
	b.Backend = &framework.Backend{
		Help: strings.TrimSpace(backendHelp),

		PathsSpecial: &logical.Paths{
			Unauthenticated: []string{
				"keys/jws/certificate",
				"keys/client/certificate",
			},

			LocalStorage: []string{
				"verify/",
			},
			SealWrapStorage: []string{
				"config/keys/jws",
				"config/keys/client",
			},
		},

		Paths: []*framework.Path{
			pathConfigJws(&b),
			pathFetchJwsCertificate(&b),
			//pathConfigClient(&b),
			pathConfigDestination(&b),
			pathConfigDestinations(&b),
			pathDestination(&b),
			pathVerify(&b),
			//pathFetchClientCertificate(&b),
		},

		//Secrets:     []*framework.Secret{},
		//Invalidate:  b.invalidate,
		BackendType: logical.TypeLogical,
	}

	return &b
}

type backend struct {
	*framework.Backend
	Lock sync.RWMutex
}

const backendHelp = `
The webhook backend sends signed HTTP requests to other services, allowing Vault to perform the AAA
work for privileged requests.
`
