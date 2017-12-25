package etc

import (
	"fmt"
	"net"
	"os"
	"strconv"
)

func (conf *Configuration) parseEnv() error {

	if pidFile := os.Getenv("PID_FILE"); pidFile != "" {
		conf.PidFile = pidFile
	}

	retryStartup := os.Getenv("RETRY_STARTUP")
	if retryStartup != "" {
		value, err := strconv.ParseBool(retryStartup)
		if err != nil {
			return fmt.Errorf("RETRY_STARTUP invalid, %s", err.Error())
		}
		conf.RetryStartup = value
	}

	if securityRoot := os.Getenv("SECURITY_ROOT"); securityRoot != "" {
		conf.SecurityRoot = securityRoot
	}

	//parse API listen env
	if err := parseListenEnv(conf); err != nil {
		return err
	}

	//parse logger env
	return parseLogger(conf)
}

func parseListenEnv(conf *Configuration) error {

	if apiHost := os.Getenv("API_LISTEN"); apiHost != "" {
		hostIP, hostPort, err := net.SplitHostPort(apiHost)
		if err != nil {
			return fmt.Errorf("API_LISTEN invalid, %s", err.Error())
		}
		if hostIP != "" {
			if _, err := net.LookupHost(hostIP); err != nil {
				return fmt.Errorf("API_LISTEN invalid, %s", err.Error())
			}
		}
		conf.Listen.Hosts = []string{net.JoinHostPort(hostIP, hostPort)}
	}

	if enableCors := os.Getenv("API_ENABLECORS"); enableCors != "" {
		value, err := strconv.ParseBool(enableCors)
		if err != nil {
			return fmt.Errorf("API_ENABLECORS invalid, %s", err.Error())
		}
		conf.Listen.EnableCors = value
	}
	return nil
}

func parseLogger(conf *Configuration) error {

	if logFile := os.Getenv("LOG_FILE"); logFile != "" {
		conf.Logger.LogFile = logFile
	}

	if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
		conf.Logger.LogLevel = logLevel
	}

	if logSize := os.Getenv("LOG_SIZE"); logSize != "" {
		value, err := strconv.ParseInt(logSize, 10, 64)
		if err != nil {
			return fmt.Errorf("LOG_SIZE invalid, %s", err.Error())
		}
		conf.Logger.LogSize = value
	}
	return nil
}
