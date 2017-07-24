package hwcconfig_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "code.cloudfoundry.org/hwc/hwcconfig"
)

var _ = Describe("HwcApplication", func() {
	Describe("New applications", func() {
		var apps []*HwcApplication
		Context("No context path", func() {
			BeforeEach(func() {
				apps = NewHwcApplications(
					"c:\\user\\vcap\\tmp\\wwwroot",
					"c:\\containerizer\\guid\\app",
					"/")
			})
			It("Creates one application", func() {
				Expect(apps).To(HaveLen(1))
			})
			It("creates the main app", func() {
				Expect(apps).To(ContainElement(&HwcApplication{
					Path:         "/",
					PhysicalPath: "c:\\containerizer\\guid\\app",
				}))
			})
		})
		Context("one context path", func() {
			BeforeEach(func() {
				apps = NewHwcApplications(
					"c:\\user\\vcap\\tmp\\wwwroot",
					"c:\\containerizer\\guid\\app",
					"/contextpath1")
			})
			It("creates 2 applications", func() {
				Expect(apps).To(HaveLen(2))
			})
			It("creates the main app", func() {
				Expect(apps).To(ContainElement(&HwcApplication{
					Path:         "/contextpath1",
					PhysicalPath: "c:\\containerizer\\guid\\app",
				}))
			})
			It("creates the root app", func() {
				Expect(apps).To(ContainElement(&HwcApplication{
					Path:         "/",
					PhysicalPath: "c:\\user\\vcap\\tmp\\wwwroot",
				}))
			})
		})
		Context("two nested context paths", func() {
			BeforeEach(func() {
				apps = NewHwcApplications(
					"c:\\user\\vcap\\tmp\\wwwroot",
					"c:\\containerizer\\guid\\app",
					"/contextpath1/contextpath2")
			})
			It("creates 3 applications", func() {
				Expect(apps).To(HaveLen(3))
			})
			It("creates the main app", func() {
				Expect(apps).To(ContainElement(&HwcApplication{
					Path:         "/contextpath1/contextpath2",
					PhysicalPath: "c:\\containerizer\\guid\\app",
				}))
			})
			It("creates the intermediate app", func() {
				Expect(apps).To(ContainElement(&HwcApplication{
					Path:         "/contextpath1",
					PhysicalPath: "c:\\user\\vcap\\tmp\\wwwroot",
				}))
			})
			It("creates the root app", func() {
				Expect(apps).To(ContainElement(&HwcApplication{
					Path:         "/",
					PhysicalPath: "c:\\user\\vcap\\tmp\\wwwroot",
				}))
			})
		})
	})
})
