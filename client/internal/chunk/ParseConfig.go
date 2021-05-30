package chunk

import (
	"DDSS/tools"
	"errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	Chunk struct {
		AverageSize int    `yaml:"averageSize"`
		MinSize int `yaml:"minSize"`
		MaxSize int `yaml:"maxSize"`
		Normalization int `yaml:"normalization"`
		DisableNormalization bool `yaml:"disableNormalization"`

	}
}

var ChunkConfig = new(Config)

func init() {

	yamlFile, err := ioutil.ReadFile("/Users/xylophone/code/go/DDSS/client/configs/config.yaml")
	if err != nil{
		tools.PrintError(err)

	}

	if yaml.Unmarshal(yamlFile, &ChunkConfig) != nil {
		yamler := errors.New("解析yaml出错！")
		tools.PrintError(yamler)
	}

}
