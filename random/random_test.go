package random_test

import (
	"github.com/Altoros/cf-cassandra-broker/random"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("random", func() {
	Describe(".Bytes", func() {
		It("returns bytes array of given lenght", func() {
			Ω(len(random.Bytes(10))).To(Equal(10))
		})
	})

	Describe(".Hex", func() {
		It("returns string of 2x length", func() {
			Ω(len(random.Hex(10))).To(Equal(20))
		})
	})
})
