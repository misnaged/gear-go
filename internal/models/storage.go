package models

type FrameSystemAccountInfo struct {
	Nonce       any   `json:"nonce"`
	Consumers   any   `json:"consumers"`
	Provider    any   `json:"providers"`
	Sufficients any   `json:"sufficients"`
	Data        *Data `json:"data"`
}
type Data struct {
	Free     any `json:"free"`
	Reserved any `json:"reserved"`
	Frozen   any `json:"frozen"`
	Flags    any `json:"flags"`
}

type Program struct {
	Status     string
	Terminated string
	Active     *ActiveProgram `json:"active,omitempty"`
	ProgramId  string
}
type ActiveProgram struct {
	CodeId      any `json:"code_id"`
	State       any `json:"state"`
	Gas         any `json:"gas_reservation_map"`
	Mem         any `json:"memory_infix"`
	Block       any `json:"expiration_block"`
	Allocations any `json:"allocations_tree_len"`
}

func NewProgram(status, programId, terminated string, program *ActiveProgram) *Program {
	return &Program{
		Status:     status,
		Terminated: terminated,
		Active:     program,
		ProgramId:  programId,
	}
}
