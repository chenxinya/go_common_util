package yamlReader

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)
//FeConf 如果定义结构体字段名首字母是小写的，这意味着这些字段在包外不可见,因而无法在其他包中被访问，只允许包内访问。所以定义字段的时候需要首字母大写
// 以两种方法来定义 字段
//在定义结构体字段时，除字段名称和数据类型外，还可以使用反引号为结构体字段声明元信息，这种元信息称为Tag，用于编译阶段关联到字段当中
type FeConf struct {
	FeIp string `yaml:"fe_ip"`
	FePort string `yaml:"fe_port"`
	User string `yaml:"user"`
	Password string `yaml:"password"`
	OutPath string `yaml:"out_path"`
}

func (c *FeConf) GetFeConf(path string) *FeConf {

	yamlFile, err1 := ioutil.ReadFile(path)
	if err1 != nil {
		log.Printf("yamlFile.Get err   #%v ", err1)
	}

	err := yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Printf("Unmarshal: %v", err)
	}
	return c
}

