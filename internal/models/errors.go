package models

import (
	"encoding/hex"
	"fmt"
	"github.com/itering/scale.go/types"
	"strings"
)

func GetMessageByIndex(moduleIdx, errorIdx int, meta *types.MetadataStruct) string {
	for _, v := range meta.Metadata.Modules {
		if v.Index == moduleIdx {
			for _, vv := range v.Errors {
				if vv.Index == errorIdx {
					return vv.Name
				}
			}
		}
	}
	return ""
}

func ConvertFromHexToInt(h string) (*int, error) {
	h = strings.TrimPrefix(h, "0x")

	if len(h)%2 != 0 {
		h = "0" + h
	}

	decoded, err := hex.DecodeString(h)
	if err != nil {
		return nil, fmt.Errorf("error decoding hex string: %v", err)
	}
	firstByte := int(decoded[0])
	return &firstByte, nil
}
