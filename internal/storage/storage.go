package gear_storage

type IGearStorage interface {
	DecodeStorage(decodeData any, storkey string) error
	GetStorageKeys() ([]string, error)
	BuildParams(accountId string) error
	GetStorageKey() (string, error)
	DecodeStorageDataArray() ([]map[string]any, error)
	DecodeStorageDataMap(storkey string) (map[string]any, error)
	DecodeStorageDataAny(storkey string, v any) error
	GetProgramsId() ([]string, error)
}
