package gear_calls

type IGearCalls interface {
	UploadCode(pathToWasm string) error
}
