package validator

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	_ "runtime/cgo"
	"strings"
)

type Configuration struct {
	XMLName         xml.Name `xml:"configuration"`
	SystemWebServer SystemWebServer
}

type SystemWebServer struct {
	XMLName         xml.Name `xml:"system.webServer"`
	HTTPCompression HTTPCompression
}

type HTTPCompression struct {
	XMLName     xml.Name   `xml:"httpCompression"`
	Attrs       []xml.Attr `xml:",any,attr"`
	StaticTypes struct {
		XMLName xml.Name `xml:"staticTypes"`
	}
	DynamicTypes struct {
		XMLName xml.Name `xml:"dynamicTypes"`
	}
	UnwantedTags []xml.Name `xml:",any"`
}

func ValidateWebConfig(path string, writer io.Writer) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	conf := Configuration{}
	if err = xml.Unmarshal(data, &conf); err != nil {
		return err
	}

	hcAttrs := conf.SystemWebServer.HTTPCompression.Attrs
	if len(hcAttrs) > 0 {
		fmt.Fprintf(writer, "Warning: <httpCompression> should not have any attributes but it has %+v\n", collectAttrs(hcAttrs))
	}

	unwantedTags := conf.SystemWebServer.HTTPCompression.UnwantedTags
	if len(unwantedTags) > 0 {
		fmt.Fprintf(writer, "Warning: <httpCompression> should not have any child tags other than <staticTypes> and <dynamicTypes> but it has %+v\n", collectNames(unwantedTags))
	}
	return nil
}

func collectAttrs(attrs []xml.Attr) string {
	collected := make([]string, len(attrs))
	for i, attr := range attrs {
		collected[i] = attr.Name.Local
	}
	return strings.Join(collected, ", ")
}

func collectNames(names []xml.Name) string {
	collected := make([]string, len(names))
	for i, name := range names {
		collected[i] = fmt.Sprintf("<%s>", name.Local)
	}
	return strings.Join(collected, ", ")
}
