package cos_test

import (
	"github.com/rancher-sandbox/cOS/tests/sut"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("cOS Recovery upgrade tests", func() {
	var s *sut.SUT
	BeforeEach(func() {
		s = sut.NewSUT()
		s.EventuallyConnects()
	})

	Context("upgrading recovery", func() {

		When("specific images", func() {
			It("upgrades to a specific image and reset back to the installed version", func() {
				out, err := s.Command("source /etc/os-release && echo $VERSION")
				Expect(err).ToNot(HaveOccurred())
				Expect(out).ToNot(Equal(""))

				version := out
				By("upgrading to raccos/releases-opensuse:cos-system-0.4.35-1")
				out, err = s.Command("cos-upgrade --recovery --docker-image raccos/releases-opensuse:cos-system-0.4.35-1")
				Expect(err).ToNot(HaveOccurred())
				Expect(out).Should(ContainSubstring("Upgrade done, now you might want to reboot"))
				Expect(out).Should(ContainSubstring("Upgrading recovery partition"))

				By("booting into recovery to check the OS version")
				err = s.ChangeBoot(sut.Recovery)
				Expect(err).ToNot(HaveOccurred())

				s.Reboot()

				out, err = s.Command("source /etc/os-release && echo $VERSION")
				Expect(err).ToNot(HaveOccurred())
				Expect(out).ToNot(Equal(""))
				Expect(out).ToNot(Equal(version))
				Expect(out).To(Equal("0.4.35+1\n"))

				By("setting back to active and rebooting")
				err = s.ChangeBoot(sut.Active)
				Expect(err).ToNot(HaveOccurred())

				s.Reboot()
			})
		})

		When("upgrade channel", func() {
			It("upgrades to latest image", func() {
				By("upgrading recovery and reboot")
				out, err := s.Command("cos-upgrade --recovery")
				Expect(err).ToNot(HaveOccurred())
				Expect(out).Should(ContainSubstring("Upgrade done, now you might want to reboot"))
				Expect(out).Should(ContainSubstring("Upgrading recovery partition"))

				err = s.ChangeBoot(sut.Recovery)
				Expect(err).ToNot(HaveOccurred())

				s.Reboot()

				By("checking recovery version")
				out, err = s.Command("source /etc/os-release && echo $VERSION")
				Expect(err).ToNot(HaveOccurred())
				Expect(out).ToNot(Equal(""))
				Expect(out).ToNot(Equal("0.4.35+1\n"))

				By("switch back to active and reboot")
				err = s.ChangeBoot(sut.Active)
				Expect(err).ToNot(HaveOccurred())

				s.Reboot()
			})
		})

	})
})