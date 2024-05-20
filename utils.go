package xsql

import (
	"unsafe"
)

// bytesToString []byte转string
func bytesToString(b []byte) string {
	if len(b) == 0 {
		return ""
	}
	return unsafe.String(&b[0], len(b))
}

// stringToBytes string 转[]byte
func stringToBytes(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

func genEntity(length int) []any {
	entity := make([]any, 0, length)
	for i := 0; i < length; i++ {
		entity = append(entity, new([]byte))
	}

	return entity
}

func genRecord(data []any, fields []string) map[string]any {
	record := make(map[string]any, len(fields))
	for k, v := range data {
		record[fields[k]] = bytesToString(*v.(*[]byte))
	}

	return record
}
