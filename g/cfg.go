package g

import (
	"encoding/json"
	"github.com/toolkits/file"
	"log"
	"sync"
)

type HttpConfig struct {
	Enabled bool   `json:"enabled"`
	Listen  string `json:"listen"`
}

type GlobalConfig struct {
	Debug     bool        `json:"debug"`
	Hosts     string      `json:"hosts"`
	Database  string      `json:"database"`
	MaxIdle   int         `json:"maxIdle"`
	Listen    string      `json:"listen"`
    ExternalNodes   string  `json:"external_nodes"`
	Trustable []string    `json:"trustable"`
	Http      *HttpConfig `json:"http"`
}

var (
	ConfigFile string
	config     *GlobalConfig
	configLock = new(sync.RWMutex)
)

func Config() *GlobalConfig {
	configLock.RLock()
	defer configLock.RUnlock()
	return config
}

func ParseConfig(cfg string) {
	if cfg == "" {
		log.Fatalln("use -c to specify configuration file")
	}

	if !file.IsExist(cfg) {
		log.Fatalln("config file:", cfg, "is not existent")
	}

	ConfigFile = cfg

	configContent, err := file.ToTrimString(cfg)
	if err != nil {
		log.Fatalln("read config file:", cfg, "fail:", err)
	}

	var c GlobalConfig
	err = json.Unmarshal([]byte(configContent), &c)
	if err != nil {
		log.Fatalln("parse config file:", cfg, "fail:", err)
	}
    if !file.IsExist( c.ExternalNodes ) {
        log.Printf("WARN: the external_nodes file [%s] is not exist!", c.ExternalNodes)
        c.ExternalNodes = ""
    }

	configLock.Lock()
	defer configLock.Unlock()

	config = &c

	log.Println("read config file:", cfg, "successfully")
}
