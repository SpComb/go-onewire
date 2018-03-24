package w1

import (
	"fmt"
)

type ErrorStatus uint8

func (err ErrorStatus) Error() string {
	return fmt.Sprintf("errno -%d", uint8(err))
}
