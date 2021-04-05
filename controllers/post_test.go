package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type test struct { // 定义test结构体
	body string
	want ResCode
}

func TestCreatePostHandler(t *testing.T) {
	r := gin.Default()
	url := "/api/v1/post"
	r.POST(url, CreatePostHandler)

	// 方法一: 判断响应内容中是不是包含指定的字符串
	// assert.Contains(t, w.Body.String(), "用户未登录")

	bodyValid := `{
		"community_id":1,
		"title":"test",
		"content":"just a test"
	}`

	bodyInvalid := `{
		"title":"test",
		"content":"just a test"
	}`

	tests := map[string]test{ // 测试用例使用map存储
		"ValidParam":   {body: bodyValid, want: CodeNeedLogin},
		"InvalidParam": {body: bodyInvalid, want: CodeInvalidParam},
	}
	res := new(ResponseData)
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) { // 使用t.Run()执行子测试
			req, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader([]byte(tc.body)))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			// 判断状态码是否正确
			assert.Equal(t, 200, w.Code)
			json.Unmarshal(w.Body.Bytes(), res)
			// 判断响应码是否正确
			if !reflect.DeepEqual(res.Code, tc.want) {
				t.Errorf("excepted:%#v, got:%#v", tc.want, res.Code)
			}
		})
	}

	// 方法二: 将响应的内容反序列化到ResponseData 然后判断字段与预期是否一致
	// res := new(ResponseData)
	// if err := json.Unmarshal(w.Body.Bytes(), res); err != nil {
	// 	t.Fatalf("json.Unmarshal w.Body failed, err:%v\n", err)
	// }
	// assert.Equal(t, res.Code, CodeNeedLogin)
}
