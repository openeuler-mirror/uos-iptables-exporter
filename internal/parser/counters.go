package parser

import (
	"fmt"
	"strings"
)

func parseCounters(token string) (uint64, uint64, bool) {
	token = strings.Trim(token, "[]")
	if token == "" { // 新增空值检查
		return 0, 0, false
	}

	parts := strings.Split(token, ":") // 改为冒号分割
	if len(parts) != 2 {
		return 0, 0, false
	}

	// 处理可能存在的空格
	parts[0] = strings.TrimSpace(parts[0])
	parts[1] = strings.TrimSpace(parts[1])

	var packets, bytes uint64
	_, err1 := fmt.Sscanf(parts[0], "%d", &packets)
	_, err2 := fmt.Sscanf(parts[1], "%d", &bytes)
	if err1 != nil || err2 != nil {
		return 0, 0, false
	}
	return packets, bytes, true
}
