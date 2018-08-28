package lats_test

import (
	"github.com/cloudfoundry/sonde-go/events"
	"github.com/gogo/protobuf/proto"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Sending metrics through loggregator", func() {
	Describe("Firehose", func() {
		var (
			msgChan   <-chan *events.Envelope
			errorChan <-chan error
		)
		BeforeEach(func() {
			msgChan, errorChan = ConnectToFirehose()
		})

		AfterEach(func() {
			Expect(errorChan).To(BeEmpty())
		})

		It("receives a counter event with correct total", func() {
			envelope := createCounterEvent()
			EmitToMetronV1(envelope)

			receivedEnvelope := FindMatchingEnvelope(msgChan, envelope)
			Expect(receivedEnvelope).NotTo(BeNil())

			Expect(receivedEnvelope.GetCounterEvent().GetTotal()).To(Equal(envelope.GetCounterEvent().GetDelta()))

			// Clear it so the next assertion is valid
			receivedEnvelope.GetCounterEvent().Total = proto.Uint64(0)

			Expect(receivedEnvelope.GetCounterEvent()).To(Equal(envelope.GetCounterEvent()))
			EmitToMetronV1(envelope)

			receivedEnvelope = FindMatchingEnvelope(msgChan, envelope)
			Expect(receivedEnvelope).NotTo(BeNil())

			Expect(receivedEnvelope.GetCounterEvent().GetTotal()).To(Equal(uint64(10)))
		})

		It("receives a value metric", func() {
			envelope := createValueMetric()
			EmitToMetronV1(envelope)

			receivedEnvelope := FindMatchingEnvelope(msgChan, envelope)
			Expect(receivedEnvelope).NotTo(BeNil())

			Expect(receivedEnvelope.GetValueMetric()).To(Equal(envelope.GetValueMetric()))
		})

		It("receives a container metric", func() {
			envelope := createContainerMetric("7e6a0e79-4ac1-4521-95b7-a8e5ab5c7959")
			EmitToMetronV1(envelope)

			receivedEnvelope := FindMatchingEnvelope(msgChan, envelope)
			Expect(receivedEnvelope).NotTo(BeNil())

			Expect(receivedEnvelope.GetContainerMetric()).To(Equal(envelope.GetContainerMetric()))
		})
	})

	Describe("Stream", func() {
		It("receives a container metric", func() {
			msgChan, errorChan := ConnectToStream("7e6a0e79-4ac1-4521-95b7-a8e5ab5c7959")
			envelope := createContainerMetric("7e6a0e79-4ac1-4521-95b7-a8e5ab5c7959")
			EmitToMetronV1(createContainerMetric("7b446f11-9c05-4aa4-8e17-eb830300d0ed"))
			EmitToMetronV1(envelope)

			receivedEnvelope, err := FindMatchingEnvelopeByID("7e6a0e79-4ac1-4521-95b7-a8e5ab5c7959", msgChan)
			Expect(err).NotTo(HaveOccurred())
			Expect(receivedEnvelope).NotTo(BeNil())

			Expect(receivedEnvelope.GetContainerMetric()).To(Equal(envelope.GetContainerMetric()))
			Expect(errorChan).To(BeEmpty())
		})
	})
})
