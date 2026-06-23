package ext2

// inode.go

import (
	"bytes"
	"encoding/binary"
)

type Inode struct {
	IUid   int32
	IGid   int32
	IS     int32
	IAtime [19]byte
	ICtime [19]byte
	IMtime [19]byte
	IBlock [15]int32
	IType  [1]byte
	IPerm  [3]byte
}

func SizeInode() int {
	return binary.Size(Inode{})
}

func InodeToBytes(inode Inode) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, inode)
	return buf.Bytes(), err
}

func BytesToInode(data []byte) (Inode, error) {
	var inode Inode
	buf := bytes.NewReader(data)
	err := binary.Read(buf, binary.LittleEndian, &inode)
	return inode, err
}

func NewInode() Inode {
	var inode Inode
	for i := range inode.IBlock {
		inode.IBlock[i] = -1
	}
	return inode
}
