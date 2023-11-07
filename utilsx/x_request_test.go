package utilsx

import (
	"testing"
	"time"
)

func TestGetRequest(t *testing.T) {
	// 第一种写法，可直接传入完整url
	req := NewHttpRequest("https://qqlykm.cn/api/free/history/get")
	req.SetTimeout(5 * time.Second)

	// 第二种写法，可传入请求地址后使用setUri方法设置请求路径
	// req := NewHttpRequest("https://qqlykm.cn")
	// req.SetUri("api/free/history/get")

	// 第三种种写法，可传入请求域名后使用setUri方法设置请求路径，默认使用https协议
	// req := NewHttpRequest("qqlykm.cn")
	// req.SetUri("api/free/history/get")

	req.SetTimeout(1 * time.Second).Do(HTTP_METHOD_GET)
}
