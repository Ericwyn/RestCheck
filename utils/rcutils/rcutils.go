package rcutils

import (
	"strings"
)

// Rest Check 的工具方法

//--------------------------

// BuildHttpUrl
// 拼接成完整 URL，并且去除中间多余的 /
func BuildHttpUrl(baseUrl string, apiPath string) string {
	var protocol string
	if strings.HasPrefix(strings.ToLower(baseUrl), "https://") {
		protocol = "https://"
		baseUrl = strings.TrimPrefix(strings.ToLower(baseUrl), protocol)
	} else if strings.HasPrefix(strings.ToLower(baseUrl), "http://") {
		protocol = "http://"
		baseUrl = strings.TrimPrefix(strings.ToLower(baseUrl), protocol)
	}

	baseUrl = strings.TrimSuffix(baseUrl, "/")
	apiPath = strings.TrimPrefix(apiPath, "/")

	return protocol + baseUrl + "/" + apiPath
}

// BuildHeader
// 拼接多个来源的 Http Header 参数
// 所有 header 名称都要转小写
// 如果 special header 里面有同样的 header, 就会覆盖掉 base header
func BuildHeader(baseHeader map[string]string, specialHeader map[string]string) map[string]string {
	header := make(map[string]string)
	for k, v := range baseHeader {
		header[strings.ToLower(k)] = v
	}
	for k, v := range specialHeader {
		header[strings.ToLower(k)] = v
	}
	return header

}
