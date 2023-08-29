package parser

import (
	"fmt"
	"strconv"
	"strings"
)

type TTL struct {
	Years  int
	Months int
	Days   int
}

func ParseResponse(buf []byte) (lhs uint64, rhs string) {
	if len(buf) < 2 {
		return
	}
	buf = (buf[:len(buf)-1])[1:]
	splittedBuf := strings.Split(string(buf), ",")
	if len(splittedBuf) != 2 {
		return
	}
	lhs, err := strconv.ParseUint(splittedBuf[0], 10, 64)
	if err != nil {
		return
	}
	rhs = splittedBuf[1]
	return
}

func ParseTTL(ttl string) (result *TTL, err error) {
	result = new(TTL)
	start := 0
	for i := 0; i < len(ttl); i++ {
		if ttl[i] >= '0' && ttl[i] <= '9' {
			continue
		}
		num, err := strconv.Atoi(ttl[start:i])
		if err != nil {
			return nil, fmt.Errorf("unexpected ttl format")
		}
		if strings.EqualFold(string(ttl[i]), "y") {
			if num < 1 || num > 100 {
				return nil, fmt.Errorf("years count must be between 1 and 100")
			}
			result.Years = num
		} else if strings.EqualFold(string(ttl[i]), "m") {
			if num < 1 || num > 11 {
				return nil, fmt.Errorf("months count must be between 1 and 11")
			}
			result.Months = num
		} else {
			if num < 1 || num > 30 {
				return nil, fmt.Errorf("days count must be between 1 and 30")
			}
			result.Days = num
		}
		start = i + 1
	}
	return
}
