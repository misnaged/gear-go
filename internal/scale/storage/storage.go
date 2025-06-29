package gear_storage

import gear_rpc "github.com/misnaged/gear-go/internal/rpc"

type IGearStorage interface {
	DecodeStorage(gearRPC gear_rpc.IGearRPC, decodeData any, isBigData bool) error
	BuildParams(accountId string) error
}
