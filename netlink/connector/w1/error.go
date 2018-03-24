package w1

import (
	"syscall"
)

type ErrorStatus uint8

func (err ErrorStatus) Errno() syscall.Errno {
	return syscall.Errno(err)
}

func (err ErrorStatus) Error() string {
	return err.Errno().Error()
}
