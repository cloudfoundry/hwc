package hwcconfig

import (
	"os"
	"path/filepath"
)

type HwcConfig struct {
	Instance      string
	Port          int
	RootPath      string
	TempDirectory string

	AspnetConfigPath          string
	WebConfigPath             string
	ApplicationHostConfigPath string
}

func New(port int, rootPath, tmpPath, uuid string) (error, *HwcConfig) {
	config := &HwcConfig{
		Instance:      uuid,
		Port:          port,
		RootPath:      rootPath,
		TempDirectory: tmpPath,
	}

	dest := filepath.Join(config.TempDirectory, "config")
	err := os.MkdirAll(dest, 0700)
	if err != nil {
		return err, nil
	}

	config.ApplicationHostConfigPath = filepath.Join(dest, "ApplicationHost.config")
	config.AspnetConfigPath = filepath.Join(dest, "Aspnet.config")
	config.WebConfigPath = filepath.Join(dest, "Web.config")

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
