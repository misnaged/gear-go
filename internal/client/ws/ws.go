package gear_ws

import (
	"context"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/gorilla/websocket"
	"github.com/misnaged/gear-go/config"

	//nolint:typecheck
	gear_client "github.com/misnaged/gear-go/internal/client"

	"github.com/misnaged/gear-go/internal/models"
	"github.com/misnaged/gear-go/pkg/logger"
	"sync"
	"time"
)

type wsClient struct {
	config *config.Scheme
	closed chan struct{}
	mu     sync.RWMutex
	dealer *websocket.Dialer
	id     any
}

func (ws *wsClient) PropagateAddress() string {
	return fmt.Sprintf("%s://%s:%d", ws.config.Client.Transport, ws.config.Client.Host, ws.config.Client.Port)
}

func (ws *wsClient) SetId(id any) {
	ws.id = id
}
func NewWsClient(config *config.Scheme) (gear_client.IClient, error) {
	wsc := &wsClient{
		closed: make(chan struct{}),
		config: config,
	}
	//address := wsc.PropagateAddress()

	// Setting default id //
	if wsc.id == nil {
		wsc.id = "1"
	}
	// --------------------- //
	wsc.dealer = &websocket.Dialer{}
	return wsc, nil
}

//func (ws *wsClient) unsubscribe(subId string, conn *websocket.Conn) {
//	var params []any
//	params = append(params, subId)
//	rpcRequest := &models.RpcGenericRequest{
//		Jsonrpc: "2.0",
//		Id:      "1",
//		Method:  "author_unwatchExtrinsic",
//		Params:  params,
//	}
//	body, err := rpcRequest.MarshalBody()
//	if err != nil {
//		logger.Log().Errorf("marshal json rpc request body failed: %v", err)
//		return
//	}
//
//	ws.mu.Lock()
//	err = conn.WriteMessage(websocket.TextMessage, body)
//	ws.mu.Unlock()
//}

func (ws *wsClient) Subscribe(params any, method string) {
	pCtx := context.Background()
	ctx, cancel := context.WithTimeout(pCtx, time.Second*5)
	defer cancel()
	address := ws.PropagateAddress()
	conn, _, err := ws.dealer.DialContext(ctx, address, nil)
	if err != nil {
		//return nil, fmt.Errorf("failed to dial to ws: %w", err)
		logger.Log().Infof("failed to dial websocket:%v", err)
		return
	}

	rpcRequest := &models.RpcGenericRequest{
		Jsonrpc: "2.0",
		Id:      "1",
		Method:  method,
		Params:  params,
	}
	body, err := rpcRequest.MarshalBody()
	if err != nil {
		logger.Log().Errorf("marshal json rpc request body failed: %v", err)
		return
	}

	ws.mu.Lock()
	err = conn.WriteMessage(websocket.TextMessage, body)
	ws.mu.Unlock()

	if err != nil {
		logger.Log().Errorf("send subscription request failed: %v", err)
		return
	}
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ws.closed:
				return
			case <-ctx.Done():
				logger.Log().Info("Done!")
				//ws.unsubscribe("1", conn) //TODO: add unsubscribe if neeeded
				return
			default:
				_, message, err := conn.ReadMessage()
				if err != nil {
					logger.Log().Errorf("read message failed: %v", err)
					return
				}
				var response any
				//TODO: add map type (as well as fully functional subscription result parser)
				if err = json.Unmarshal(message, &response); err != nil {
					logger.Log().Errorf("failed to unmarshal message: %v, body: %s", err, message)
					return
				}
				switch response.(type) {
				case map[string]interface{}:
					aaa := response.(map[string]interface{})["params"]
					if aaa != nil {
						bbb := aaa.(map[string]interface{})
						logger.Log().Infof("received message: %v", bbb["subscription"])
					}
				}

			}
		}
	}()
	wg.Wait()

	// nolint:errcheck
	conn.Close()

}

type Resp struct {
	Params struct {
		Subscription string `json:"subscription"`
		Result       string `json:"result"`
	} `json:"params"`
}
type Subscription struct {
	Subscription string `json:"subscription"`
}

func (ws *wsClient) PostRequest(params any, method string) (*models.RpcGenericResponse, error) {
	address := ws.PropagateAddress()
	conn, _, err := ws.dealer.Dial(address, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to dial websocket:%v", err)
	}

	// nolint:errcheck
	defer conn.Close()

	rpcRequest := &models.RpcGenericRequest{
		Jsonrpc: "2.0",
		Id:      "1",
		Method:  method,
		Params:  params,
	}

	body, err := rpcRequest.MarshalBody()
	if err != nil {
		return nil, fmt.Errorf("marshal json rpc request body failed: %v", err)
	}

	err = conn.WriteMessage(websocket.TextMessage, body)
	if err != nil {
		return nil, fmt.Errorf("failed to write to WebSocket: %w", err)
	}

	_, responseData, err := conn.ReadMessage()
	if err != nil {
		return nil, fmt.Errorf("failed to read response from WebSocket: %w", err)
	}

	var resp models.RpcGenericResponse
	err = json.Unmarshal(responseData, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON-RPC response: %w", err)
	}
	return &resp, nil
}
