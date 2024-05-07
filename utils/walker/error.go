package walker

import (
	"fmt"

	"github.com/layer5io/meshkit/errors"
)

var (
	ErrInvalidSizeFileCode       = "meshkit-11241"
	ErrCloningRepoCode           = "meshkit-11242"
	ErrInvokeFileInterceptorCode = ""
)

func ErrCloningRepo(err error) error {
	return errors.New(ErrCloningRepoCode, errors.Alert, []string{"could not clone the repo"}, []string{err.Error()}, []string{}, []string{})
}

func ErrInvalidSizeFile(err error) error {
	return errors.New(ErrInvalidSizeFileCode, errors.Alert, []string{err.Error()}, []string{"Could not read the file while walking the repo"}, []string{"Given file size is either 0 or exceeds the limit of 50 MB"}, []string{""})
}

func ErrInvokeFileInterceptor(err error, fileName string) error {
	return errors.New(ErrInvokeFileInterceptorCode, errors.Alert, []string{fmt.Sprintf("error invoke file interceptor for %s ", fileName)}, []string{err.Error()}, []string{}, []string{})
}
