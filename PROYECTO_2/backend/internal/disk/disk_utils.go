package disk

// disk_utils.go

import (
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func WriteMBR(path string, mbr MBR) error {
	f, err := os.OpenFile(path, os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	return binary.Write(f, binary.LittleEndian, mbr)
}

func ReadMBR(path string) (MBR, error) {
	var mbr MBR
	f, err := os.Open(path)
	if err != nil {
		return mbr, err
	}
	defer f.Close()
	err = binary.Read(f, binary.LittleEndian, &mbr)
	return mbr, err
}

func WriteEBR(path string, pos int64, ebr EBR) error {
	f, err := os.OpenFile(path, os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Seek(pos, 0)
	if err != nil {
		return err
	}
	return binary.Write(f, binary.LittleEndian, ebr)
}

func ReadEBR(path string, pos int64) (EBR, error) {
	var ebr EBR
	f, err := os.Open(path)
	if err != nil {
		return ebr, err
	}
	defer f.Close()
	_, err = f.Seek(pos, 0)
	if err != nil {
		return ebr, err
	}
	err = binary.Read(f, binary.LittleEndian, &ebr)
	return ebr, err
}

func WriteBytes(path string, pos int64, data []byte) error {
	f, err := os.OpenFile(path, os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Seek(pos, 0)
	if err != nil {
		return err
	}
	_, err = f.Write(data)
	return err
}

func ReadBytes(path string, pos int64, n int) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	_, err = f.Seek(pos, 0)
	if err != nil {
		return nil, err
	}
	data := make([]byte, n)
	_, err = f.Read(data)
	return data, err
}

func CreateDisk(path string, sizeBytes int64) error {

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("error creando directorios: %v", err)
	}
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("error creando archivo: %v", err)
	}
	defer f.Close()

	buf := make([]byte, 1024)
	written := int64(0)
	for written < sizeBytes {
		toWrite := int64(1024)
		if written+toWrite > sizeBytes {
			toWrite = sizeBytes - written
		}
		_, err = f.Write(buf[:toWrite])
		if err != nil {
			return fmt.Errorf("error escribiendo disco: %v", err)
		}
		written += toWrite
	}
	return nil
}

func DiskExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func GetDiskSize(path string) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

func WriteStruct(path string, pos int64, data interface{}) error {
	f, err := os.OpenFile(path, os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Seek(pos, 0)
	if err != nil {
		return err
	}
	return binary.Write(f, binary.LittleEndian, data)
}

func ReadStruct(path string, pos int64, data interface{}) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Seek(pos, 0)
	if err != nil {
		return err
	}
	return binary.Read(f, binary.LittleEndian, data)
}

func GetFitChar(fit string) (byte, error) {
	switch strings.ToUpper(fit) {
	case "BF":
		return 'B', nil
	case "FF":
		return 'F', nil
	case "WF":
		return 'W', nil
	default:
		return 0, fmt.Errorf("tipo de ajuste inválido: %s (use BF, FF o WF)", fit)
	}
}

func FindFreeSpace(mbr MBR, size int32, fit byte) int32 {
	mbrSize := int32(SizeMBR())
	diskSize := mbr.MbrTamano

	type segment struct{ start, end int32 }
	var used []segment

	for _, p := range mbr.MbrPartitions {
		if p.PartStart > 0 && p.PartS > 0 {
			used = append(used, segment{start: p.PartStart, end: p.PartStart + p.PartS})
		}
	}

	for i := 0; i < len(used)-1; i++ {
		for j := i + 1; j < len(used); j++ {
			if used[i].start > used[j].start {
				used[i], used[j] = used[j], used[i]
			}
		}
	}

	var free []struct{ start, size int32 }
	lastEnd := mbrSize

	for _, s := range used {
		if s.start > lastEnd {
			free = append(free, struct{ start, size int32 }{lastEnd, s.start - lastEnd})
		}
		if s.end > lastEnd {
			lastEnd = s.end
		}
	}
	if lastEnd < diskSize {
		free = append(free, struct{ start, size int32 }{lastEnd, diskSize - lastEnd})
	}

	if len(free) == 0 {
		return -1
	}

	switch fit {
	case 'F':
		for _, f := range free {
			if f.size >= size {
				return f.start
			}
		}
	case 'B':
		best := -1
		bestSize := int32(-1)
		for i, f := range free {
			if f.size >= size {
				if bestSize == -1 || f.size < bestSize {
					best = i
					bestSize = f.size
				}
			}
		}
		if best != -1 {
			return free[best].start
		}
	case 'W':
		worst := -1
		worstSize := int32(-1)
		for i, f := range free {
			if f.size >= size && f.size > worstSize {
				worst = i
				worstSize = f.size
			}
		}
		if worst != -1 {
			return free[worst].start
		}
	}
	return -1
}

func FindFreePartitionSlot(mbr MBR) int {
	for i, p := range mbr.MbrPartitions {
		if p.PartStart == 0 && p.PartS == 0 {
			return i
		}
	}
	return -1
}

func CountPrimaryExtended(mbr MBR) int {
	count := 0
	for _, p := range mbr.MbrPartitions {
		if p.PartStart > 0 {
			count++
		}
	}
	return count
}

func HasExtended(mbr MBR) bool {
	for _, p := range mbr.MbrPartitions {
		if p.PartType[0] == 'E' && p.PartStart > 0 {
			return true
		}
	}
	return false
}

func GetExtended(mbr MBR) (Partition, int, bool) {
	for i, p := range mbr.MbrPartitions {
		if p.PartType[0] == 'E' && p.PartStart > 0 {
			return p, i, true
		}
	}
	return Partition{}, -1, false
}

func FindPartitionByName(mbr MBR, name string) (Partition, int, bool) {
	for i, p := range mbr.MbrPartitions {
		if p.PartStart > 0 && strings.TrimRight(string(p.PartName[:]), "\x00") == name {
			return p, i, true
		}
	}
	return Partition{}, -1, false
}

func FindLogicalByName(diskPath string, ext Partition, name string) (EBR, int64, bool) {
	pos := int64(ext.PartStart)
	for pos != -1 {
		ebr, err := ReadEBR(diskPath, pos)
		if err != nil {
			break
		}
		ebrName := strings.TrimRight(string(ebr.PartName[:]), "\x00")
		if ebrName == name && ebr.PartS > 0 {
			return ebr, pos, true
		}
		if ebr.PartNext == -1 {
			break
		}
		pos = int64(ebr.PartNext)
	}
	return EBR{}, -1, false
}

func GetNextCorrelative(mbr MBR) int32 {
	max := int32(0)
	for _, p := range mbr.MbrPartitions {
		if p.PartCorrelative > max {
			max = p.PartCorrelative
		}
	}
	return max + 1
}
