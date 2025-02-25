package reponse

import (
	"fmt"
	"time"
)

type JsonTime time.Time

func (j JsonTime) MarshalJSON() ([]byte, error) {
	var stmp = fmt.Sprintf("\"%s\"", time.Time(j).Format("2025-2-25"))
	return []byte(stmp),nil
}

type UserResponse struct {
	Id       int32  `json:"id"`
	NickName string `json:"name"`
	// BirthDay time.Time `json:"birthday"`法一
	// BirthDay string `json:"birthday"`
	BirthDay JsonTime `json:"birthday"`//法三加上上面的MarshalJSON
	Gender   string    `json:"gender"`
	Mobile   string    `json:"mobile"`
}
