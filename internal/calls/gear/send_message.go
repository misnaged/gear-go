package gear_calls

import (
	"fmt"
)

func (gc *GearCalls) SendMessage(dest, value, payload string, gaslimit int, keepAlive bool) (string, error) {
	args := []any{dest, payload, gaslimit, value, keepAlive}
	call, err := gc.c.CallBuilder("send_message", "Gear", args)
	if err != nil {
		return "", fmt.Errorf("error calling extrinsic params builder: %w", err)
	}
	return call, nil

}
