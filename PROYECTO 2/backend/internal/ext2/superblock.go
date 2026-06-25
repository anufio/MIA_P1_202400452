package ext2

// superblock.go

import (
	"bytes"
	"encoding/binary"
)

type SuperBlock struct {
	SFilesystemType  int32
	SInodesCount     int32
	SBlocksCount     int32
	SFreeBlocksCount int32
	SFreeInodesCount int32
	SMtime           [19]byte
	SUmtime          [19]byte
	SMntCount        int32
	SMagic           int32
	SInodeS          int32
	SBlockS          int32
	SFirstIno        int32
	SFirstBlo        int32
	SBmInodeStart    int32
	SBmBlockStart    int32
	SInodeStart      int32
	SBlockStart      int32
}

func SizeSuperBlock() int {
	return binary.Size(SuperBlock{})
}

func SuperBlockToBytes(sb SuperBlock) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, sb)
	return buf.Bytes(), err
}

func BytesToSuperBlock(data []byte) (SuperBlock, error) {
	var sb SuperBlock
	buf := bytes.NewReader(data)
	err := binary.Read(buf, binary.LittleEndian, &sb)
	return sb, err
}
