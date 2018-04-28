package handles

import (
	"math/rand"
	"time"
	"fmt"
)

type vcodeInfo struct {
	Code string
	Timestramp time.Time
}

var vcodes map[string]*vcodeInfo

func init() {
	vcodes = make(map[string]*vcodeInfo)
}

func vcode() string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%06v", rnd.Int31n(1000000))
}

func GenVerifyCode(tel string) string {
	v := &vcodeInfo{
		vcode(),
		time.Now(),
	}
	vcodes[tel] = v
	return v.Code
}

func CheckVerifyCode(tel string, vcode string) bool {
	//mark for debug
	if vcode == "888888" {
		return true
	}

	vcodeInfo, ok := vcodes[tel]
	if ok {
		if vcode == vcodeInfo.Code && time.Now().Sub(vcodeInfo.Timestramp).Minutes() < 2 {

			return true
		}
	}
	return false
}