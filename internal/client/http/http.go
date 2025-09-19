package gear_http

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/misnaged/gear-go/config"
	//nolint:typecheck
	gear_client "github.com/misnaged/gear-go/internal/client"
	"github.com/misnaged/gear-go/internal/models"
	"net/http"
	"time"
)

type HttpClient struct {
	address string
	client  *http.Client
	config  *config.Scheme
	id      any
}

func (cli *HttpClient) propagateAddress() {
	cli.address = fmt.Sprintf("%s://%s:%d", cli.config.Client.Transport, cli.config.Client.Host, cli.config.Client.Port)
}
func NewHttpClient(timeout time.Duration, config *config.Scheme) gear_client.IClient {
	c := &http.Client{
		Timeout: timeout,
	}
	httpClient := &HttpClient{client: c, config: config}
	// setting id for `Id` json-rpc field
	if httpClient.id == nil {
		httpClient.id = "1"
	}
	httpClient.propagateAddress()
	// ---------------------- //
	return httpClient
}
func (cli *HttpClient) SetId(id any) {
	cli.id = id
}

func (cli *HttpClient) PostRequest(params any, method string) (*models.RpcGenericResponse, error) {
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

	req, err := http.NewRequest(http.MethodPost, cli.address, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("post request has failed: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := cli.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("client Do failed %w", err)
	}
	// nolint:errcheck
	defer resp.Body.Close()

	var respRPC *models.RpcGenericResponse
	if err = json.NewDecoder(resp.Body).Decode(&respRPC); err != nil {
		return nil, fmt.Errorf("failed to read all bytes: %w", err)
	}
	// TODO: needs better and separate error handling
	if respRPC.Error != nil {
		errorMsg := fmt.Sprintf("response for method: %s has failed due to: %v", method, respRPC.Error)
		return respRPC, errors.New(errorMsg)
	}
	return respRPC, nil
}
