package gear_ws

import (
	"fmt"
	"github.com/goccy/go-json"
	"github.com/gorilla/websocket"
	"github.com/misnaged/gear-go/config"
	gear_client "github.com/misnaged/gear-go/internal/client"
	"github.com/misnaged/gear-go/internal/models"
	"sync"
)

type wsClient struct {
	config *config.Scheme
	conn   *websocket.Conn
	closed chan struct{}
	mu     sync.RWMutex
}

func (ws *wsClient) PropagateAddress() string {
	return fmt.Sprintf("%s://%s:%v", ws.config.Client.Transport, ws.config.Client.Host, ws.config.Client.Port)
}

func NewWsClient(config *config.Scheme) (gear_client.IClient, error) {
	wsc := &wsClient{
		closed: make(chan struct{}),
		config: config,
	}
	address := wsc.PropagateAddress()
	fmt.Println("address:", address)
	conn, _, err := websocket.DefaultDialer.Dial(address, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to dial to ws: %w", err)
	}

	wsc.conn = conn
	return wsc, nil
}

func (ws *wsClient) PostRequest(params any, method string) (*models.RpcGenericResponse, error) {
	ws.mu.Lock()
	defer ws.mu.Unlock()

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

	err = ws.conn.WriteMessage(websocket.TextMessage, body)
	if err != nil {
		return nil, fmt.Errorf("failed to write to WebSocket: %w", err)
	}

	_, responseData, err := ws.conn.ReadMessage()
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
