package hwcconfig

import (
	"fmt"
	"os"
	"path/filepath"
)

type HwcConfig struct {
	Instance                      string
	Port                          int
	TempDirectory                 string
	IISCompressedFilesDirectory   string
	ASPCompiledTemplatesDirectory string

	Applications              []*HwcApplication
	AspnetConfigPath          string
	WebConfigPath             string
	ApplicationHostConfigPath string
}

func New(port int, rootPath, tmpPath, contextPath, uuid string) (error, *HwcConfig) {
	config := &HwcConfig{
		Instance:                      uuid,
		Port:                          port,
		TempDirectory:                 tmpPath,
		IISCompressedFilesDirectory:   filepath.Join(tmpPath, "IIS Temporary Compressed Files"),
		ASPCompiledTemplatesDirectory: filepath.Join(tmpPath, "ASP Compiled Templates"),
	}

	defaultRootPath := filepath.Join(config.TempDirectory, "wwwroot")
	err := os.MkdirAll(defaultRootPath, 0700)
	if err != nil {
		return err, nil
	}

	configPath := filepath.Join(config.TempDirectory, "config")
	err = os.MkdirAll(configPath, 0700)
	if err != nil {
		return err, nil
	}

	err = os.MkdirAll(config.IISCompressedFilesDirectory, 0700)
	if err != nil {
		return err, nil
	}

	appPoolPath := fmt.Sprintf("AppPool%d", port)
	cachePath := filepath.Join(tmpPath, "IIS Temporary Compressed Files", appPoolPath)

	err = os.MkdirAll(cachePath, 0700)
	if err != nil {
		return err, nil
	}

	err = os.MkdirAll(config.ASPCompiledTemplatesDirectory, 0700)
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
