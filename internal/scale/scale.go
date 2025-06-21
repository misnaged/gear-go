package gear_scale

import (
	"errors"
	"fmt"
	"github.com/misnaged/gear-go/config"
	gear_rpc "github.com/misnaged/gear-go/internal/rpc"
	"github.com/misnaged/substrate-api-rpc/rpc"
	"github.com/vedhavyas/go-subkey/v2"
	"github.com/vedhavyas/go-subkey/v2/ed25519"
	"github.com/vedhavyas/go-subkey/v2/sr25519"

	"github.com/itering/scale.go/types"
	"github.com/misnaged/substrate-api-rpc/keyring"
)

const (
	VaraPrefix      uint16 = 137
	SubstratePrefix uint16 = 42
)

type Scale struct {
	gearRpc  gear_rpc.IGearRPC
	metadata *types.MetadataStruct
	config   *config.Scheme
	customTx rpc.ICustomTranscation
}

func NewScale(gearRpc gear_rpc.IGearRPC, config *config.Scheme) *Scale {
	return &Scale{
		gearRpc: gearRpc,
		config:  config,
	}
}

/*
  accountID, _ := GetKeyByName(sr25519.Scheme{}, "//Alice")
  fmt.Println(accountID)

  //Output: "d43593c715fdd31c61141abd04a99fd6822c8558854ccde39a5684e7a56da27d"
*/

func GetKeyByName(scheme subkey.Scheme, keyname string) (string, error) {
	kp, err := subkey.DeriveKeyPair(scheme, keyname)
	if err != nil {
		return "", fmt.Errorf("could not derive key pair: %w", err)
	}

	/*addr, _, err := subkey.SS58Decode(subkey.SS58Encode([]byte(accountId), format))
	if err != nil {
		return fmt.Errorf("could not decode account ID: %w", err)
	}
	fmt.Printf("Alice's AccountID: %x Public: %x Address %s \n ", kp.AccountID(), kp.Public(), kp.SS58Address(addr))
	*/
	return fmt.Sprintf("%x", kp.AccountID()), nil

}

func KeyRingByName(scheme subkey.Scheme, keyname string) (keyring.IKeyRing, error) {
	keyByName, err := GetKeyByName(scheme, keyname)
	if err != nil {
		return nil, fmt.Errorf("could not get key by name: %w", err)
	}

	switch scheme {
	case sr25519.Scheme{}:
		return keyring.New(keyring.Sr25519Type, keyByName), nil
	case ed25519.Scheme{}:
		return keyring.New(keyring.Ed25519Type, keyByName), nil
	default:
		return nil, errors.New("key scheme not supported")
	}
}
func KeyRingDefault(key string, category keyring.Category) keyring.IKeyRing {
	return keyring.New(category, key)
}
