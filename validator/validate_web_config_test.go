package validator_test

import (
	"code.cloudfoundry.org/hwc/validator"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("ValidateWebConfig", func() {
	var (
		buf *gbytes.Buffer
	)

	BeforeEach(func() {
		buf = gbytes.NewBuffer()
	})

	Context("when the web.config is valid and all the xml elements are allowed", func() {
		BeforeEach(func() {
			webConfig := "../fixtures/webconfigs/Web.config.good"
			Expect(validator.ValidateWebConfig(webConfig, buf)).To(Succeed())
		})

		It("does not print any warnings", func() {
			Eventually(buf.Contents()).Should(BeEmpty())
		})
	})

	Context("when the web.config is valid but some xml elements are not allowed", func() {
		BeforeEach(func() {
			webConfig := "../fixtures/webconfigs/Web.config.bad"
			Expect(validator.ValidateWebConfig(webConfig, buf)).To(Succeed())
		})

		It("contains an error message about <httpCompression> attributes", func() {
			Eventually(buf).Should(gbytes.Say("Warning: <httpCompression> should not have any attributes but it has nastykey, anotherbadkey"))
		})

		It("contains an error message about <httpCompression> tags", func() {
			Eventually(buf).Should(gbytes.Say("Warning: <httpCompression> should not have any child tags other than <staticTypes>" +
				" and <dynamicTypes> but it has <scheme>"))
		})
	})

	Context("when the web.config does not exist", func() {
		It("returns an error", func() {
			webConfig := "some/file/that/does/not/exist"
			Expect(validator.ValidateWebConfig(webConfig, buf)).NotTo(Succeed())
		})
	})

	Context("when the web.config has invalid xml", func() {
		It("returns an error", func() {
			webConfig := "../fixtures/webconfigs/Web.config.invalid"
			Expect(validator.ValidateWebConfig(webConfig, buf)).NotTo(Succeed())
		})
	})
})
