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
)

func (ws *wsClient) PropagateAddress() string {
	return fmt.Sprintf("%s://%s:%d", ws.config.Client.Transport, ws.config.Client.Host, ws.config.Client.Port)
}

func (ws *wsClient) SetId(id any) {
	ws.id = id
}

type wsClient struct {
	config       *config.Scheme
	closed       chan struct{}
	cancel       context.CancelFunc
	mu           sync.RWMutex
	dealer       *websocket.Dialer
	conn         *websocket.Conn
	sem          chan struct{}
	subscribes   sync.Map
	responseChan chan *models.SubscriptionResponse
	id           any
}

func (ws *wsClient) readLoop() {
	for {
		select {
		case <-ws.closed:
			logger.Log().Info("client disconnected")
			break
		default:

			_, message, err := ws.conn.ReadMessage()
			if err != nil {
				logger.Log().Errorf("wsClient.readLoop: read message failed: %v", err)
				return
			}
			var response models.SubscriptionResponse
			if err = json.Unmarshal(message, &response); err != nil {
				logger.Log().Errorf("failed to unmarshal message: %v, body: %s", err, message)
				return
			}

			ws.responseChan <- &response
		}
	}
}
func (ws *wsClient) Acquire() {
	ws.sem <- struct{}{}
}
func (ws *wsClient) Release() {
	<-ws.sem
}
func (ws *wsClient) WriteResponse(resp *models.SubscriptionResponse) {
	ws.responseChan <- resp
}
func NewWsClient(config *config.Scheme) (gear_client.IWsClient, error) {
	wsc := &wsClient{
		closed:       make(chan struct{}),
		config:       config,
		sem:          make(chan struct{}, 1),
		responseChan: make(chan *models.SubscriptionResponse, 100),
	}
	address := wsc.PropagateAddress()
	//pCtx := context.Background()
	//ctx, cancel := context.WithTimeout(pCtx, time.Second*1)
	//defer cancel()
	conn, _, err := websocket.DefaultDialer.Dial(address, nil)
	if err != nil {
		return nil, fmt.Errorf("dial websocket failed: %v", err)
	}
	wsc.conn = conn

	// Setting default id //
	if wsc.id == nil {
		wsc.id = "1"
	}
	// --------------------- //
	go wsc.readLoop()

	return wsc, nil
}
func (ws *wsClient) Cancel() {
	ws.cancel()
}
func (ws *wsClient) Subscribe(params any, method string) (<-chan *models.SubscriptionResponse, error) {

	respChan := make(chan *models.SubscriptionResponse, 1)
	rpcRequest := &models.RpcGenericRequest{
		Jsonrpc: "2.0",
		Id:      "1",
		Method:  method,
		Params:  params,
	}
	body, err := json.Marshal(rpcRequest)
	if err != nil {
		return nil, fmt.Errorf("marshal json rpc request body failed: %v", err)
	}
	ws.Acquire()
	err = ws.conn.WriteMessage(websocket.TextMessage, body)
	ws.Release()
	if err != nil {
		return nil, fmt.Errorf("send subscription request failed: %v", err)
	}
	go func() {
		for resp := range ws.responseChan {
			if resp.Result != nil {
				if resp.Result.(string) != "" {
					ws.subscribes.Store(resp.Result.(string), respChan)
				}
			}
			respChan <- resp
		}
	}()

	return respChan, nil
}
func aa() {
	//if resp.Params != nil {
	//	e, err := models.GetFieldFromAny("result", resp.Params)
	//	if err != nil {
	//		logger.Log().Errorf("failed to get field: %v", err)
	//		break
	//	}
	//	switch e.(type) {
	//	case map[string]any:
	//		m, ok := e.(map[string]any)
	//		if !ok {
	//			logger.Log().Error("not a map[string]any")
	//			break
	//		}
	//		subHash, err := models.GetFieldFromAny("subscription", resp.Params)
	//		if err != nil {
	//			logger.Log().Errorf("failed to get field: %v", err)
	//			break
	//		}
	//		if subHash.(string) == "" {
	//			logger.Log().Errorf("failed to get field: %v", err)
	//			break
	//		}
	//		if m["finalized"] != nil {
	//			ws.subscribes.Delete(subHash.(string))
	//			close(respChan)
	//			return
	//		}
	//	}
	//}

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

	body, err := json.Marshal(rpcRequest)
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
