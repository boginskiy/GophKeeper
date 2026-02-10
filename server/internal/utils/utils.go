package utils

import (
	"encoding/json"
	"time"
)

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

func DefinErr(err ...error) error {
	for i := range err {
		if err[i] != nil {
			return err[i]
		}
	}
	return nil
}

func ConversDtToTableView(t time.Time) time.Time {
	return t.Add(3 * time.Hour).UTC()
}

func Deserialization(data []byte, obj any) error {
	return json.Unmarshal(data, obj)
}

func Serialization(obj any) ([]byte, error) {
	return json.MarshalIndent(obj, "", "  ")
}
