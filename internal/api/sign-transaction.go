package gear_api

import (
	"errors"
	"fmt"
	scalecodec "github.com/itering/scale.go"
	"github.com/itering/scale.go/types"
	gear_utils "github.com/misnaged/gear-go/internal/utils"
	"github.com/misnaged/substrate-api-rpc/keyring"
	rpcModels "github.com/misnaged/substrate-api-rpc/model"
	"github.com/misnaged/substrate-api-rpc/rpc"
)

func (api *Api) SignTransaction(moduleName, callName string, kr keyring.IKeyRing, params []scalecodec.ExtrinsicParam) (string, error) {
	if kr == nil {
		return "", errors.New("failed to sign transaction: signer keyring is nil")
	}
	if params == nil {
		return "", errors.New("failed to sign transaction: params is nil")
	}
	if err := api.MetadataCheck(); err != nil {
		return "", fmt.Errorf("failed to sign transaction: %w", err)
	}

	genesisHash, err := api.getChainGetBlockHash()
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}
	version, err := api.getStateGetRuntimeVersion()
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	opt := &types.ScaleDecoderOption{Metadata: api.metadata, Spec: -1}
	callIndex := gear_utils.GetCallLookupIndexByModuleAndCallNames(api.metadata, moduleName, callName)

	resp, err := api.gearRpc.SystemAccountNextIndex(kr.PublicKey())
	if err != nil {
		return "", fmt.Errorf("failed to send SystemAccountNextIndex request: %w", err)
	}
	api.customTx = rpc.NewCustomTransaction(
		callIndex,
		genesisHash,
		int(resp.Result.(float64)),
		version,
		api.metadata,
		kr,
		opt,
		params,
	)
	signed, err := api.customTx.SignTransactionCustom() //TODO: Era is ALWAYS immortal. Need to change!
	if err != nil {
		return "", fmt.Errorf("failed to sign transaction: %w", err)
	}
	return signed, nil
}

func (api *Api) getStateGetRuntimeVersion() (*rpcModels.RuntimeVersion, error) {
	runtimeVersion, err := api.gearRpc.StateGetRuntimeVersionLatest()
	if err != nil {
		return nil, fmt.Errorf("request state_getRuntimeVersion failed: %v", err)
	}
	rtm := &rpcModels.RuntimeVersion{}

	switch runtimeVersion.Result.(type) {
	case map[string]any:
		for key, val := range runtimeVersion.Result.(map[string]any) {
			switch key {
			//todo panic check! do not forget!
			case "apis":
				if val.([]any) == nil {
					return nil, errors.New("apis has a wrong type")
				}
				rtm.Apis = val.([]any)
			case "implName":
				rtm.ImplName = val.(string)
			case "implVersion":
				rtm.ImplVersion = int(val.(float64))
			case "specName":
				rtm.SpecName = val.(string)
			case "specVersion":
				rtm.SpecVersion = int(val.(float64))
			case "transactionVersion":
				rtm.TransactionVersion = int(val.(float64))
			}
		}
		return rtm, nil

	default:
		return nil, errors.New("unknown genesis hash type")
	}
}
func (api *Api) getChainGetBlockHash() (string, error) {

	genesisHash, err := api.gearRpc.ChainGetBlockHash(0)
	if err != nil {
		return "", fmt.Errorf("request chain_getBlockHash failed: %v", err)
	}
	switch genesisHash.Result.(type) {
	case string:
		return genesisHash.Result.(string), nil
	default:
		fmt.Printf("%T\n", genesisHash.Result)
		return "", errors.New("genesisHash is not string")
	}
}

func ParamFromTypes(fromTypes *types.ExtrinsicParam) *scalecodec.ExtrinsicParam {
	return &scalecodec.ExtrinsicParam{
		Name:     fromTypes.Name,
		Type:     fromTypes.Type,
		TypeName: fromTypes.Type,
		Value:    fromTypes.Value,
	}
}
func ParamsFromTypes(fromTypes *types.ExtrinsicParam) []scalecodec.ExtrinsicParam {
	var params []scalecodec.ExtrinsicParam
	params = append(params, *ParamFromTypes(fromTypes))
	return params
}
