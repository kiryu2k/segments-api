package parser

import (
	"strconv"
	"strings"
)

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
