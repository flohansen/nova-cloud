package proc

import (
	"runtime"
)

func GetCores() int32 {
	return int32(runtime.NumCPU())
}
