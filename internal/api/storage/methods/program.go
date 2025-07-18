package gear_storage_methods

import (
	"fmt"
	gear_rpc "github.com/misnaged/gear-go/internal/rpc"
	"strings"
)

func (stor *Storage) GetProgramsId(gearRPC gear_rpc.IGearRPC) ([]string, error) {
	prefix, err := stor.GetStorageKey()
	if err != nil {

		return nil, fmt.Errorf("failed to get storage key: %w", err)
	}
	storkeys, err := stor.getPagedKeys(gearRPC)
	if err != nil {
		return nil, fmt.Errorf("failed to get paged keys: %w", err)
	}
	var programIds []string
	for i := range storkeys {
		programIds = append(programIds, fmt.Sprintf("0x%s", strings.TrimPrefix(storkeys[i], prefix)))
	}

	return programIds, nil
}
