package api

import (
	"bufio"
	"bytes"
	"encoding/json"
	"github.com/Ericwyn/GoTools/file"
	"github.com/Ericwyn/RestCheck/conf"
	"github.com/Ericwyn/RestCheck/utils/rcutils"
	"io"
	"net/http"
	"os"
	"strings"
)

type ApiStorage struct {
	Name            string // 文件夹名称, 可能带有路径, 比如 note/listNote
	DirPath         string
	RequestFilePath string // .http 配置文件的路径
}

// BuildRequest
// 基于默认环境构造 http 请求
func (apiStorage *ApiStorage) BuildRequest() (*ApiRequest, error) {
	config := conf.LoadProjectConfig()

	if config == nil {
		return nil, nil
	}

	return apiStorage.BuildEnvRequest(config.DefaultEnv)
}

// BuildEnvRequest
// 基于指定环境构造 http 请求
func (apiStorage *ApiStorage) BuildEnvRequest(env string) (*ApiRequest, error) {
	request, err := apiStorage.buildBaseRequest()
	if err != nil {
		return nil, err
	}

	// 载入配置
	config := conf.LoadProjectConfig()
	apiEnv := config.Envs[env]

	request.RequestUrl = rcutils.BuildHttpUrl(apiEnv.BaseUrl, request.RequestUrl)
	request.Header = rcutils.BuildHeader(apiEnv.Header, request.Header)
	request.Env = env

	return &request, nil
}

// buildBaseRequest
// 读取 apiStorage 里面的 RequestFilePath
// 识别 .http 文件
func (apiStorage *ApiStorage) buildBaseRequest() (ApiRequest, error) {
	file, err := os.Open(apiStorage.RequestFilePath)
	if err != nil {
		return ApiRequest{}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan() // read the first line
	firstLine := scanner.Text()
	parts := strings.Split(firstLine, " ")
	httpMethod := parts[0]
	requestUrl := parts[1]

	header := make(map[string]string)
	var body strings.Builder
	isBody := false
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			isBody = true
			continue
		}
		if isBody {
			body.WriteString(line)
			body.WriteString("\n")
		} else {
			parts := strings.SplitN(line, ": ", 2)
			if len(parts) == 2 {
				header[parts[0]] = parts[1]
			}
		}
	}

	// 载入 baseConfig 里面的信息
	return ApiRequest{
		Name:       apiStorage.Name,
		RequestUrl: requestUrl,
		HttpMethod: httpMethod,
		Header:     header,
		Body:       body.String(),
	}, nil

}

type ApiRequest struct {
	Name       string
	Env        string
	RequestUrl string
	HttpMethod string
	Header     map[string]string
	Body       string
}

// DoRequest
// 实现 http 请求, 并且返回 response 文字信息
func (apiRequest *ApiRequest) DoRequest() (*ApiTestResponse, error) {
	req, err := http.NewRequest(apiRequest.HttpMethod, apiRequest.RequestUrl, bytes.NewBuffer([]byte(apiRequest.Body)))
	if err != nil {
		return nil, err
	}

	for key, value := range apiRequest.Header {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	headers := make(map[string]string)
	for key, values := range resp.Header {
		headers[key] = values[0]
	}

	return &ApiTestResponse{
		Request:    *apiRequest,
		Body:       string(body),
		Header:     headers,
		StatusCode: resp.StatusCode,
	}, nil
}

type ApiTestResponse struct {
	Request    ApiRequest
	Body       string            // response body
	Header     map[string]string // response 的 body
	StatusCode int               // response 的 statusCode
}

//// parseJsonStr
//// 把 response 的结果用 json 解析
//func (ApiTestResponse *ApiTestResponse) parseJsonStr() string {
//
//}

// ToIndentJsonStr
// 将 response 的结果转换成美化后的 json 字符串
func (apiResponse *ApiTestResponse) ToIndentJsonStr() string {
	var result map[string]interface{}
	err := json.Unmarshal([]byte(apiResponse.Body), &result)

	if err != nil {
		return "RestCheck Error: " + err.Error()
	}

	prettyJSON, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		return "RestCheck Error: " + err.Error()
	}

	return string(prettyJSON)

	//// 使用新的库来实现, 保证输出的 json 字段是有序的
	//var json = jsoniter.ConfigCompatibleWithStandardLibrary
	//data, err := json.MarshalIndent(result, "", "  ")
	//if err != nil {
	//	fmt.Println("error:", err)
	//}
	//return string(data)

}

// ListApiStorage 遍历当前目录下的 Api Storage 信息
func ListApiStorage(baseDirPath string) []*ApiStorage {
	//baseDirPath := conf.GetProjectConfigPath()
	resArr := make([]*ApiStorage, 0)

	//baseDirPath := paths.GetBasePath()
	baseDir := file.OpenFile(baseDirPath)
	for _, children := range baseDir.Children() {
		// 看看文件夹是否是
		if children.IsDir() {
			apiStorage := tryBuildApiStorage(children.AbsPath())
			if apiStorage != nil {
				resArr = append(resArr, apiStorage)
			}
		}
	}

	return resArr
}

// tryBuildApiStorage 把一个文件夹抽象成 ApiStorage 对象
func tryBuildApiStorage(dirPath string) *ApiStorage {
	openFile := file.OpenFile(dirPath)

	if !openFile.IsDir() {
		return nil
	}

	requestFilePath := openFile.AbsPath() + "/request.http"
	requestFile := file.OpenFile(requestFilePath)
	if requestFile.Exits() && requestFile.IsFile() {
		return &ApiStorage{
			Name:            openFile.Name(),
			DirPath:         openFile.AbsPath(),
			RequestFilePath: requestFilePath,
		}
	}

	return nil
}
