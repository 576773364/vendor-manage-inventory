package utils

import (
	"fmt"
)

// ConstructSchemeKey 通过零售商名称构造补货方案的 key
func ConstructSchemeKey(name string) string {
	return fmt.Sprintf("Scheme-%s", name)
}

