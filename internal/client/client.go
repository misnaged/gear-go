package gear_client

import (
	"github.com/misnaged/gear-go/internal/models"
)

type IClient interface {
	PostRequest(params any, method string) (*models.RpcGenericResponse, error)
	SetId(id any)
	PropagateAddress() string
	Subscribe(params any, method string)
}
