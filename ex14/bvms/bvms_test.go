package bvms

import (
	"testing"
	"fmt"
)

var verifySignDatas = []struct {
	sig string
	msg string
	addr string
	want bool
	net int
}{
	//error
	{"a","","",false, 0},
	{"test","","",false, 0},
	{"IM30GkRoKvK57U3xUsw3XOWnnLqrltf5Z3GuKAOJCu4oVJiQZ2RBfrdzpHZ7u9YkdDW2lWptLJS7ARp3nAY4PeQ=","test","1nLsDVQyMQosdnBgLStXH1si8BP5xzKvbG",false, 0},

	//succ
	{"IM30GkRoKvK57U3xUsw3XOWnnLqrltf5Z3GuKAOJCu4oVJiQZ2RBfrdzpHZ7u9YkdDW2lWptLJS7ARp3nAY4PeQ=","test","mnLsDVQyMQosdnBgLStXH1si8BP5xzKvbG",true, RegressionNetID},
	{"IMSmaKvkv+RX6RUM87v4Brw5KvLdU/GFEZm6sl4lFWBVE2G7EVylg7Hrv4azrXkfvReJZ6WRkYTOt7Sjm1dOcZc=","this is a test", "1L3MNP8e4NSt5RFWCE6UuyyPRYW9kaAZ8C", true, MainNetID},
}

var netnames = []string{
	"UnknownNetID",
	"MainNetID",
	"RegressionNetID",
	"TestNet3ID",
	"SimNetID",
}

func getNetName(id int) string {
	if id >= 0 && id < len(netnames) {
		return netnames[id]
	}
	return fmt.Sprintf("ErrorNetId[%d]", id)
}

func TestVerifyMessages(t *testing.T) {
	for _, vsd := range verifySignDatas {

		err, net := VerifyMessage(vsd.addr, vsd.sig, vsd.msg)
		var pass bool
		if err != nil {
			pass = false
		}else {
			pass = true
		}

		if pass != vsd.want {
			t.Errorf("result want:%v, %v.", pass, err)
		}else if pass == true && net != vsd.net {
			t.Errorf("net want:%s, got:%s.", getNetName(vsd.net), getNetName(net) )
		}
	}
}

var addressDatas = [] struct {
	address string
	want bool
} {
	//error
	{"", false},
	{"123", false},
	{"bc1qwqdg6squsna38e46795at95yu9atm8azzmyvckulcc7kytlcckxswvvzej",false},
	{"mnLsDVQyMQosdnBgLStXH1si8BP5xzKvbGf", false},
	{"364CuAH97GvzFFfnhrbLQArPDUyCKmokZz", false},

	{"1111111111111111111114oLvT2",true},
	{"1234567yy2wy6LuJbf9G2NfCtjRvqCHGFR",true},
	{"188888888ikzoy8jmR2byT4WoQsycLxeUH", true},
	{"mnLsDVQyMQosdnBgLStXH1si8BP5xzKvbG", true},
	{"1L3MNP8e4NSt5RFWCE6UuyyPRYW9kaAZ8C", true},
}

func TestIsValidAddress(t *testing.T) {
	for _, d := range addressDatas {
		got := IsValidAddress(d.address)
		if got != d.want {
			t.Errorf("%s want:%v, got:%v", d.address, d.want, got)
		}
	}
}