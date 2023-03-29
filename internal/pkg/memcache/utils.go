package memcache

import (
	"fmt"
	"strconv"
	"strings"
)

func isKeyValid(key string) error {
	if len(key) == 0 {
		return fmt.Errorf("empty")
	}
	if len(key) > 250 {
		return fmt.Errorf("too many symbols")
	}
	if strings.Contains(key, " ") ||
		strings.Contains(key, "\r") ||
		strings.Contains(key, "\n") {
		return fmt.Errorf("forbidden symbols")
	}
	return nil
}

func extractValueSize(line string) (int, error) {
	parts := strings.Split(strings.TrimSpace(line), " ")
	if len(parts) != 4 {
		return 0, fmt.Errorf("header of the get command is not correct")
	}
	return strconv.Atoi(parts[3])
}

func isEnd(line string) bool {
	return line == "END\r\n"
}

func isAValue(line string) bool {
	return strings.HasPrefix(line, "VALUE")
}

func isStored(line string) bool {
	return line == "STORED\r\n"
}

func isDeleted(line string) bool {
	return line == "DELETED\r\n"
}

func isNotFount(line string) bool {
	return line == "NOT_FOUND\r\n"
}
