//nolint:typecheck
package gear_client

import (
	"github.com/misnaged/gear-go/internal/models"
)

type ResponseType string

type IClient interface {
	PostRequest(params any, method string) (*models.RpcGenericResponse, error)
	SetId(id any)
}

type IWsClient interface {
	IClient
	AddResponseTypesAndMakeWsConnectionsPool(responseTypes ...string) error
	NewSubscriptionFunc(method string, params any, responseType ResponseType) (chan *models.SubscriptionResponse, error)
	CloseAllConnection() error
	CloseChannelByResponseType(respType ResponseType)
	Cancel()
}
