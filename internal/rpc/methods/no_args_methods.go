/*
 no args methods represents collection of argumentless json-rpc methods

 Please do keep in mind that ALL SUBSCRIBE methods are handled separately
*/

package gear_rpc_method

import (
	"github.com/misnaged/gear-go/internal/models"
	gear_rpc "github.com/misnaged/gear-go/internal/rpc"
)

func (gearRPC *GearRpc) NoArgRpcRequest(method gear_rpc.NoArgsMethod) (*models.RpcGenericResponse, error) {
	return gearRPC.client.PostRequest(nil, method.String())
}
