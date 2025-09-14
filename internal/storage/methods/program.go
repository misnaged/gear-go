package gear_storage_methods

import (
	"errors"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/misnaged/gear-go/internal/models"
	"strings"
)

func (stor *Storage) GetProgramsId() ([]string, error) {
	prefix, err := stor.GetStorageKey()
	if err != nil {

		return nil, fmt.Errorf("failed to get storage key: %w", err)
	}
	storkeys, err := stor.getPagedKeys()
	if err != nil {
		return nil, fmt.Errorf("failed to get paged keys: %w", err)
	}
	var programIds []string
	for i := range storkeys {
		programIds = append(programIds, programId(storkeys[i], prefix))
	}

	return programIds, nil
}
func programId(s string, prefix string) string {
	return fmt.Sprintf("0x%s", strings.TrimPrefix(s, prefix))
}

func (stor *Storage) GetPrograms() ([]*models.Program, error) {
	var programs []*models.Program
	prefix, err := stor.GetStorageKey()
	if err != nil {

		return nil, fmt.Errorf("failed to get storage key: %w", err)
	}
	keys, err := stor.getPagedKeys()
	if err != nil {
		return nil, fmt.Errorf("error getting paged keys: %v", err)
	}
	for i := range keys {
		m, err := stor.DecodeStorageDataMap(keys[i])
		if err != nil {
			return nil, fmt.Errorf("error decoding storage data map: %v", err)
		}
		for key, value := range m {
			var program models.ActiveProgram

			if key != "Terminated" {

				b, err := json.Marshal(value)
				if err != nil {
					return nil, fmt.Errorf("error marshalling data: %v", err)
				}
				err = json.Unmarshal(b, &program)
				if err != nil {
					return nil, fmt.Errorf("error unmarshalling data: %v", err)
				}
				programs = append(programs, models.NewProgram(key, programId(keys[i], prefix), "", &program))
			} else {
				if value.(string) != "" {
					programs = append(programs, models.NewProgram(key, programId(keys[i], prefix), value.(string), nil))
				}
			}
		}
	}
	return programs, nil
}

func (stor *Storage) GetProgramByCodeId(codeId string) (*models.Program, error) {
	programs, err := stor.GetPrograms()
	if err != nil {
		return nil, fmt.Errorf("error getting all programs: %w", err)
	}
	for _, program := range programs {
		if program.Active != nil {
			if program.Active.CodeId == codeId {
				return program, nil
			}
		}
	}
	return nil, errors.New("active program not found, possibly due to early termination")
}
