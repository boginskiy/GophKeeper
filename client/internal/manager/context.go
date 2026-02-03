package manager

import "google.golang.org/grpc/metadata"

func TakeValueFromHeader(header metadata.MD, field string, idx int) string {
	values := header.Get(field)
	if len(values) > 0 {
		return values[idx]
	}
	return ""
}

func CreateHeaderWithValue(key, value string) metadata.MD {
	return metadata.Pairs(key, value)
}
