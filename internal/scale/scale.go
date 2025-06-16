package gear_scale

import (
	"errors"
	"fmt"
	"github.com/misnaged/gear-go/config"
	gear_client "github.com/misnaged/gear-go/internal/client"

	scalecodec "github.com/itering/scale.go"
	"github.com/itering/scale.go/types"
	"github.com/itering/scale.go/utiles"
	"github.com/itering/substrate-api-rpc/hasher"

	"github.com/itering/substrate-api-rpc/keyring"
	"github.com/itering/substrate-api-rpc/util"
	"github.com/itering/substrate-api-rpc/util/ss58"
)

type Scale struct {
	metadata  *types.MetadataStruct
	extrinsic *scalecodec.GenericExtrinsic
	client    gear_client.IClient
	config    *config.Scheme
}

const TxVersionInfo = "84"

func NewScale(client gear_client.IClient, config *config.Scheme) *Scale {
	return &Scale{
		client: client,
		config: config,
	}
}

func (s *Scale) SignTransaction(callIndex, key string, args ...interface{}) (string, error) {
	decoder := &scalecodec.MetadataDecoder{}
	var postReq any
	resp, err := s.client.PostRequest(postReq, "state_getMetadata")
	if err != nil {
		return "", fmt.Errorf("post request err: %v", err)
	}

	decoder.Init(utiles.HexToBytes(resp.Result.(string)))
	err = decoder.Process()
	if err != nil {
		return "", fmt.Errorf("failed to decode metadata: %w", err)
	}

	kr := keyring.New(keyring.Sr25519Type, key)

	metadataStruct := decoder.Metadata
	opt := &types.ScaleDecoderOption{Metadata: &metadataStruct}

	var params []scalecodec.ExtrinsicParam
	for _, v := range args {
		params = append(params, scalecodec.ExtrinsicParam{Value: v})
	}

	encodeCall := types.EncodeWithOpt("Call", map[string]interface{}{"call_index": callIndex, "params": params}, opt)

	accountId, err := s.getAccountId(key)
	if err != nil {
		return "", fmt.Errorf("failed to get account id: %w", err)
	}

	genericExtrinsic := &scalecodec.GenericExtrinsic{
		VersionInfo: TxVersionInfo,
		Signer:      map[string]interface{}{"Id": kr.PublicKey()},
		Era:         "00",
		Nonce:       int(*accountId),
		Params:      params,
		CallCode:    callIndex,
	}
	genericExtrinsic.SignedExtensions = make(map[string]interface{})
	if util.StringInSlice("ChargeAssetTxPayment", metadataStruct.Extrinsic.SignedIdentifier) {
		genericExtrinsic.SignedExtensions["ChargeAssetTxPayment"] = map[string]interface{}{"tip": 0, "asset_id": nil}
	}
	if util.StringInSlice("CheckMetadataHash", metadataStruct.Extrinsic.SignedIdentifier) {
		genericExtrinsic.SignedExtensions["CheckMetadataHash"] = "Disabled"
	}

	payload, err := s.buildExtrinsicPayload(encodeCall, genericExtrinsic, &decoder.Metadata)
	if err != nil {
		return "", err
	}

	// if payload length > 256, Blake256 hash payload
	if len(util.HexToBytes(payload)) > 256 {
		payload = util.BytesToHex(hasher.HashByCryptoName(util.HexToBytes(payload), "Blake2_256"))
	}
	genericExtrinsic.SignatureRaw = map[string]interface{}{string(kr.Type()): utiles.AddHex(kr.Sign(util.AddHex(payload)))}

	encodedExtrinsic, err := genericExtrinsic.Encode(opt)
	if err != nil {
		return "", fmt.Errorf("failed to encode extrinsic: %w", err)
	}
	return util.AddHex(encodedExtrinsic), nil
}

func (s *Scale) buildExtrinsicPayload(encodeCall string, genericExtrinsic *scalecodec.GenericExtrinsic, meta *types.MetadataStruct) (string, error) {
	genesisHash, err := s.getChainGetBlockHash()
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}
	version, err := s.getStateGetRuntimeVersion()
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	//
	data := encodeCall
	data = data + types.Encode("EraExtrinsic", genericExtrinsic.Era)   // era
	data = data + types.Encode("Compact<U32>", genericExtrinsic.Nonce) // nonce
	if len(meta.Extrinsic.SignedIdentifier) > 0 && utiles.SliceIndex("ChargeTransactionPayment", meta.Extrinsic.SignedIdentifier) > -1 {
		data = data + types.Encode("Compact<Balance>", genericExtrinsic.Tip) // tip
	}

	for identifier, extension := range genericExtrinsic.SignedExtensions {
		for _, ext := range meta.Extrinsic.SignedExtensions {
			if ext.Identifier == identifier {
				data = data + types.Encode(ext.TypeString, extension)
			}
		}
	}
	data = data + types.Encode("U32", version.SpecVersion)        // specVersion
	data = data + types.Encode("U32", version.TransactionVersion) // transactionVersion
	data = data + util.TrimHex(types.Encode("Hash", genesisHash)) // genesisHash
	data = data + util.TrimHex(types.Encode("Hash", genesisHash)) // blockHash

	if _, ok := genericExtrinsic.SignedExtensions["CheckMetadataHash"]; ok {
		data = data + util.TrimHex("00") // CheckMetadataHash
	}
	return data, nil
}

func (s *Scale) getStateGetRuntimeVersion() (*RuntimeVersion, error) {
	genesisHash, err := s.client.PostRequest(nil, "state_getRuntimeVersion")
	if err != nil {
		return nil, fmt.Errorf("request state_getRuntimeVersion failed: %v", err)
	}
	rtm := &RuntimeVersion{}

	switch genesisHash.Result.(type) {
	case map[string]any:
		for key, val := range genesisHash.Result.(map[string]any) {
			switch key {
			//todo panic check! do not forget!
			case "apis":
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
func (s *Scale) getChainGetBlockHash() (string, error) {
	var params []int
	params = append(params, 0)
	genesisHash, err := s.client.PostRequest(params, "chain_getBlockHash")
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

func (s *Scale) getAccountId(accountKey string) (*uint64, error) {
	var params []string
	params = append(params, ss58.Encode(accountKey, 42))
	account, err := s.client.PostRequest(params, "system_accountNextIndex")
	if err != nil {
		return nil, fmt.Errorf("request system_accountNextIndex failed: %v", err)
	}
	switch account.Result.(type) {
	case float64:
		fmt.Println("account id:", account.Result.(float64))
		resultFromFloat := account.Result.(float64)
		rusultConverted := uint64(resultFromFloat)
		return &rusultConverted, nil
	case uint64:
		fmt.Println("account id:", account.Result.(*uint64))
		result := account.Result.(uint64)
		return &result, nil
	default:
		return nil, errors.New("account id is not uint64")
	}
}

type RuntimeVersion struct {
	Apis               []any  `json:"apis"`
	AuthoringVersion   int    `json:"authoringVersion"`
	ImplName           string `json:"implName"`
	ImplVersion        int    `json:"implVersion"`
	SpecName           string `json:"specName"`
	SpecVersion        int    `json:"specVersion"`
	TransactionVersion int    `json:"transactionVersion"`
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
