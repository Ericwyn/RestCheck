package paths

import (
	"github.com/Ericwyn/RestCheck/utils/log"
	"os"
	"path/filepath"
	"strings"
)

var runnerPath = ""

// GetExecutablePath
// 获取文件执行路径, goland 里面会获取到 /tmp/Goland
func GetExecutablePath() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	return exPath
}

// GetRunnerPath
// 获取项目所在路径, 就是我们正常理解里面的 ./
func GetRunnerPath() string {
	if runnerPath == "" {
		//返回绝对路径  filepath.Dir(os.Args[0])去除最后一个元素的路径

		//log.D("os.Args[0]:" + os.Args[0])

		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			log.E("无法获取程序运行目录")
			log.E(err)
		}

		//将\替换成/
		runnerPath = strings.Replace(dir, "\\", "/", -1)

		//log.D("程序运行目录:" + runnerPath)

		// 如果运行的目录是在 Temp 下面的话, 那么看看 ./ 目录是什么
		if strings.Contains(runnerPath, "AppData/Local/Temp") ||
			strings.HasPrefix(runnerPath, "/tmp") {
			//log.D("程序运行在 Temp 目录")
			dir, err := filepath.Abs(filepath.Dir("./"))
			if err != nil {
				log.E("获取 ./ 目录绝对路径失败")
				log.E(err)
			}

			runnerPath = strings.Replace(dir, "\\", "/", -1)
			//log.D("程序运行目录更新为: " + runnerPath)
		}
	}

	return runnerPath
}
