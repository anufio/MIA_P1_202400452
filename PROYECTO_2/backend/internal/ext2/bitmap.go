package ext2

// bitmap.go

import (
	"fmt"
	"os"
)

func ReadBitmap(diskPath string, startByte int32, count int32) ([]byte, error) {
	f, err := os.Open(diskPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	_, err = f.Seek(int64(startByte), 0)
	if err != nil {
		return nil, err
	}
	data := make([]byte, count)
	_, err = f.Read(data)
	return data, err
}

func WriteBitmap(diskPath string, startByte int32, data []byte) error {
	f, err := os.OpenFile(diskPath, os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Seek(int64(startByte), 0)
	if err != nil {
		return err
	}
	_, err = f.Write(data)
	return err
}

func GetFreeBit(bitmap []byte) int {
	for i, b := range bitmap {
		if b == '0' {
			return i
		}
	}
	return -1
}

func SetBit(diskPath string, startByte int32, index int32, value byte) error {
	if index < 0 {
		return fmt.Errorf("índice de bitmap inválido: %d", index)
	}
	f, err := os.OpenFile(diskPath, os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Seek(int64(startByte)+int64(index), 0)
	if err != nil {
		return err
	}
	_, err = f.Write([]byte{value})
	return err
}

func GetBit(diskPath string, startByte int32, index int32) (byte, error) {
	if index < 0 {
		return 0, fmt.Errorf("índice de bitmap inválido: %d", index)
	}
	f, err := os.Open(diskPath)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	_, err = f.Seek(int64(startByte)+int64(index), 0)
	if err != nil {
		return 0, err
	}
	buf := make([]byte, 1)
	_, err = f.Read(buf)
	return buf[0], err
}

func firstFreeAddress(diskPath string, bitmapStart, count, tableStart, structSize int32) int32 {
	bitmap, err := ReadBitmap(diskPath, bitmapStart, count)
	if err != nil {
		return -1
	}
	idx := GetFreeBit(bitmap)
	if idx == -1 {
		return -1
	}
	return tableStart + int32(idx)*structSize
}

func UpdateFirstFreeInode(diskPath string, sb *SuperBlock) {
	sb.SFirstIno = firstFreeAddress(diskPath, sb.SBmInodeStart, sb.SInodesCount, sb.SInodeStart, int32(SizeInode()))
}

func UpdateFirstFreeBlock(diskPath string, sb *SuperBlock) {
	sb.SFirstBlo = firstFreeAddress(diskPath, sb.SBmBlockStart, sb.SBlocksCount, sb.SBlockStart, int32(SizeBlock()))
}

func AllocateInode(diskPath string, sb *SuperBlock) (int32, error) {
	bitmap, err := ReadBitmap(diskPath, sb.SBmInodeStart, sb.SInodesCount)
	if err != nil {
		return -1, err
	}
	idx := GetFreeBit(bitmap)
	if idx == -1 {
		return -1, fmt.Errorf("no hay inodos libres")
	}
	err = SetBit(diskPath, sb.SBmInodeStart, int32(idx), '1')
	if err != nil {
		return -1, err
	}
	sb.SFreeInodesCount--
	UpdateFirstFreeInode(diskPath, sb)
	return int32(idx), nil
}

func AllocateBlock(diskPath string, sb *SuperBlock) (int32, error) {
	bitmap, err := ReadBitmap(diskPath, sb.SBmBlockStart, sb.SBlocksCount)
	if err != nil {
		return -1, err
	}
	idx := GetFreeBit(bitmap)
	if idx == -1 {
		return -1, fmt.Errorf("no hay bloques libres")
	}
	err = SetBit(diskPath, sb.SBmBlockStart, int32(idx), '1')
	if err != nil {
		return -1, err
	}
	sb.SFreeBlocksCount--
	UpdateFirstFreeBlock(diskPath, sb)
	return int32(idx), nil
}

func FreeInode(diskPath string, sb *SuperBlock, idx int32) error {
	bit, err := GetBit(diskPath, sb.SBmInodeStart, idx)
	if err != nil {
		return err
	}
	if bit == '1' {
		if err := SetBit(diskPath, sb.SBmInodeStart, idx, '0'); err != nil {
			return err
		}
		sb.SFreeInodesCount++
	}
	UpdateFirstFreeInode(diskPath, sb)
	return nil
}

func FreeBlock(diskPath string, sb *SuperBlock, idx int32) error {
	bit, err := GetBit(diskPath, sb.SBmBlockStart, idx)
	if err != nil {
		return err
	}
	if bit == '1' {
		if err := SetBit(diskPath, sb.SBmBlockStart, idx, '0'); err != nil {
			return err
		}
		sb.SFreeBlocksCount++
	}
	UpdateFirstFreeBlock(diskPath, sb)
	return nil
}

func InodeOffset(sb SuperBlock, idx int32) int64 {
	return int64(sb.SInodeStart) + int64(idx)*int64(SizeInode())
}

func BlockOffset(sb SuperBlock, idx int32) int64 {
	return int64(sb.SBlockStart) + int64(idx)*int64(SizeBlock())
}
