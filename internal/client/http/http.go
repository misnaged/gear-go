package gear_http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/misnaged/gear-go/config"
	"github.com/misnaged/gear-go/internal/client"
	"github.com/misnaged/gear-go/internal/models"
	"net/http"
	"time"
)

type HttpClient struct {
	client *http.Client
	config *config.Scheme
}

func NewHttpClient(timeout time.Duration, config *config.Scheme) gear_client.IClient {
	c := &http.Client{
		Timeout: timeout,
	}
	return &HttpClient{client: c, config: config}

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

	defer resp.Body.Close()
	var respRPC *models.RpcGenericResponse
	if err = json.NewDecoder(resp.Body).Decode(&respRPC); err != nil {
		return nil, fmt.Errorf("failed to read all bytes: %w", err)
	}
	return respRPC, nil
}
