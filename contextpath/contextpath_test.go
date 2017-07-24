package contextpath_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"code.cloudfoundry.org/hwc/contextpath"
	"github.com/cloudfoundry-community/go-cfenv"
)

var _ = Describe("Contextpath", func() {
	Describe("New", func() {
		createContextPath := func(URIs []string) (string, error) {
			cfapp := &cfenv.App{
				ApplicationURIs: URIs,
			}
			return contextpath.New(cfapp)
		}
		testContextPath := func(URIs []string, expectedPath string) {
			path, err := createContextPath(URIs)
			Expect(err).ToNot(HaveOccurred())
			Expect(path).To(Equal(expectedPath))
		}
		Context("no bound routes", func() {
			It("should have '/' context path", func() {
				testContextPath([]string{}, "/")
			})
		})
		Context("application URI without a path", func() {
			It("should have '/' context path", func() {
				testContextPath([]string{"myapp.apps.pcf.example.com"}, "/")
			})
		})
		Context("application URI with trailing slash and without a path", func() {
			It("should have '/' context path", func() {
				testContextPath([]string{"myapp.apps.pcf.example.com/"}, "/")
			})
		})
		Context("application URI with a /contextPath path", func() {
			It("should have '/contextPath' context path", func() {
				testContextPath([]string{"myapp.apps.pcf.example.com/contextPath"}, "/contextPath")
			})
		})
		Context("application URI with a /contextPath/contextPath2 path", func() {
			It("should have '/contextPath/contextPath2' context path", func() {
				testContextPath([]string{"myapp.apps.pcf.example.com/contextPath/contextPath2"},
					"/contextPath/contextPath2")
			})
		})
		Context("application URIs with the same /contextPath path", func() {
			It("should have '/contextPath' context path", func() {
				testContextPath([]string{
					"myapp.apps.pcf.example.com/contextPath",
					"myapp.apps.pcf.example.com/contextPath/",
					"example.com/contextPath"},
					"/contextPath")
			})
		})
		Context("application URIs with different paths", func() {
			It("should error", func() {
				_, err := createContextPath([]string{
					"myapp.apps.pcf.example.com/contextPath",
					"myapp.apps.pcf.example.com/contextPath1/contextPath2",
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal(
					"Application may not contain conflicting route paths: /contextPath, /contextPath1/contextPath2"))
			})
		})
	})
})
