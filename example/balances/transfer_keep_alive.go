package main

import (
	gear_go "github.com/misnaged/gear-go"
	"github.com/misnaged/gear-go/pkg/logger"
	"os"
)

const BobAccountId = "8eaf04151687736326c9fea17e25fc5287613693c912909cb226aa4794f26a48"

func main() {

	gear, err := gear_go.NewGear()
	if err != nil {
		logger.Log().Errorf("NewGear failed: %v", err)
		os.Exit(1)
	}
	args := []any{map[string]interface{}{"Id": BobAccountId}, 12345}
	s, err := gear.GetCalls().CallBuilder("transfer_keep_alive", "Balances", args)
	if err != nil {
		logger.Log().Errorf(" gear.calls.Builder failed: %v", err)
		os.Exit(1)
	}

	_, err = gear.GetClient().PostRequest([]string{s}, "author_submitExtrinsic")
}
