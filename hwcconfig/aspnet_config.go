package hwcconfig

import "os"

func (c *HwcConfig) generateAspNetConfig() error {
	file, err := os.Create(c.AspnetConfigPath)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(aspnetConfigTemplate)
	if err != nil {
		return err
	}
	return nil
}

const aspnetConfigTemplate = `<?xml version="1.0" encoding="UTF-8"?>
<configuration>
  <runtime>
    <legacyUnhandledExceptionPolicy enabled="false" />
    <legacyImpersonationPolicy enabled="true"/>
    <alwaysFlowImpersonationPolicy enabled="false"/>
    <SymbolReadingPolicy enabled="1" />
  </runtime>
  <startup useLegacyV2RuntimeActivationPolicy="true" />
</configuration>
`
