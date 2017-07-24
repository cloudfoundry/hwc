package hwcconfig

import (
	"os"
	"path/filepath"
)

type HwcConfig struct {
	Instance      string
	Port          int
	TempDirectory string

	Applications              []*HwcApplication
	AspnetConfigPath          string
	WebConfigPath             string
	ApplicationHostConfigPath string
}

func New(port int, rootPath, tmpPath, contextPath, uuid string) (error, *HwcConfig) {
	config := &HwcConfig{
		Instance:      uuid,
		Port:          port,
		TempDirectory: tmpPath,
	}

	defaultRootPath := filepath.Join(tmpPath, "wwwroot")
	err := os.MkdirAll(defaultRootPath, 0700)
	if err != nil {
		return err, nil
	}

	configPath := filepath.Join(config.TempDirectory, "config")
	err = os.MkdirAll(configPath, 0700)
	if err != nil {
		return err, nil
	}

	config.Applications = NewHwcApplications(defaultRootPath, rootPath, contextPath)
	config.ApplicationHostConfigPath = filepath.Join(configPath, "ApplicationHost.config")
	config.AspnetConfigPath = filepath.Join(configPath, "Aspnet.config")
	config.WebConfigPath = filepath.Join(configPath, "Web.config")

	err = config.generateApplicationHostConfig()
	if err != nil {
		return err, nil
	}

	err = config.generateAspNetConfig()
	if err != nil {
		return err, nil
	}

	err = config.generateWebConfig()
	if err != nil {
		return err, nil
	}

	return nil, config
}
