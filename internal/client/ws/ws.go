package gear_ws

import (
	"context"
	"errors"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/gorilla/websocket"
	"github.com/misnaged/gear-go/config"

	//nolint:typecheck
	gear_client "github.com/misnaged/gear-go/internal/client"

	"github.com/misnaged/gear-go/internal/models"
	"github.com/misnaged/gear-go/pkg/logger"
)

type wsClient struct {
	address        string
	config         *config.Scheme
	closed         chan struct{}
	cancel         context.CancelFunc
	sem            chan struct{}
	responsePool   map[gear_client.ResponseType]chan *models.SubscriptionResponse
	connectionPool map[gear_client.ResponseType]*websocket.Conn
	id             any
	responseTypes  []gear_client.ResponseType
}

func NewWsClient(config *config.Scheme) (gear_client.IWsClient, error) {
	wsc := &wsClient{
		closed: make(chan struct{}),
		config: config,
		sem:    make(chan struct{}, 1),
	}
	wsc.propagateAddress()
	if wsc.id == nil {
		wsc.id = "1"
	}

	// --------------------- //

	return wsc, nil
}

func (ws *wsClient) newResponseType(typeName string) error {
	if len(ws.responseTypes) > 0 {
		for _, t := range ws.responseTypes {
			if string(t) == typeName {
				return fmt.Errorf(" gear.wsClient.NewResponseType failed: type name %s already exists", typeName)
			}
		}
	}
	ws.responseTypes = append(ws.responseTypes, gear_client.ResponseType(typeName))
	return nil
}
func (ws *wsClient) propagateAddress() {
	ws.address = fmt.Sprintf("%s://%s:%d", ws.config.Client.Transport, ws.config.Client.Host, ws.config.Client.Port)
}

func (ws *wsClient) SetId(id any) {
	ws.id = id
}

func (ws *wsClient) readLoop(respType gear_client.ResponseType) {
	for {
		select {
		case <-ws.closed:
			logger.Log().Info("client disconnected")
			return
		default:
			_, message, err := ws.connectionPool[respType].ReadMessage()
			if err != nil {
				logger.Log().Errorf("wsClient.readLoop: read message failed: %v", err)
				return
			}
			var response models.SubscriptionResponse
			if err = json.Unmarshal(message, &response); err != nil {
				logger.Log().Errorf("failed to unmarshal message: %v, body: %s", err, message)
				return
			}

			ws.responsePool[respType] <- &response
		}
	}
}
func (ws *wsClient) Acquire() {
	ws.sem <- struct{}{}
}
func (ws *wsClient) Release() {
	<-ws.sem
}

func (ws *wsClient) AddResponseTypesAndMakeWsConnectionsPool(responseTypes ...string) error {
	if ws.config.Subscriptions.Enabled {

		ws.responsePool = make(map[gear_client.ResponseType]chan *models.SubscriptionResponse)
		ws.connectionPool = make(map[gear_client.ResponseType]*websocket.Conn)

		if responseTypes == nil && len(responseTypes) <= 0 {
			return fmt.Errorf(" gear.wsClient.NewWsClient: no response types provided")
		}

		for _, r := range responseTypes {
			if err := ws.newResponseType(r); err != nil {
				return fmt.Errorf("%w", err)
			}
		}
		for _, rType := range ws.responseTypes {
			err := ws.newConnection(rType)
			if err != nil {
				return fmt.Errorf(":%w", err)
			}
			ws.responsePool[rType] = make(chan *models.SubscriptionResponse, ws.config.Subscriptions.BuffSize)

		}

		return nil
	}
	return errors.New("subscriptions not enabled in config")
}

func (ws *wsClient) Cancel() {
	ws.cancel()
}
func (ws *wsClient) newConnection(respType gear_client.ResponseType) error {
	conn, _, err := websocket.DefaultDialer.Dial(ws.address, nil) // todo: switch to dialer with context
	if err != nil {
		return fmt.Errorf("failed to dial websocket:%w", err)
	}
	ws.connectionPool[respType] = conn
	return nil
}

func (ws *wsClient) NewSubscriptionFunc(method string, params any, responseType gear_client.ResponseType) (chan *models.SubscriptionResponse, error) {

	if ws.responseTypes == nil {
		return nil, fmt.Errorf(" gear.wsClient.NewSubscriptionFunc: no response types provided")
	}

	for _, rType := range ws.responseTypes {
		if responseType == rType {
			go ws.readLoop(responseType)
			return ws.subscribe(params, method, responseType)
		}
	}
	return nil, fmt.Errorf("gear.wsClient.NewSubscriptionFunc: response type %s not supported", responseType)

}

func (ws *wsClient) CloseAllConnection() error {
	for _, conn := range ws.connectionPool {
		err := conn.Close()
		if err != nil {
			return fmt.Errorf("failed to close connection:%w", err)
		}
	}
	return nil
}

func (ws *wsClient) CloseChannelByResponseType(respType gear_client.ResponseType) {
	if ws.responsePool[respType] == nil {
		fmt.Println(" gear.wsClient.CloseChannelByResponseType: no response types provided", respType)
	}
	close(ws.responsePool[respType])
}

func (ws *wsClient) subscribe(params any, method string, responseType gear_client.ResponseType) (chan *models.SubscriptionResponse, error) {
	rpcRequest := &models.RpcGenericRequest{
		Jsonrpc: "2.0",
		Id:      "1",
		Method:  method,
		Params:  params,
	}
	body, err := json.Marshal(rpcRequest)
	if err != nil {
		return nil, err
	}
	ws.Acquire()
	err = ws.connectionPool[responseType].WriteMessage(websocket.TextMessage, body)
	ws.Release()
	if err != nil {
		return nil, err
	}

	return ws.responsePool[responseType], nil
}

func (ws *wsClient) PostRequest(params any, method string) (*models.RpcGenericResponse, error) {
	conn, _, err := websocket.DefaultDialer.Dial(ws.address, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to dial websocket:%w", err)
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
		return nil, fmt.Errorf("marshal json rpc request body failed: %w", err)
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
