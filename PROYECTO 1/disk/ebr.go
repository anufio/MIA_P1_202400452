package disk

// ebr.go

import (
	"bytes"
	"encoding/binary"
)

type EBR struct {
	PartMount [1]byte
	PartFit   [1]byte
	PartStart int32
	PartS     int32
	PartNext  int32
	PartName  [16]byte
}

func SizeEBR() int {
	return binary.Size(EBR{})
}

func EBRToBytes(ebr EBR) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, ebr)
	return buf.Bytes(), err
}

func BytesToEBR(data []byte) (EBR, error) {
	var ebr EBR
	buf := bytes.NewReader(data)
	err := binary.Read(buf, binary.LittleEndian, &ebr)
	return ebr, err
}
