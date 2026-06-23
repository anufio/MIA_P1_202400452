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
	return writeBitmapReport("Bitmap de Inodos", bitmap, outPath)
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
	return writeBitmapReport("Bitmap de Bloques", bitmap, outPath)

}

func writeBitmapReport(title string, bitmap []byte, outPath string) error {
	ext := strings.ToLower(filepath.Ext(outPath))

	if ext == ".txt" {
		return writeBitmapText(title, bitmap, outPath)
	}

	dot := buildBitmapDot(title, bitmap)
	return renderGraphviz(dot, outPath)
}

func writeBitmapText(title string, bitmap []byte, outPath string) error {
	if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
		return fmt.Errorf("error creando directorio: %v", err)
	}
	content := buildBitmapText(title, bitmap)

	if err := os.WriteFile(outPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("error escribiendo archivo: %v", err)
	}
	return nil

}

func buildBitmapText(title string, bitmap []byte) string {
	var text strings.Builder

	text.WriteString(title + "\n")
	text.WriteString(strings.Repeat("=", len(title)) + "\n\n")

	for i := 0; i < len(bitmap); i += 50 {
		text.WriteString(fmt.Sprintf("%5d: ", i+1))

		for j := 0; j < 50 && i+j < len(bitmap); j++ {
			text.WriteByte(bitmap[i+j])
		}

		text.WriteString("\n")
	}

	return text.String()
}

func buildBitmapDot(title string, bitmap []byte) string {
	var dot strings.Builder

	dot.WriteString("digraph G {\n")
	dot.WriteString("node [shape=plaintext]\n")
	dot.WriteString("rankdir=TB\n")
	dot.WriteString("BM [\n")
	dot.WriteString("label=<\n")
	dot.WriteString("<TABLE BORDER=\"0\" CELLBORDER=\"1\" CELLSPACING=\"0\" CELLPADDING=\"4\">\n")
	dot.WriteString(fmt.Sprintf("<TR><TD COLSPAN=\"2\" BGCOLOR=\"#3d85c8\"><FONT COLOR=\"white\"><B>%s</B></FONT></TD></TR>\n", title))
	dot.WriteString("<TR BGCOLOR=\"#cfe2f3\"><TD><B>Posición</B></TD><TD><B>Bits</B></TD></TR>\n")

	for i := 0; i < len(bitmap); i += 50 {
		var bits strings.Builder

		for j := 0; j < 50 && i+j < len(bitmap); j++ {
			bits.WriteByte(bitmap[i+j])
		}

		dot.WriteString(fmt.Sprintf(
			"<TR><TD>%d</TD><TD><FONT FACE=\"Courier\">%s</FONT></TD></TR>\n",
			i+1,
			bits.String(),
		))
	}
	dot.WriteString("</TABLE>>\n")
	dot.WriteString("]\n")
	dot.WriteString("}\n")

	return dot.String()
}
