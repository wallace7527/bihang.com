//这个例子用来展示urlencode编解码

package main

import (
	"fmt"
	"net/url"
)

func main() {
	msg, err := url.QueryUnescape("%e6%82%a8%e7%9a%84%e9%aa%8c%e8%af%81%e7%a0%81%e6%98%af%ef%bc%9a123%e3%80%90%e7%be%8e%e8%81%94%e3%80%91")
	if err != nil {
		fmt.Println("解析出错")
	}else {
		fmt.Println(msg)
	}

	fmt.Println(url.PathEscape(msg))
	fmt.Println(url.QueryEscape(" ff  tt"))
	fmt.Println(url.PathEscape(" ff  tt"))

	fmt.Println(url.PathEscape("您的验证码是：wang【美联】"))
	fmt.Println(url.QueryUnescape("%E6%82%A8%E7%9A%84%E9%AA%8C%E8%AF%81%E7%A0%81%E6%98%AF%EF%BC%9Awang%E3%80%90%E7%BE%8E%E8%81%94%E3%80%91"))
}

