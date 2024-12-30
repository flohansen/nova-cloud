package proc

import (
	"os"
	"strconv"
	"strings"
)

type MemInfo struct {
	MemTotal int32
	MemFree  int32
}

func GetMemInfo() MemInfo {
	b, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return MemInfo{}
	}

	info := MemInfo{}

	lines := strings.Split(string(b), "\n")
	for _, line := range lines {
		parts := strings.Split(line, ":")

		if len(parts) != 2 {
			continue
		}

		name := parts[0]
		valueString := strings.Replace(strings.TrimSpace(parts[1]), " kB", "", 1)

		value, err := strconv.ParseInt(valueString, 10, 32)
		if err != nil {
			continue
		}

		switch name {
		case "MemTotal":
			info.MemTotal = int32(value)
		case "MemFree":
			info.MemFree = int32(value)
		}
	}

	return info
}
