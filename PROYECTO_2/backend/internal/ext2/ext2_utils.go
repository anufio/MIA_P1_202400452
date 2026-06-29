package ext2

// ext2_utils.go

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"os"
	"strings"
	"time"
)

func CalculateStructures(partitionSize int32) (int32, int32) {
	sbSize := float64(SizeSuperBlock())
	inodeSize := float64(SizeInode())
	blockSize := float64(SizeBlock())
	n := (float64(partitionSize) - sbSize) / (4.0 + inodeSize + 3.0*blockSize)
	numStructures := int32(math.Floor(n))
	return numStructures, numStructures * 3
}

func GetCurrentTime() [19]byte {
	now := time.Now()
	str := now.Format("2006-01-02 15:04:05")
	var arr [19]byte
	copy(arr[:], str)
	return arr
}

func WriteInode(diskPath string, sb SuperBlock, idx int32, inode Inode) error {
	f, err := os.OpenFile(diskPath, os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Seek(InodeOffset(sb, idx), 0)
	if err != nil {
		return err
	}
	return binary.Write(f, binary.LittleEndian, inode)
}

func ReadInode(diskPath string, sb SuperBlock, idx int32) (Inode, error) {
	var inode Inode
	f, err := os.Open(diskPath)
	if err != nil {
		return inode, err
	}
	defer f.Close()
	_, err = f.Seek(InodeOffset(sb, idx), 0)
	if err != nil {
		return inode, err
	}
	err = binary.Read(f, binary.LittleEndian, &inode)
	return inode, err
}

func WriteFolderBlock(diskPath string, sb SuperBlock, idx int32, fb FolderBlock) error {
	f, err := os.OpenFile(diskPath, os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Seek(BlockOffset(sb, idx), 0)
	if err != nil {
		return err
	}
	return binary.Write(f, binary.LittleEndian, fb)
}

func ReadFolderBlock(diskPath string, sb SuperBlock, idx int32) (FolderBlock, error) {
	var fb FolderBlock
	f, err := os.Open(diskPath)
	if err != nil {
		return fb, err
	}
	defer f.Close()
	_, err = f.Seek(BlockOffset(sb, idx), 0)
	if err != nil {
		return fb, err
	}
	err = binary.Read(f, binary.LittleEndian, &fb)
	return fb, err
}

func WriteFileBlock(diskPath string, sb SuperBlock, idx int32, fb FileBlock) error {
	f, err := os.OpenFile(diskPath, os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Seek(BlockOffset(sb, idx), 0)
	if err != nil {
		return err
	}
	return binary.Write(f, binary.LittleEndian, fb)
}

func ReadFileBlock(diskPath string, sb SuperBlock, idx int32) (FileBlock, error) {
	var fb FileBlock
	f, err := os.Open(diskPath)
	if err != nil {
		return fb, err
	}
	defer f.Close()
	_, err = f.Seek(BlockOffset(sb, idx), 0)
	if err != nil {
		return fb, err
	}
	err = binary.Read(f, binary.LittleEndian, &fb)
	return fb, err
}

func WritePointerBlock(diskPath string, sb SuperBlock, idx int32, pb PointerBlock) error {
	f, err := os.OpenFile(diskPath, os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Seek(BlockOffset(sb, idx), 0)
	if err != nil {
		return err
	}
	return binary.Write(f, binary.LittleEndian, pb)
}

func ReadPointerBlock(diskPath string, sb SuperBlock, idx int32) (PointerBlock, error) {
	var pb PointerBlock
	f, err := os.Open(diskPath)
	if err != nil {
		return pb, err
	}
	defer f.Close()
	_, err = f.Seek(BlockOffset(sb, idx), 0)
	if err != nil {
		return pb, err
	}
	err = binary.Read(f, binary.LittleEndian, &pb)
	return pb, err
}

func WriteSuperBlock(diskPath string, partStart int32, sb SuperBlock) error {
	f, err := os.OpenFile(diskPath, os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Seek(int64(partStart), 0)
	if err != nil {
		return err
	}
	return binary.Write(f, binary.LittleEndian, sb)
}

func ReadSuperBlock(diskPath string, partStart int32) (SuperBlock, error) {
	var sb SuperBlock
	f, err := os.Open(diskPath)
	if err != nil {
		return sb, err
	}
	defer f.Close()
	_, err = f.Seek(int64(partStart), 0)
	if err != nil {
		return sb, err
	}
	err = binary.Read(f, binary.LittleEndian, &sb)
	return sb, err
}

func InitBitmap(diskPath string, startByte int32, count int32) error {
	f, err := os.OpenFile(diskPath, os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Seek(int64(startByte), 0)
	if err != nil {
		return err
	}
	data := make([]byte, count)
	for i := range data {
		data[i] = '0'
	}
	_, err = f.Write(data)
	return err
}

func GetFileContent(diskPath string, sb SuperBlock, inodeIdx int32) (string, error) {
	inode, err := ReadInode(diskPath, sb, inodeIdx)
	if err != nil {
		return "", err
	}
	if inode.IType[0] != '1' {
		return "", fmt.Errorf("el inodo no es un archivo")
	}
	var content bytes.Buffer
	size := int(inode.IS)
	read := 0

	for i := 0; i < 12 && read < size; i++ {
		if inode.IBlock[i] == -1 {
			break
		}
		fb, err := ReadFileBlock(diskPath, sb, inode.IBlock[i])
		if err != nil {
			return content.String(), err
		}
		toRead := 64
		if read+toRead > size {
			toRead = size - read
		}
		content.Write(fb.BContent[:toRead])
		read += toRead
	}

	if inode.IBlock[12] != -1 && read < size {
		pb, err := ReadPointerBlock(diskPath, sb, inode.IBlock[12])
		if err != nil {
			return content.String(), err
		}
		for _, ptr := range pb.BPointers {
			if ptr == -1 || read >= size {
				break
			}
			fb, err := ReadFileBlock(diskPath, sb, ptr)
			if err != nil {
				return content.String(), err
			}
			toRead := 64
			if read+toRead > size {
				toRead = size - read
			}
			content.Write(fb.BContent[:toRead])
			read += toRead
		}
	}

	if inode.IBlock[13] != -1 && read < size {
		pb1, err := ReadPointerBlock(diskPath, sb, inode.IBlock[13])
		if err != nil {
			return content.String(), err
		}
		for _, ptr1 := range pb1.BPointers {
			if ptr1 == -1 || read >= size {
				break
			}
			pb2, err := ReadPointerBlock(diskPath, sb, ptr1)
			if err != nil {
				return content.String(), err
			}
			for _, ptr2 := range pb2.BPointers {
				if ptr2 == -1 || read >= size {
					break
				}
				fb, err := ReadFileBlock(diskPath, sb, ptr2)
				if err != nil {
					return content.String(), err
				}
				toRead := 64
				if read+toRead > size {
					toRead = size - read
				}
				content.Write(fb.BContent[:toRead])
				read += toRead
			}
		}
	}
	return content.String(), nil
}

func WriteFileContent(diskPath string, sb *SuperBlock, inode *Inode, content []byte) error {

	if err := FreeInodeBlocks(diskPath, sb, inode); err != nil {
		return err
	}

	size := len(content)
	written := 0
	directIdx := 0
	for written < size {
		remaining := size - written
		blockContent := make([]byte, 64)
		toWrite := 64
		if remaining < 64 {
			toWrite = remaining
		}
		copy(blockContent, content[written:written+toWrite])
		var fb FileBlock
		copy(fb.BContent[:], blockContent)
		blockIdx, err := AllocateBlock(diskPath, sb)
		if err != nil {
			return err
		}
		if err := WriteFileBlock(diskPath, *sb, blockIdx, fb); err != nil {
			return err
		}
		if directIdx < 12 {
			inode.IBlock[directIdx] = blockIdx
			directIdx++
		} else if directIdx < 12+16 {
			if inode.IBlock[12] == -1 {
				pbIdx, err := AllocateBlock(diskPath, sb)
				if err != nil {
					return err
				}
				pb := NewPointerBlock()
				WritePointerBlock(diskPath, *sb, pbIdx, pb)
				inode.IBlock[12] = pbIdx
			}
			pb, _ := ReadPointerBlock(diskPath, *sb, inode.IBlock[12])
			pb.BPointers[directIdx-12] = blockIdx
			WritePointerBlock(diskPath, *sb, inode.IBlock[12], pb)
			directIdx++
		} else if directIdx < 12+16+16*16 {

			if inode.IBlock[13] == -1 {
				pbIdx, err := AllocateBlock(diskPath, sb)
				if err != nil {
					return err
				}
				pb := NewPointerBlock()
				WritePointerBlock(diskPath, *sb, pbIdx, pb)
				inode.IBlock[13] = pbIdx
			}
			outerIdx := (directIdx - 12 - 16) / 16
			innerIdx := (directIdx - 12 - 16) % 16
			pb1, _ := ReadPointerBlock(diskPath, *sb, inode.IBlock[13])
			if pb1.BPointers[outerIdx] == -1 {
				pbIdx, err := AllocateBlock(diskPath, sb)
				if err != nil {
					return err
				}
				pb := NewPointerBlock()
				WritePointerBlock(diskPath, *sb, pbIdx, pb)
				pb1.BPointers[outerIdx] = pbIdx
				WritePointerBlock(diskPath, *sb, inode.IBlock[13], pb1)
			}
			pb2, _ := ReadPointerBlock(diskPath, *sb, pb1.BPointers[outerIdx])
			pb2.BPointers[innerIdx] = blockIdx
			WritePointerBlock(diskPath, *sb, pb1.BPointers[outerIdx], pb2)
			directIdx++
		}
		written += toWrite
	}
	inode.IS = int32(size)
	return nil
}

func FreeInodeBlocks(diskPath string, sb *SuperBlock, inode *Inode) error {

	for i := 0; i < 12; i++ {
		if inode.IBlock[i] != -1 {
			FreeBlock(diskPath, sb, inode.IBlock[i])
			inode.IBlock[i] = -1
		}
	}
	if inode.IBlock[12] != -1 {
		pb, err := ReadPointerBlock(diskPath, *sb, inode.IBlock[12])
		if err == nil {
			for _, ptr := range pb.BPointers {
				if ptr != -1 {
					FreeBlock(diskPath, sb, ptr)
				}
			}
		}
		FreeBlock(diskPath, sb, inode.IBlock[12])
		inode.IBlock[12] = -1
	}
	if inode.IBlock[13] != -1 {
		pb1, err := ReadPointerBlock(diskPath, *sb, inode.IBlock[13])
		if err == nil {
			for _, ptr1 := range pb1.BPointers {
				if ptr1 != -1 {
					pb2, err := ReadPointerBlock(diskPath, *sb, ptr1)
					if err == nil {
						for _, ptr2 := range pb2.BPointers {
							if ptr2 != -1 {
								FreeBlock(diskPath, sb, ptr2)
							}
						}
					}
					FreeBlock(diskPath, sb, ptr1)
				}
			}
		}
		FreeBlock(diskPath, sb, inode.IBlock[13])
		inode.IBlock[13] = -1
	}
	inode.IS = 0
	return nil
}

func FindInodeByPath(diskPath string, sb SuperBlock, path string) (int32, error) {
	path = strings.TrimSpace(path)
	if path == "" || path == "/" {
		return 0, nil
	}
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) == 1 && parts[0] == "" {
		return 0, nil
	}
	currentInode := int32(0)
	for _, part := range parts {
		if part == "" {
			continue
		}
		inode, err := ReadInode(diskPath, sb, currentInode)
		if err != nil {
			return -1, fmt.Errorf("error leyendo inodo %d: %v", currentInode, err)
		}
		if inode.IType[0] != '0' {
			return -1, fmt.Errorf("'%s' no es un directorio", part)
		}
		found := false
		for i := 0; i < 12; i++ {
			if inode.IBlock[i] == -1 {
				continue
			}
			fb, err := ReadFolderBlock(diskPath, sb, inode.IBlock[i])
			if err != nil {
				continue
			}
			for _, c := range fb.BContent {
				name := strings.TrimRight(string(c.BName[:]), "\\x00")
				name = strings.TrimRight(name, "\x00")
				if name == part && c.BInode != -1 {
					currentInode = c.BInode
					found = true
					break
				}
			}
			if found {
				break
			}
		}
		if !found {
			return -1, fmt.Errorf("no existe '%s' en la ruta", part)
		}
	}
	return currentInode, nil
}

func AddEntryToFolder(diskPath string, sb *SuperBlock, folderInodeIdx int32, name string, entryInodeIdx int32) error {
	inode, err := ReadInode(diskPath, *sb, folderInodeIdx)
	if err != nil {
		return err
	}
	var nameArr [12]byte
	copy(nameArr[:], name)
	for i := 0; i < 12; i++ {
		if inode.IBlock[i] == -1 {
			continue
		}
		fb, err := ReadFolderBlock(diskPath, *sb, inode.IBlock[i])
		if err != nil {
			continue
		}
		for j := range fb.BContent {
			if fb.BContent[j].BInode == -1 {
				fb.BContent[j].BName = nameArr
				fb.BContent[j].BInode = entryInodeIdx
				if err := WriteFolderBlock(diskPath, *sb, inode.IBlock[i], fb); err != nil {
					return err
				}
				inode.IMtime = GetCurrentTime()
				if err := WriteInode(diskPath, *sb, folderInodeIdx, inode); err != nil {
					return err
				}
				return nil
			}
		}
	}
	for i := 0; i < 12; i++ {
		if inode.IBlock[i] == -1 {
			blockIdx, err := AllocateBlock(diskPath, sb)
			if err != nil {
				return fmt.Errorf("no hay bloques libres para carpeta")
			}
			fb := NewFolderBlock()
			fb.BContent[0].BName = nameArr
			fb.BContent[0].BInode = entryInodeIdx
			if err := WriteFolderBlock(diskPath, *sb, blockIdx, fb); err != nil {
				return err
			}
			inode.IBlock[i] = blockIdx
			inode.IMtime = GetCurrentTime()
			if err := WriteInode(diskPath, *sb, folderInodeIdx, inode); err != nil {
				return err
			}
			return nil
		}
	}
	return fmt.Errorf("directorio lleno (sin implementar indirecto para carpetas)")
}

func RemoveEntryFromFolder(diskPath string, sb *SuperBlock, folderInodeIdx int32, name string) error {
	inode, err := ReadInode(diskPath, *sb, folderInodeIdx)
	if err != nil {
		return err
	}
	for i := 0; i < 12; i++ {
		if inode.IBlock[i] == -1 {
			continue
		}
		fb, err := ReadFolderBlock(diskPath, *sb, inode.IBlock[i])
		if err != nil {
			continue
		}
		for j := range fb.BContent {
			entryName := strings.TrimRight(string(fb.BContent[j].BName[:]), "\\x00")
			entryName = strings.TrimRight(entryName, "\x00")
			if entryName == name && fb.BContent[j].BInode != -1 {
				fb.BContent[j].BInode = -1
				fb.BContent[j].BName = [12]byte{}
				WriteFolderBlock(diskPath, *sb, inode.IBlock[i], fb)
				inode.IMtime = GetCurrentTime()
				WriteInode(diskPath, *sb, folderInodeIdx, inode)
				return nil
			}
		}
	}
	return fmt.Errorf("entrada '%s' no encontrada en el directorio", name)
}

func HasWritePermission(inode Inode, uid, gid int32) bool {
	perm := strings.TrimRight(string(inode.IPerm[:]), "\\x00")
	perm = strings.TrimRight(perm, "\x00")
	if len(perm) < 3 {
		return false
	}
	if uid == 1 {
		return true
	}
	uPerm := int(perm[0] - '0')
	gPerm := int(perm[1] - '0')
	oPerm := int(perm[2] - '0')
	if inode.IUid == uid {
		return uPerm&2 != 0
	} else if inode.IGid == gid {
		return gPerm&2 != 0
	} else {
		return oPerm&2 != 0
	}
}

func HasReadPermission(inode Inode, uid, gid int32) bool {
	perm := strings.TrimRight(string(inode.IPerm[:]), "\\x00")
	perm = strings.TrimRight(perm, "\x00")
	if len(perm) < 3 {
		return false
	}
	if uid == 1 {
		return true
	}
	uPerm := int(perm[0] - '0')
	gPerm := int(perm[1] - '0')
	oPerm := int(perm[2] - '0')
	if inode.IUid == uid {
		return uPerm&4 != 0
	} else if inode.IGid == gid {
		return gPerm&4 != 0
	} else {
		return oPerm&4 != 0
	}
}

func HasExecutePermission(inode Inode, uid, gid int32) bool {
	perm := strings.TrimRight(string(inode.IPerm[:]), "\\x00")
	perm = strings.TrimRight(perm, "\x00")
	if len(perm) < 3 {
		return false
	}
	if uid == 1 {
		return true
	}
	uPerm := int(perm[0] - '0')
	gPerm := int(perm[1] - '0')
	oPerm := int(perm[2] - '0')
	if inode.IUid == uid {
		return uPerm&1 != 0
	} else if inode.IGid == gid {
		return gPerm&1 != 0
	} else {
		return oPerm&1 != 0
	}
}
