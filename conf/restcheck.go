package conf

type ProjectConfig struct {
	ProjectName string // 项目名称
	Header      map[string]string
	BaseUrl     string
	Envs        []EnvSetting
	DefaultEnv  string
}

type EnvSetting struct {
	Name    string // 环境名称
	Header  map[string]string
	BaseUrl string
}

// 检查 commend config 是否正确
// 如果没问题就返回 "", 否则返回异常信息
func (config *ProjectConfig) check() string {
	// 检查 Envs 是否为空, 需要有一个默认的 envs

	// 检查 DefaultEnv 是否在 envs 里面

	return ""
}
