package backends

import (
	"os"

	"github.com/craigjames16/cfstate/backends/local"
	"github.com/craigjames16/cfstate/backends/s3"
)

type Backend interface {
	GetState() (state []byte, err error)
	UpdateState(state []byte) (err error)
}

func GetBackend() (backend Backend) {
	var (
		s3Backend s3.S3Backend
		local     local.LocalBackend
	)
	switch os.Getenv("CFSTATE_BACKEND") {
	case "local":
		backend = local.NewBackend()
	case "s3":
		backend = s3Backend.NewBackend()
	default:
		panic("CFSTATE_BACKEND not set")
	}
	return backend
}
