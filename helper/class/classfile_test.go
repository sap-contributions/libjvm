package class_test

import (
	"bytes"
	"encoding/binary"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/paketo-buildpacks/libjvm/helper/class"
	"github.com/sclevine/spec"
)

func testClassFile(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect
	)

	context("valid Java Class file", func() {
		m := bytes.NewBuffer(nil)
		binary.Write(m, binary.BigEndian, uint8(0xCA))
		binary.Write(m, binary.BigEndian, uint8(0xFE))
		binary.Write(m, binary.BigEndian, uint8(0xBA))
		binary.Write(m, binary.BigEndian, uint8(0xBE))
		binary.Write(m, binary.BigEndian, uint16(0))
		binary.Write(m, binary.BigEndian, uint16(55))

		it("returns correct Java version", func() {
			v, err := class.JVMVersionFromClassFile(m)
			Expect(err).ToNot(HaveOccurred())
			Expect(v).To(Equal("11"))
		})
	})

	context("invalid Java Class file", func() {
		var m *bytes.Buffer
		it.Before(func() {
			m = bytes.NewBuffer(nil)
			binary.Write(m, binary.BigEndian, uint8(0xCA))
			binary.Write(m, binary.BigEndian, uint8(0xFE))
			binary.Write(m, binary.BigEndian, uint8(0xBA))
			binary.Write(m, binary.BigEndian, uint8(0xBE))
			binary.Write(m, binary.BigEndian, uint16(0))
		})

		it("fails to read major version", func() {
			v, err := class.JVMVersionFromClassFile(m)
			Expect(v).To(BeEmpty())
			Expect(err).To(MatchError("failed to read major version from Java Class file, error: EOF"))
		})

		it("fails with the version unknown", func() {
			binary.Write(m, binary.BigEndian, uint16(100))
			v, err := class.JVMVersionFromClassFile(m)
			Expect(v).To(BeEmpty())
			Expect(err).To(MatchError("unknown Java version (hex): 0x64"))
		})
	})

}
