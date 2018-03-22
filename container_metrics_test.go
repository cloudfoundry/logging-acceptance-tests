package lats_test

import (
	"github.com/cloudfoundry/sonde-go/events"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Container Metrics Endpoint", func() {
	It("can receive container metrics", func() {
		envelope := createContainerMetric("7e6a0e79-4ac1-4521-95b7-a8e5ab5c7959")
		EmitToMetronV1(envelope)

		f := func() []*events.ContainerMetric {
			return RequestContainerMetrics("7e6a0e79-4ac1-4521-95b7-a8e5ab5c7959")
		}
		Eventually(f).Should(ContainElement(envelope.ContainerMetric))
	})
})
