package gormx

import (
	"fmt"
	"time"
)

type LocalTime time.Time

func (t *LocalTime) MarshalJSON() ([]byte, error) {
	tTime := time.Time(*t)
	return []byte(fmt.Sprintf("\"%v\"", tTime.Format("2006-01-02 15:04:05"))), nil
}
func (t *LocalTime) Format(layout string) string {
	tTime := time.Time(*t)
	return tTime.Format(layout)
}
