package conf

import (
	"encoding/json"
	"github.com/Ericwyn/RestCheck/utils/log"
	"github.com/Ericwyn/RestCheck/utils/paths"
	"os"
)

var projectConfigCache *ProjectConfig

var runnerPath string

func SaveProjectConfig(config ProjectConfig) {
	if checkErr := config.check(); checkErr != "" {
		log.E("project config 保存失败, " + checkErr)
		return
	}

	jsonStr, _ := json.MarshalIndent(config, "  ", "  ")
	os.WriteFile(GetProjectConfigPath(), jsonStr, 0755)
}

// GetBasePath
// 整个 rest check 工程的根目录, 该目录下有 restcheck.json 配置
func GetBasePath() string {
	if runnerPath == "" {
		runnerPath = paths.GetRunnerPath()
	}
	return runnerPath
}

func SetBasePath(path string) {
	runnerPath = path
}

func GetProjectConfigPath() string {

	projectConfigJsonPath := GetBasePath() + "/restcheck.json"

	return projectConfigJsonPath
}

// LoadProjectConfig
// 读取当前目录的 commend 信息
func LoadProjectConfig() *ProjectConfig {
	if projectConfigCache != nil {
		return projectConfigCache
	}

	// 读取 project config
	//fileBytes, err := io.ReadFile(projectConfigJsonPath)
	bytes, err := os.ReadFile(GetProjectConfigPath())
	if err != nil {
		//log.E("载入项目配置失败, 无法读取当前目录下的 restcheck.json 配置")
		return nil
	}

	var projectConfig ProjectConfig
	err = json.Unmarshal(bytes, &projectConfig)
	if err != nil {
		//log.E("载入项目配置失败, 无法读取当前目录下的 restcheck.json 配置")
		return nil
	}

	projectConfigCache = &projectConfig
	return projectConfigCache
}
