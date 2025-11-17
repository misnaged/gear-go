package calls

import (
	"errors"
	"fmt"
	scalecodec "github.com/itering/scale.go"
	"github.com/itering/scale.go/types"
	gear_utils "github.com/misnaged/gear-go/internal/utils"
	"github.com/misnaged/gear-go/pkg/logger"
	"github.com/misnaged/substrate-api-rpc/keyring"
	rpcModels "github.com/misnaged/substrate-api-rpc/model"
	"github.com/misnaged/substrate-api-rpc/rpc"
)

func (calls *Calls) SignTransactionWithKeyring(moduleName, callName string, kr keyring.IKeyRing, params []scalecodec.ExtrinsicParam) (string, error) {
	if kr == nil {
		return "", fmt.Errorf("%w", errors.New("failed to sign transaction: keyring is nil"))
	}

	if err := calls.Meta.MetadataCheck(); err != nil {
		return "", fmt.Errorf("failed to sign transaction: %w", err)
	}
	genesisHash, err := calls.getChainGetBlockHash()
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}
	version, err := calls.getStateGetRuntimeVersion()
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	opt := &types.ScaleDecoderOption{Metadata: calls.Meta.GetMetadata(), Spec: -1}
	callIndex := gear_utils.GetCallLookupIndexByModuleAndCallNames(calls.Meta.GetMetadata(), moduleName, callName)

	resp, err := calls.GearRpc.SystemAccountNextIndex(kr.PublicKey())
	if err != nil {
		return "", fmt.Errorf("failed to send SystemAccountNextIndex request: %w", err)
	}
	logger.Log().Infof("nonce is: %d  %s %s \n", int(resp.Result.(float64)), callName, moduleName)
	calls.customTx = rpc.NewCustomTransaction(
		callIndex,
		genesisHash,
		"00", //todo: always immortal
		int(resp.Result.(float64)),
		version,
		calls.Meta.GetMetadata(),
		kr,
		opt,
		params,
	)
	signed, err := calls.customTx.SignTransactionCustom()
	if err != nil {
		return "", fmt.Errorf("failed to sign transaction: %w", err)
	}
	return signed, nil
}

func (calls *Calls) SignTransaction(moduleName, callName string, params []scalecodec.ExtrinsicParam) (string, error) {
	if calls.KeyRing == nil {
		return "", fmt.Errorf("%w", errors.New("failed to sign transaction: params is nil"))
	}
	if err := calls.Meta.MetadataCheck(); err != nil {
		return "", fmt.Errorf("failed to sign transaction: %w", err)
	}
	genesisHash, err := calls.getChainGetBlockHash()
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}
	version, err := calls.getStateGetRuntimeVersion()
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	opt := &types.ScaleDecoderOption{Metadata: calls.Meta.GetMetadata(), Spec: -1}
	callIndex := gear_utils.GetCallLookupIndexByModuleAndCallNames(calls.Meta.GetMetadata(), moduleName, callName)

	resp, err := calls.GearRpc.SystemAccountNextIndex(calls.KeyRing.PublicKey())
	if err != nil {
		return "", fmt.Errorf("failed to send SystemAccountNextIndex request: %w", err)
	}
	logger.Log().Infof("nonce is: %d  %s %s \n", int(resp.Result.(float64)), callName, moduleName)
	calls.customTx = rpc.NewCustomTransaction(
		callIndex,
		genesisHash,
		"00", //todo: always immortal
		int(resp.Result.(float64)),
		version,
		calls.Meta.GetMetadata(),
		calls.KeyRing,
		opt,
		params,
	)
	signed, err := calls.customTx.SignTransactionCustom()
	if err != nil {
		return "", fmt.Errorf("failed to sign transaction: %w", err)
	}
	return signed, nil
}

func (calls *Calls) getStateGetRuntimeVersion() (*rpcModels.RuntimeVersion, error) {
	runtimeVersion, err := calls.GearRpc.StateGetRuntimeVersionLatest()
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
func (calls *Calls) getChainGetBlockHash() (string, error) {

	genesisHash, err := calls.GearRpc.ChainGetBlockHash(0)
	if err != nil {
		return "", fmt.Errorf("request chain_getBlockHash failed: %v", err)
	}
	switch genesisHash.Result.(type) {
	case string:
		return genesisHash.Result.(string), nil
	default:
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
