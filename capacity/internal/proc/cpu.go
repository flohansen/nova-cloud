package proc

import "os"

func GetCores() int32 {
	cores, _ := os.ReadDir("/sys/devices/virtual/cpuid")
	return int32(len(cores))
}
