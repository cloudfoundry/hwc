package hwcconfig_test

import (
	"io/ioutil"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "code.cloudfoundry.org/hwc/hwcconfig"
)

var _ = Describe("HwcApplication", func() {
	var defaultRootPath string
	var rootPath string
	BeforeEach(func() {
		var err error
		defaultRootPath, err = ioutil.TempDir("", "wwwroot")
		Expect(err).NotTo(HaveOccurred())
		rootPath, err = ioutil.TempDir("", "app")
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		os.RemoveAll(defaultRootPath)
		os.RemoveAll(rootPath)
	})

	Describe("New applications", func() {
		var apps []*HwcApplication
		Context("No context path", func() {
			BeforeEach(func() {
				apps = NewHwcApplications(
					defaultRootPath,
					rootPath,
					"/")
			})
			It("Creates one application", func() {
				Expect(apps).To(HaveLen(1))
			})
			It("creates the main app", func() {
				Expect(apps).To(ContainElement(&HwcApplication{
					Path:         "/",
					PhysicalPath: rootPath,
				}))
			})
		})
		Context("one context path", func() {
			BeforeEach(func() {
				apps = NewHwcApplications(
					defaultRootPath,
					rootPath,
					"/contextpath1")
			})
			It("creates 2 applications", func() {
				Expect(apps).To(HaveLen(2))
			})
			It("creates the main app", func() {
				Expect(apps).To(ContainElement(&HwcApplication{
					Path:         "/contextpath1",
					PhysicalPath: rootPath,
				}))
			})
			It("creates the root app", func() {
				Expect(apps).To(ContainElement(&HwcApplication{
					Path:         "/",
					PhysicalPath: defaultRootPath,
				}))
			})
		})
		Context("two nested context paths", func() {
			BeforeEach(func() {
				apps = NewHwcApplications(
					defaultRootPath,
					rootPath,
					"/contextpath1/contextpath2")
			})
			It("creates 3 applications", func() {
				Expect(apps).To(HaveLen(3))
			})
			It("creates the main app", func() {
				Expect(apps).To(ContainElement(&HwcApplication{
					Path:         "/contextpath1/contextpath2",
					PhysicalPath: rootPath,
				}))
			})
			It("creates the intermediate app", func() {
				Expect(apps).To(ContainElement(&HwcApplication{
					Path:         "/contextpath1",
					PhysicalPath: defaultRootPath,
				}))
			})
			It("creates the root app", func() {
				Expect(apps).To(ContainElement(&HwcApplication{
					Path:         "/",
					PhysicalPath: defaultRootPath,
				}))
			})
		})
	})
})
