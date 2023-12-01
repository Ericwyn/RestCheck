package commend

import (
	"fmt"
	"github.com/Ericwyn/RestCheck/conf"
	"github.com/Ericwyn/RestCheck/utils/log"
)

// InitProject 初始化一个 rest check commend
func InitProject(projectName string) {
	config := conf.LoadProjectConfig()
	if config != nil {
		log.E("当前目录下已有初始化项目配置")
		return
	}

	if projectName == "" || projectName == "null" {
		fmt.Println("请输入正确的项目名称")
		return
	}

	newConf := conf.BuildDefaultProjectConfig()

	newConf.ProjectName = projectName
	//projectName.
	conf.SaveProjectConfig(newConf)
}

func ShowProjectMsg() {
	projectConfig := conf.LoadProjectConfig()
	if projectConfig == nil {
		log.E("载入项目配置失败, 无法读取当前目录下的 restcheck.json 配置")
		return
	}

	fmt.Println("ProjectName: " + projectConfig.ProjectName)
	fmt.Println("BaseUrl: " + projectConfig.BaseUrl)
	fmt.Println("DefaultEnv: " + projectConfig.DefaultEnv)
	fmt.Println("ProjectConfigPath: " + conf.GetProjectConfigPath())
}
