package commend

import (
	"errors"
	"fmt"
	"github.com/Ericwyn/GoTools/date"
	"github.com/Ericwyn/GoTools/file"
	"github.com/Ericwyn/RestCheck/api"
	"github.com/Ericwyn/RestCheck/conf"
	"github.com/Ericwyn/RestCheck/utils/log"
	"github.com/Ericwyn/RestCheck/utils/paths"
	"os"
	"path"
	"sort"
	"strings"
	"time"
)

// InitApi 初始化 API 配置
func InitApi(apiName string) {

}

func CheckAllApi(env string, saveResult bool) {

}

func CheckApi(env string, apiName string, saveResult bool) {
	log.D("开始测试 API, apiName: ", apiName, ", env: ", env, ", saveResult: ", saveResult)
	config := conf.LoadProjectConfig()
	if config == nil {
		log.E("无法载入项目配置")
		return
	}

	apiStorages := api.ListApiStorage(paths.GetRunnerPath())

	var runApi *api.ApiStorage
	// 找到名称为 apiName 的 apiStorage
	for _, storage := range apiStorages {
		if storage.Name == apiName {
			runApi = storage
			break
		}
	}

	if runApi == nil {
		log.E("未找到名称为 " + apiName + " 的 api")
		return
	}

	var request *api.ApiRequest
	var err error
	if env == "" {
		request, err = runApi.BuildRequest()
	} else {
		request, err = runApi.BuildEnvRequest(env)
	}

	if err != nil {
		log.E("构建 ApiRequest 失败", err)
		return
	}

	if request == nil {
		log.E("构建 ApiRequest 失败")
		return
	}

	resp, err := request.DoRequest()
	if err != nil {
		log.E("请求失败", err)
		return
	}

	handlerResult(resp, saveResult)
	//fmt.Println(resp.ToIndentJsonStr())
}

// 处理 API 请求结果
func handlerResult(response *api.ApiTestResponse, saveResult bool) {
	if saveResult {
		saveTestResult(response)
	}

	checkTestResult(response)

	printTestResult(response)
}

// 保存 API 请求结果
func saveTestResult(response *api.ApiTestResponse) {
	// 需要保存的位置
	// 基础目录 + Api 名称 + 环境名称
	resultSaveDirPath := path.Join(conf.GetBasePath(), response.Request.Name, response.Request.Env)

	os.MkdirAll(resultSaveDirPath, 0755)

	// 20231129_141731_results.txt
	resultFileName := date.Format(time.Now(), "yyyyMMdd_HHmmss_result.txt")

	str := response.ToIndentJsonStr()

	finalSavePath := path.Join(resultSaveDirPath, resultFileName)
	log.D("开始保存 API 请求结果 ", response.Request.Name+", "+response.Request.Env+", "+finalSavePath)
	err := os.WriteFile(finalSavePath, []byte(str), 0755)
	if err != nil {
		log.E("保存 API 请求结果失败", err)
	}
}

// 检查 API 请求结果
func checkTestResult(response *api.ApiTestResponse) {
	lastResultFile, err := findFinalResult(response)
	if err != nil {
		log.E("未找到可校验匹配的测试结果, ", err)
		return
	}

	result, err := lastResultFile.Read()

	log.D("开始校验 API 请求结果")

	// 此次的测试结果
	resultNow := response.ToIndentJsonStr()

	var resultFileName string
	// 校验
	if resultNow == string(result) {
		log.D(response.Request.Name, "测试结果一致, 匹配: ", lastResultFile.Name())
		// 将测试结果写到 yyyyMMdd_HHmmss_check_success.txt
		resultFileName = date.Format(time.Now(), "yyyyMMdd_HHmmss") + "_check_success.txt"
	} else {
		log.D(response.Request.Name, "测试结果不一致!!! 匹配: ", lastResultFile.Name())
		// 将测试结果写到 yyyyMMdd_HHmmss_check_fail.txt
		resultFileName = date.Format(time.Now(), "yyyyMMdd_HHmmss") + "_check_fail.txt"
	}

	finalSavePath := path.Join(lastResultFile.ParentPath(), resultFileName)
	//log.D("开始保存 API 请求结果 ", response.Request.Name+", "+response.Request.Env+", "+finalSavePath)
	os.WriteFile(finalSavePath, []byte(resultNow), 0755)
	if err != nil {
		log.E("保存 API 请求结果失败", err)
	}
}

// 找到最后一次 result 来做比较
func findFinalResult(response *api.ApiTestResponse) (*file.File, error) {
	// 需要保存的位置
	// 基础目录 + Api 名称 + 环境名称
	resultSaveDirPath := path.Join(conf.GetBasePath(), response.Request.Name, response.Request.Env)

	os.MkdirAll(resultSaveDirPath, 0755)

	// 遍历 dir, 找到最后一个 yyyyMMdd_HHmmss_result.txt 文件
	dir := file.OpenFile(resultSaveDirPath)
	log.D("开始遍历 ", resultSaveDirPath)
	files := dir.Children()
	// files 按照 Name 倒序
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() > files[j].Name()
	})

	var finalResultFiles *file.File

	for _, f := range files {
		log.D(f.Name())
		if strings.HasSuffix(f.Name(), "_result.txt") {
			finalResultFiles = &f
			return finalResultFiles, nil
		}
	}

	if finalResultFiles == nil {
		return nil, errors.New("无法获取此前测试结果")
	}

	//bytes, err := finalResultFiles.Read()
	//if err != nil {
	//	log.E("读取此前测试结果失败 ", err)
	//	return nil, err
	//}

	return finalResultFiles, nil
}

func printTestResult(response *api.ApiTestResponse) {
	fmt.Println("============================================")
	fmt.Println("Api: " + response.Request.Name)
	fmt.Println("Env: " + response.Request.Env)
	fmt.Println("Path: " + response.Request.RequestUrl)
	fmt.Println("============================")
	fmt.Println(response.ToIndentJsonStr())
	fmt.Println("============================================")
}
