package ext2

// block.go

import (
	"bytes"
	"encoding/binary"
)

type Content struct {
	BName  [12]byte
	BInode int32
}

type FolderBlock struct {
	BContent [4]Content
}

type FileBlock struct {
	BContent [64]byte
}

type PointerBlock struct {
	BPointers [16]int32
}

func SizeFolderBlock() int {
	return binary.Size(FolderBlock{})
}

func SizeFileBlock() int {
	return binary.Size(FileBlock{})
}

func SizePointerBlock() int {
	return binary.Size(PointerBlock{})
}

func SizeBlock() int {
	return 64
}

func FolderBlockToBytes(fb FolderBlock) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, fb)
	return buf.Bytes(), err
}

func BytesToFolderBlock(data []byte) (FolderBlock, error) {
	var fb FolderBlock
	buf := bytes.NewReader(data)
	err := binary.Read(buf, binary.LittleEndian, &fb)
	return fb, err
}

func FileBlockToBytes(fb FileBlock) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, fb)
	return buf.Bytes(), err
}

func BytesToFileBlock(data []byte) (FileBlock, error) {
	var fb FileBlock
	buf := bytes.NewReader(data)
	err := binary.Read(buf, binary.LittleEndian, &fb)
	return fb, err
}

func PointerBlockToBytes(pb PointerBlock) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, pb)
	return buf.Bytes(), err
}

func BytesToPointerBlock(data []byte) (PointerBlock, error) {
	var pb PointerBlock
	buf := bytes.NewReader(data)
	err := binary.Read(buf, binary.LittleEndian, &pb)
	return pb, err
}

func NewPointerBlock() PointerBlock {
	var pb PointerBlock
	for i := range pb.BPointers {
		pb.BPointers[i] = -1
	}
	return pb
}

func NewFolderBlock() FolderBlock {
	var fb FolderBlock
	for i := range fb.BContent {
		fb.BContent[i].BInode = -1
	}
	return fb
}
