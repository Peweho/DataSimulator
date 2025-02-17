package etc

import (
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"simulator/util"
)

type Config struct {
	Time       Ftime    `yaml:"time"`
	Data       []Data   `yaml:"data"`
	Mq         MqClient `yaml:"mq"`
	DataBaseId string   `yaml:"dataBaseId"`
	Env        string   `yaml:"env"`
}

type Ftime struct {
	Start    string `yaml:"start"`
	End      string `yaml:"end"`
	Interval int    `yaml:"interval"`
}

type Data struct {
	Title     string   `yaml:"title"`
	Min       float64  `yaml:"min"`
	Max       float64  `yaml:"max"`
	Model     string   `yaml:"model"`
	Params    []string `yaml:"params"`
	Id        string   `yaml:"id"`
	Frequency int64    `yaml:"frequency"`
}

type MqClient struct {
	Addr      string `yaml:"addr"`
	Topic     string `yaml:"topic"`
	Partition int    `yaml:"partition"`
	TimeOut   int    `yaml:"timeout"`
}

var (
	cfg *Config
)

func GetConfig(path string) *Config {
	if cfg != nil && path == "" {
		return cfg
	}
	cfg = new(Config)
	var configFilePath string
	if filepath.IsAbs(path) {
		configFilePath = path
	} else {
		executablePath, err := os.Executable()
		if err != nil {
			util.Log.Fatalf("获取绝对路径：%v", err)
		}
		executableDir := filepath.Dir(executablePath)
		configFilePath = filepath.Join(executableDir, path)
	}

	// 读取配置文件
	yamlFile, err := os.ReadFile(configFilePath)
	if err != nil {
		util.Log.Fatalf("无法读取配置文件：%v", err)
	}

	// 解析YAML文件
	if err = yaml.Unmarshal(yamlFile, cfg); err != nil {
		util.Log.Fatalf("解析YAML文件失败：%v", err)
	}
	return cfg
}

// 判断真实环境还是虚拟环境
func IsVirtual() bool {
	return cfg.Env == "virtual"
}
