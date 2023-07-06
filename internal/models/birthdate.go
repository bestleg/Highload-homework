package models

import (
	"encoding/json"
	"strings"
	"time"
)

type JsonBirthDate time.Time

func (j *JsonBirthDate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	*j = JsonBirthDate(t)
	return nil
}

func (j JsonBirthDate) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(j).Format("2006-01-02"))
}

func (j JsonBirthDate) Format(s string) string {
	t := time.Time(j)
	return t.Format(s)
}
