package internal

import (
	"errors"
	"strconv"
	"strings"
)

func ParseVersion(ver string) (i int, err error) {
	v := strings.TrimPrefix(ver, "v")
	i, err = strconv.Atoi(v)
	if i < 1 {
		return -1, errors.New("version number too low")
	}
	return
}
