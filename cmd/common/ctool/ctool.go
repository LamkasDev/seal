package ctool

import "bytes"

func ByteArray256ToString(array [256]uint8) string {
	return string(array[:bytes.IndexByte(array[:], 0)])
}
