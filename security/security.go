package security

import yaml "gopkg.in/yaml.v2"

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

//SecurityData is exported
type SecurityData struct {
	AppID     string `yaml:"appid" json:"appid"`
	AppSecret string `yaml:"appsecret" json:"appsecret"`
	OriginID  string `yaml:"originid" json:"originid"`
	Token     string `yaml:"token" json:"token"`
	AESKey    string `yaml:"aeskey" json:"aeskey"`
}

//New is exported
//wechat security file load.
func New(rootPath string, originID string) (*SecurityData, error) {

	if strings.TrimSpace(originID) == "" {
		return nil, fmt.Errorf("security originid invalid.")
	}

	securityRoot, err := filepath.Abs(rootPath)
	if err != nil {
		return nil, fmt.Errorf("security root path invalid, %s", err)
	}

	securityFile := filepath.Join(securityRoot, originID) + ".yaml"
	fd, err := os.OpenFile(securityFile, os.O_RDONLY, 0777)
	if err != nil {
		return nil, fmt.Errorf("security file %s open failure, %s", securityFile, err)
	}

	defer fd.Close()
	buf, err := ioutil.ReadAll(fd)
	if err != nil {
		return nil, fmt.Errorf("security file %s read error, %s", securityFile, err)
	}

	securityData := &SecurityData{OriginID: originID}
	if err := yaml.Unmarshal(buf, securityData); err != nil {
		return nil, fmt.Errorf("security %s invalid, %s", securityFile, err)
	}
	return securityData, nil
}
