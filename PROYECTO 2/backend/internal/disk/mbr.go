package disk

// mbr.go

import (
	"bytes"
	"encoding/binary"
	"time"
)

type Partition struct {
	PartStatus      [1]byte
	PartType        [1]byte
	PartFit         [1]byte
	PartStart       int32
	PartS           int32
	PartName        [16]byte
	PartCorrelative int32
	PartId          [4]byte
}

type MBR struct {
	MbrTamano        int32
	MbrFechaCreacion [19]byte
	MbrDskSignature  int32
	DskFit           [1]byte
	MbrPartitions    [6]Partition
}

func SizeMBR() int {
	return binary.Size(MBR{})
}

func SizePartition() int {
	return binary.Size(Partition{})
}

func MBRToBytes(mbr MBR) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, mbr)
	return buf.Bytes(), err
}

func BytesToMBR(data []byte) (MBR, error) {
	var mbr MBR
	buf := bytes.NewReader(data)
	err := binary.Read(buf, binary.LittleEndian, &mbr)
	return mbr, err
}

func GetFechaCreacion() [19]byte {
	now := time.Now()
	str := now.Format("2006-01-02 15:04:05")
	var arr [19]byte
	copy(arr[:], str)
	return arr
}

func PartitionToBytes(p Partition) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, p)
	return buf.Bytes(), err
}

func BytesToPartition(data []byte) (Partition, error) {
	var p Partition
	buf := bytes.NewReader(data)
	err := binary.Read(buf, binary.LittleEndian, &p)
	return p, err
}
