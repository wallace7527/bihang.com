package handles

import (
	"time"
	"fmt"
)



type JsonTime time.Time


//实现它的json序列化方法
func (this JsonTime) MarshalJSON() ([]byte, error) {
	var stamp = fmt.Sprintf("\"%s\"", time.Time(this).Format("2006-01-02 15:04:05"))
	return []byte(stamp), nil
}


type DataResult struct {
	Tm JsonTime			`json:"write_time"`
	State int32			`json:"state"`
	Message interface{}  `json:"msg"`
}

type VerifyCode struct {
	Vcode string			`json:"vcode"`
	//Timestramp JsonTime		`json:"timestramp"`
}

type LoginSucc struct {
	Uid int64 		`json:"id"`
	Token string 	`json:"token"`
}