package utils

import "time"

func ConvertStrDt(dt string) time.Time {
	t, err := time.Parse(time.RFC3339, dt)
	if err != nil {
		return time.Time{}
	}
	return t
}

func ConvertDtStr(dt time.Time) string {
	return dt.Format(time.RFC3339)
}
