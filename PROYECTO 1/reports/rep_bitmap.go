package reports

// rep_bitmap.go

import (
	"MIA_P1_202400452/disk"
	"MIA_P1_202400452/ext2"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func ReportBMNode(diskPath string, partition disk.Partition, outPath string) error {
	sb, err := readEXT2SuperBlock(diskPath, partition)
	if err != nil {
		return err
	}

	bitmap, err := ext2.ReadBitmap(diskPath, sb.SBmInodeStart, sb.SInodesCount)
	if err != nil {
		return fmt.Errorf("error al leer bitmap de inodos: %v", err)
	}
	if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
		return fmt.Errorf("error creando directorio: %v", err)
	}

	var sbContent strings.Builder
	sbContent.WriteString("Bitmap de Inodos\n")
	sbContent.WriteString("================\n\n")

	for i := 0; i < len(bitmap); i += 20 {
		sbContent.WriteString(fmt.Sprintf("%4d: ", i+1))
		for j := 0; j < 20 && i+j < len(bitmap); j++ {
			sbContent.WriteString(fmt.Sprintf("%c", bitmap[i+j]))
		}
		sbContent.WriteString("\n")
	}

	if err := os.WriteFile(outPath, []byte(sbContent.String()), 0644); err != nil {
		return fmt.Errorf("error escribiendo archivo: %v", err)
	}
	return nil
}

func ReportBMBlock(diskPath string, partition disk.Partition, outPath string) error {
	sb, err := readEXT2SuperBlock(diskPath, partition)
	if err != nil {
		return err
	}

	bitmap, err := ext2.ReadBitmap(diskPath, sb.SBmBlockStart, sb.SBlocksCount)
	if err != nil {
		return fmt.Errorf("error al leer bitmap de bloques: %v", err)
	}

	if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
		return fmt.Errorf("error creando directorio: %v", err)
	}

	var sbContent strings.Builder
	sbContent.WriteString("Bitmap de Bloques\n")
	sbContent.WriteString("=================\n\n")

	for i := 0; i < len(bitmap); i += 20 {
		sbContent.WriteString(fmt.Sprintf("%4d: ", i+1))
		for j := 0; j < 20 && i+j < len(bitmap); j++ {
			sbContent.WriteString(fmt.Sprintf("%c", bitmap[i+j]))
		}
		sbContent.WriteString("\n")
	}

	if err := os.WriteFile(outPath, []byte(sbContent.String()), 0644); err != nil {
		return fmt.Errorf("error escribiendo archivo: %v", err)
	}
	return nil
}
