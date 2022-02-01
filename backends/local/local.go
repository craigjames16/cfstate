package local

import (
	"fmt"
	"os"
	"time"
)

var (
	BUCKET_NAME     string
	STATE_FILE_NAME string
)

type LocalBackend struct{}

func init() {
	var (
		envExists bool
	)
	STATE_FILE_NAME, envExists = os.LookupEnv("CFSTATE_STATE_FILE_NAME")

	if !envExists {
		panic(fmt.Errorf("CFSTATE_STATE_FILE_NAME not set"))
	}
}

func (lb LocalBackend) NewBackend() LocalBackend {
	return LocalBackend{}
}

func (lb LocalBackend) UpdateState(stateFile []byte) (err error) {
	var (
		now                time.Time = time.Now()
		sec                int64     = now.Unix()
		stateFileName      string    = fmt.Sprintf("./%s.json", STATE_FILE_NAME)
		savedStateFileName string    = fmt.Sprintf("./prev_states/state-%d.json", sec)
	)

	if err = os.Rename(stateFileName, savedStateFileName); err != nil {
		return err
	}

	if err = os.WriteFile(stateFileName, stateFile, 0644); err != nil {
		return err
	}

	return err

}

func (lb LocalBackend) GetState() (output []byte, err error) {
	stateFileName := fmt.Sprintf("./%s.json", STATE_FILE_NAME)
	output, err = os.ReadFile(stateFileName)
	if err != nil {
		return nil, err
	}

	return output, err
}
