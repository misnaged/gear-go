//nolint:typecheck
package gear_http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/misnaged/gear-go/config"
	"github.com/misnaged/gear-go/internal/client"
	"github.com/misnaged/gear-go/internal/models"
	"github.com/misnaged/gear-go/pkg/logger"
	"net/http"
	"time"
)

type HttpClient struct {
	client *http.Client
	config *config.Scheme
	id     any
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
	// ---------------------- //
	return httpClient
}
func (cli *HttpClient) SetId(id any) {
	cli.id = id
}
func (cli *HttpClient) PropagateAddress() string {
	return fmt.Sprintf("%s://%s:%d", cli.config.Client.Transport, cli.config.Client.Host, cli.config.Client.Port)
}

func (cli *HttpClient) PostRequest(params any, method string) (*models.RpcGenericResponse, error) {
	address := cli.PropagateAddress()
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
	req, err := http.NewRequest(http.MethodPost, address, bytes.NewReader(body))
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
func (cli *HttpClient) Subscribe(params any, method string) {
	logger.Log().Error("http client DO NOT support subscribe methods")
}
