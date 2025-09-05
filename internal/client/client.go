//nolint:typecheck
package gear_client

import (
	"github.com/misnaged/gear-go/internal/models"
)

type IClient interface {
	PostRequest(params any, method string) (*models.RpcGenericResponse, error)
	SetId(id any)
	PropagateAddress() string
}

type IWsClient interface {
	IClient
	Subscribe(params any, method string) (<-chan *models.SubscriptionResponse, error)
	Cancel()
}
