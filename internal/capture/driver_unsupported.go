//go:build !darwin || !cgo

package capture

import "errors"

func New() (Driver, error) {
	return nil, errors.New("keyboard capture requires macOS and a build with CGO_ENABLED=1")
}
