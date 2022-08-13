package models

import "encoding/json"

type Jsonescp struct{}

func (jsn *Jsonescp) JsonEscape(i string) string {
	b, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}
	s := string(b)
	return s[1 : len(s)-1]
}
