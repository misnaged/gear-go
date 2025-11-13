package gear_storage_methods

import (
	"encoding/hex"
	"fmt"
	"github.com/misnaged/substrate-api-rpc/hasher"
)

func encodeAccountId(accoundSeed string) ([]byte, error) {
	buffer := make([]byte, 32)
	_, err := hex.Decode(buffer, []byte(accoundSeed))
	if err != nil {
		return nil, fmt.Errorf("error decoding account: %v", err)
	}
	b := hasher.HashByCryptoName(buffer, "Blake2_128")
	return append(append([]byte{}, b...), buffer...), nil
}
func (stor *Storage) AddAccountIdToStorageParams(accoundSeed string) error {
	b, err := encodeAccountId(accoundSeed)
	if err != nil {
		return fmt.Errorf("error encoding account id: %v", err)
	}
	stor.params = append(stor.params, b)
	return nil
}
