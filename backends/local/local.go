package local

var (
	BUCKET_NAME     string
	STATE_FILE_NAME string
)

type LocalBackend struct{}

func (lb LocalBackend) NewBackend() LocalBackend {
	return LocalBackend{}
}

func (lb LocalBackend) UpdateState(stateFile []byte) (err error) {

}

func (lb LocalBackend) GetState() (output []byte, err error) {
}
