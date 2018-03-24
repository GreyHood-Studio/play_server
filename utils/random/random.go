package random

import "hash/fnv"

func hashingStrToInt(hid string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(hid))
	return h.Sum32()
}