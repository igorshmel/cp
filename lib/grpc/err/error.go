package grpc_err

import (
	"errors"
	"fmt"
	"strings"

	"google.golang.org/grpc/status"
)

func Get(err error) error {

	if err == nil {
		return nil
	}

	if isCustom(err) {
		return err
	}

	stat := status.Convert(err)

	if stat == nil {
		return err
	}

	return errors.New(fmt.Sprintf("code: %s, msg: %s", stat.Code(), stat.Message()))
}

func isCustom(err error) bool {

	if strings.HasPrefix(err.Error(), "code: ") {
		return true
	}

	return false
}
