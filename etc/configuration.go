package etc

import "github.com/8micro/gutils/logger"
import "gopkg.in/yaml.v2"

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"sync"
)

var (
	SystemConfig *Configuration = nil
)

var (
	ErrConfigFileNotFound      = errors.New("config file not found.")
	ErrConfigGenerateFailure   = errors.New("config file generated failure.")
	ErrConfigFormatInvalid     = errors.New("config file format invalid.")
	ErrConfigServerDataInvalid = errors.New("config server data invalid.")
)

// Configuration is exported
type Configuration struct {
	sync.RWMutex
	Version      string `yaml:"version" json:"version"`
	PidFile      string `yaml:"pidfile" json:"pidfile"`
	RetryStartup bool   `yaml:"retrystartup" json:"retrystartup"`
	SecurityRoot string `yaml:"securityroot" json:"securityroot"`

	Listen struct {
		Hosts      []string `yaml:"hosts" json:"hosts"`
		EnableCors bool     `yaml:"enablecors" json:"enablecors"`
	} `yaml:"listen" json:"listen"`

	Logger struct {
		LogFile  string `yaml:"logfile" json:"logfile"`
		LogLevel string `yaml:"loglevel" json:"loglevel"`
		LogSize  int64  `yaml:"logsize" json:"logsize"`
	} `yaml:"logger" json:"logger"`
}

// New is exported
func New(file string) error {

	buf, err := readConfigurationFile(file)
	if err != nil {
		return fmt.Errorf("config error %s", err.Error())
	}

	conf := &Configuration{RetryStartup: true}
	if err := yaml.Unmarshal(buf, conf); err != nil {
		return ErrConfigFormatInvalid
	}

	if err = conf.parseEnv(); err != nil {
		return fmt.Errorf("config parse env %s", err.Error())
	}

	SystemConfig = parseDefaultParmeters(conf)
	log.Printf("[#etc#] version: %s\n", SystemConfig.Version)
	log.Printf("[#etc#] pidfile: %s\n", SystemConfig.PidFile)
	log.Printf("[#etc#] retrystartup: %s\n", strconv.FormatBool(SystemConfig.RetryStartup))
	log.Printf("[#etc#] securityroot: %s\n", SystemConfig.SecurityRoot)
	log.Printf("[#etc#] listen: %+v\n", SystemConfig.Listen)
	log.Printf("[#etc#] logger: %+v\n", SystemConfig.Logger)
	return nil
}

//PidFile is exported
func PidFile() string {

	if SystemConfig != nil {
		return SystemConfig.PidFile
	}
	return ""
}

//RetryStartup is exported
func RetryStartup() bool {

	if SystemConfig != nil {
		return SystemConfig.RetryStartup
	}
	return false
}

//SecurityRoot is exported
func SecurityRoot() string {

	if SystemConfig != nil {
		return SystemConfig.SecurityRoot
	}
	return ""
}

//LoggerArgs is exported
func LoggerArgs() *logger.Args {

	if SystemConfig != nil {
		return &logger.Args{
			FileName: SystemConfig.Logger.LogFile,
			Level:    SystemConfig.Logger.LogLevel,
			MaxSize:  SystemConfig.Logger.LogSize,
		}
	}
	return nil
}

func readConfigurationFile(file string) ([]byte, error) {

	fd, err := os.OpenFile(file, os.O_RDONLY, 0777)
	if err != nil {
		return nil, err
	}

	defer fd.Close()
	buf, err := ioutil.ReadAll(fd)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func parseDefaultParmeters(conf *Configuration) *Configuration {

	if conf.SecurityRoot == "" {
		conf.SecurityRoot = "./security"
	}

	if conf.Logger.LogLevel == "" {
		conf.Logger.LogLevel = "info"
	}

	if conf.Logger.LogSize == 0 {
		conf.Logger.LogSize = 20971520
	}
	return conf
}
