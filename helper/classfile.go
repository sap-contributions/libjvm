package helper

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

var versionsMap = map[uint16]string{
	0x2D: "1.1",
	0x2E: "1.2",
	0x2F: "1.3",
	0x30: "1.4",
	0x31: "5",
	0x32: "6",
	0x33: "7",
	0x34: "8",
	0x35: "9",
	0x36: "10",
	0x37: "11",
	0x38: "12",
	0x39: "13",
	0x3A: "14",
	0x3B: "15",
	0x3C: "16",
	0x3D: "17",
}

type magic struct {
	p1 uint8
	p2 uint8
	p3 uint8
	p4 uint8
}

func (m *magic) ensureMagicNumbers() error {
	err := errors.New("given file is not a Java Class file")
	if m.p1 != 0xCA {
		return err
	}

	if m.p2 != 0xFE {
		return err
	}

	if m.p3 != 0xBA {
		return err
	}

	if m.p4 != 0xBE {
		return err
	}

	return nil
}

type javaClass struct {
	magic *magic
	minor uint16
	major uint16
}

func (j *javaClass) hexToStringVersion() (string, error) {
	if v, ok := versionsMap[j.major]; ok {
		return v, nil
	}

	return "", fmt.Errorf("unknown Java version (hex): 0x%x", j.major)
}

type decoder struct {
	file io.Reader
	bo   binary.ByteOrder
	jc   *javaClass
}

func (d *decoder) readMagic() error {
	err := binary.Read(d.file, d.bo, &(d.jc.magic.p1))
	if err != nil {
		return fmt.Errorf("failed to read the first block of Java Class magic block, error: %s", err.Error())
	}

	err = binary.Read(d.file, d.bo, &(d.jc.magic.p2))
	if err != nil {
		return fmt.Errorf("failed to read the second block of Java Class magic block, error: %s", err.Error())
	}

	err = binary.Read(d.file, d.bo, &(d.jc.magic.p3))
	if err != nil {
		return fmt.Errorf("failed to read the third block of Java Class magic block, error: %s", err.Error())
	}

	err = binary.Read(d.file, d.bo, &(d.jc.magic.p4))
	if err != nil {
		return fmt.Errorf("failed to read the fourth block of Java Class magic block, error: %s", err.Error())
	}

	return d.jc.magic.ensureMagicNumbers()
}

func (d *decoder) readClass() error {
	err := d.readMagic()
	if err != nil {
		return err
	}
	err = binary.Read(d.file, d.bo, &(d.jc.minor))
	if err != nil {
		return fmt.Errorf("failed to read minor version from Java Class file, error: %s", err.Error())
	}
	err = binary.Read(d.file, d.bo, &(d.jc.major))
	if err != nil {
		return fmt.Errorf("failed to read major version from Java Class file, error: %s", err.Error())
	}

	return nil
}

func JVMVersionFromClassFile(classFile io.Reader) (string, error) {
	d := decoder{
		file: classFile,
		bo:   binary.BigEndian,
		jc: &javaClass{
			magic: &magic{},
		},
	}

	err := d.readClass()
	if err != nil {
		return "", err
	}

	return d.jc.hexToStringVersion()
}
