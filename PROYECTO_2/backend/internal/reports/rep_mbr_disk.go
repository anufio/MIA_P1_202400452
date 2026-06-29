package reports

// rep_mbr_disk.go

import (
	"MIA_P2_202400452/internal/disk"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

func renderGraphviz(dotContent, outPath string) error {

	dir := filepath.Dir(outPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("error creando directorio: %v", err)
	}

	tmpFile, err := os.CreateTemp("", "report_*.dot")
	if err != nil {
		return fmt.Errorf("error creando archivo temporal: %v", err)
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath)

	if _, err := tmpFile.WriteString(dotContent); err != nil {
		tmpFile.Close()
		return fmt.Errorf("error escribiendo archivo DOT: %v", err)
	}
	tmpFile.Close()

	format := graphvizFormat(outPath)
	cmd := exec.Command("dot", "-T"+format, tmpPath, "-o", outPath)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("error ejecutando dot: %v\n%s", err, string(output))
	}

	return nil
}

func graphvizFormat(outPath string) string {
	ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(outPath)), ".")
	switch ext {
	case "jpg", "jpeg":
		return "jpg"
	case "png", "pdf", "svg":
		return ext
	default:
		return "jpg"
	}
}

func ReportMBR(diskPath, outPath string) error {
	mbr, err := disk.ReadMBR(diskPath)
	if err != nil {
		return fmt.Errorf("error al leer MBR: %v", err)
	}

	var sb strings.Builder
	sb.WriteString("digraph G {\n")
	sb.WriteString(" node [shape=plaintext]\n")
	sb.WriteString(" rankdir=TB\n")
	sb.WriteString(" MBR [\n")
	sb.WriteString(" label=<\n")
	sb.WriteString(" <TABLE BORDER=\"0\" CELLBORDER=\"1\" CELLSPACING=\"0\" CELLPADDING=\"4\">\n")
	sb.WriteString(" <TR><TD COLSPAN=\"2\" BGCOLOR=\"#3d85c8\"><FONT COLOR=\"white\"><B>REPORTE DE MBR</B></FONT></TD></TR>\n")
	sb.WriteString(fmt.Sprintf(" <TR><TD>mbr_tamano</TD><TD>%d</TD></TR>\n", mbr.MbrTamano))
	sb.WriteString(fmt.Sprintf(" <TR><TD>mbr_fecha_creacion</TD><TD>%s</TD></TR>\n", strings.TrimRight(string(mbr.MbrFechaCreacion[:]), "\x00")))
	sb.WriteString(fmt.Sprintf(" <TR><TD>mbr_dsk_signature</TD><TD>%d</TD></TR>\n", mbr.MbrDskSignature))
	sb.WriteString(fmt.Sprintf(" <TR><TD>dsk_fit</TD><TD>%s</TD></TR>\n", string(mbr.DskFit[:])))

	for i, p := range mbr.MbrPartitions {
		if p.PartStart == 0 && p.PartS == 0 {
			continue
		}
		typeStr := string(p.PartType[:])
		if typeStr == "E" {
			sb.WriteString(fmt.Sprintf(" <TR><TD COLSPAN=\"2\" BGCOLOR=\"#f4cccc\"><B>Partición Extendida %d</B></TD></TR>\n", i+1))
		} else {
			sb.WriteString(fmt.Sprintf(" <TR><TD COLSPAN=\"2\" BGCOLOR=\"#6d9eeb\"><B>Partición %d</B></TD></TR>\n", i+1))
		}
		sb.WriteString(fmt.Sprintf(" <TR><TD>part_status</TD><TD>%s</TD></TR>\n", string(p.PartStatus[:])))
		sb.WriteString(fmt.Sprintf(" <TR><TD>part_type</TD><TD>%s</TD></TR>\n", typeStr))
		sb.WriteString(fmt.Sprintf(" <TR><TD>part_fit</TD><TD>%s</TD></TR>\n", string(p.PartFit[:])))
		sb.WriteString(fmt.Sprintf(" <TR><TD>part_start</TD><TD>%d</TD></TR>\n", p.PartStart))
		sb.WriteString(fmt.Sprintf(" <TR><TD>part_size</TD><TD>%d</TD></TR>\n", p.PartS))
		sb.WriteString(fmt.Sprintf(" <TR><TD>part_name</TD><TD>%s</TD></TR>\n", strings.TrimRight(string(p.PartName[:]), "\x00")))
		sb.WriteString(fmt.Sprintf(" <TR><TD>part_correlative</TD><TD>%d</TD></TR>\n", p.PartCorrelative))
		sb.WriteString(fmt.Sprintf(" <TR><TD>part_id</TD><TD>%s</TD></TR>\n", strings.TrimRight(string(p.PartId[:]), "\x00")))
		if p.PartType[0] == 'E' {
			appendEBRRows(&sb, diskPath, p)
		}
	}
	sb.WriteString(" </TABLE>>\n")
	sb.WriteString(" ]\n")
	sb.WriteString("}\n")

	return renderGraphviz(sb.String(), outPath)
}

func appendEBRRows(sb *strings.Builder, diskPath string, ext disk.Partition) {
	pos := ext.PartStart
	extEnd := ext.PartStart + ext.PartS
	visited := make(map[int32]bool)
	ebrIndex := 1

	for pos >= ext.PartStart && pos < extEnd && !visited[pos] {
		visited[pos] = true
		ebr, err := disk.ReadEBR(diskPath, int64(pos))
		if err != nil {
			return
		}

		if ebr.PartS > 0 {
			sb.WriteString(fmt.Sprintf(" <TR><TD COLSPAN=\"2\" BGCOLOR=\"#fce5cd\"><B>EBR %d</B></TD></TR>\n", ebrIndex))
			sb.WriteString(fmt.Sprintf(" <TR><TD>ebr_position</TD><TD>%d</TD></TR>\n", pos))
			sb.WriteString(fmt.Sprintf(" <TR><TD>part_mount</TD><TD>%s</TD></TR>\n", string(ebr.PartMount[:])))
			sb.WriteString(fmt.Sprintf(" <TR><TD>part_fit</TD><TD>%s</TD></TR>\n", string(ebr.PartFit[:])))
			sb.WriteString(fmt.Sprintf(" <TR><TD>part_start</TD><TD>%d</TD></TR>\n", ebr.PartStart))
			sb.WriteString(fmt.Sprintf(" <TR><TD>part_size</TD><TD>%d</TD></TR>\n", ebr.PartS))
			sb.WriteString(fmt.Sprintf(" <TR><TD>part_next</TD><TD>%d</TD></TR>\n", ebr.PartNext))
			sb.WriteString(fmt.Sprintf(" <TR><TD>part_name</TD><TD>%s</TD></TR>\n", strings.TrimRight(string(ebr.PartName[:]), "\x00")))
			ebrIndex++
		}

		if ebr.PartNext == -1 {
			return
		}
		pos = ebr.PartNext
	}
}

func ReportDISK(diskPath, outPath string) error {
	mbr, err := disk.ReadMBR(diskPath)
	if err != nil {
		return fmt.Errorf("error al leer MBR: %v", err)
	}

	diskName := filepath.Base(diskPath)
	totalSize := float64(mbr.MbrTamano)
	mbrSizeBytes := float64(disk.SizeMBR())

	var sb strings.Builder
	sb.WriteString("digraph G {\n")
	sb.WriteString(" node [shape=plaintext]\n")
	sb.WriteString(" rankdir=LR\n")
	sb.WriteString(fmt.Sprintf(" label=\"%s\"\n", diskName))
	sb.WriteString(" labelloc=t\n")
	sb.WriteString(" DISK [\n")
	sb.WriteString(" label=<\n")
	sb.WriteString(" <TABLE BORDER=\"1\" CELLBORDER=\"1\" CELLSPACING=\"0\" CELLPADDING=\"8\">\n")
	sb.WriteString(" <TR>\n")

	mbrPct := (mbrSizeBytes / totalSize) * 100
	sb.WriteString(fmt.Sprintf(" <TD BGCOLOR=\"#cfe2f3\"><B>MBR</B><BR/>%.1f%% del disco</TD>\n", mbrPct))

	type partInfo struct {
		start, end int32
		name       string
		isExt      bool
		partition  disk.Partition
	}
	var parts []partInfo

	for _, p := range mbr.MbrPartitions {
		if p.PartStart > 0 && p.PartS > 0 {
			parts = append(parts, partInfo{
				start:     p.PartStart,
				end:       p.PartStart + p.PartS,
				name:      strings.TrimRight(string(p.PartName[:]), "\x00"),
				isExt:     p.PartType[0] == 'E',
				partition: p,
			})
		}
	}

	sort.Slice(parts, func(i, j int) bool {
		return parts[i].start < parts[j].start
	})

	lastEnd := int32(disk.SizeMBR())
	for _, p := range parts {
		if p.start > lastEnd {
			freeSize := p.start - lastEnd
			freePct := (float64(freeSize) / totalSize) * 100
			sb.WriteString(fmt.Sprintf(" <TD BGCOLOR=\"#ffffff\">Libre<BR/>%.1f%%</TD>\n", freePct))
		}
		pct := (float64(p.end-p.start) / totalSize) * 100
		if p.isExt {
			sb.WriteString(" <TD BGCOLOR=\"#f4cccc\">\n")
			sb.WriteString("  <TABLE BORDER=\"0\" CELLBORDER=\"1\" CELLSPACING=\"0\" CELLPADDING=\"5\">\n")
			sb.WriteString(fmt.Sprintf("   <TR><TD COLSPAN=\"10\"><B>Extendida</B><BR/>%s<BR/>%.1f%%</TD></TR>\n", htmlEscape(p.name), pct))
			appendLogicalDiskCells(&sb, diskPath, p.partition, totalSize)
			sb.WriteString("  </TABLE>\n")
			sb.WriteString(" </TD>\n")
		} else {
			sb.WriteString(fmt.Sprintf(" <TD BGCOLOR=\"#d9ead3\"><B>Primaria</B><BR/>%s<BR/>%.1f%%</TD>\n", htmlEscape(p.name), pct))
		}
		lastEnd = p.end
	}

	if lastEnd < mbr.MbrTamano {
		freeSize := mbr.MbrTamano - lastEnd
		freePct := (float64(freeSize) / totalSize) * 100
		sb.WriteString(fmt.Sprintf(" <TD BGCOLOR=\"#ffffff\">Libre<BR/>%.1f%%</TD>\n", freePct))
	}

	sb.WriteString(" </TR>\n")
	sb.WriteString(" </TABLE>>\n")
	sb.WriteString(" ]\n")
	sb.WriteString("}\n")

	return renderGraphviz(sb.String(), outPath)
}

type logicalDiskSegment struct {
	start int32
	end   int32
	name  string
	kind  string
}

func appendLogicalDiskCells(sb *strings.Builder, diskPath string, ext disk.Partition, totalSize float64) {
	segments := logicalDiskSegments(diskPath, ext)
	if len(segments) == 0 {
		sb.WriteString("   <TR><TD>Libre<BR/>100.0% extendida</TD></TR>\n")
		return
	}
	sb.WriteString("   <TR>\n")
	for _, segment := range segments {
		pct := (float64(segment.end-segment.start) / totalSize) * 100
		switch segment.kind {
		case "ebr":
			sb.WriteString(fmt.Sprintf("    <TD BGCOLOR=\"#fce5cd\"><B>EBR</B><BR/>%.1f%%</TD>\n", pct))
		case "logical":
			sb.WriteString(fmt.Sprintf("    <TD BGCOLOR=\"#eadcf8\"><B>Lógica</B><BR/>%s<BR/>%.1f%%</TD>\n", htmlEscape(segment.name), pct))
		default:
			sb.WriteString(fmt.Sprintf("    <TD BGCOLOR=\"#ffffff\">Libre<BR/>%.1f%%</TD>\n", pct))
		}
	}
	sb.WriteString("   </TR>\n")
}

func logicalDiskSegments(diskPath string, ext disk.Partition) []logicalDiskSegment {
	var records []struct {
		pos int32
		ebr disk.EBR
	}
	pos := ext.PartStart
	extEnd := ext.PartStart + ext.PartS
	visited := make(map[int32]bool)
	for pos >= ext.PartStart && pos < extEnd && !visited[pos] {
		visited[pos] = true
		ebr, err := disk.ReadEBR(diskPath, int64(pos))
		if err != nil {
			break
		}
		if ebr.PartS > 0 {
			records = append(records, struct {
				pos int32
				ebr disk.EBR
			}{pos: pos, ebr: ebr})
		}
		if ebr.PartNext == -1 {
			break
		}
		pos = ebr.PartNext
	}
	sort.Slice(records, func(i, j int) bool {
		return records[i].pos < records[j].pos
	})

	var segments []logicalDiskSegment
	lastEnd := ext.PartStart
	ebrSize := int32(disk.SizeEBR())
	for _, rec := range records {
		if rec.pos > lastEnd {
			segments = append(segments, logicalDiskSegment{start: lastEnd, end: rec.pos, kind: "free"})
		}
		ebrEnd := rec.pos + ebrSize
		if ebrEnd > extEnd {
			ebrEnd = extEnd
		}
		if rec.pos < ebrEnd {
			segments = append(segments, logicalDiskSegment{start: rec.pos, end: ebrEnd, kind: "ebr"})
		}
		logEnd := rec.ebr.PartStart + rec.ebr.PartS
		if logEnd > extEnd {
			logEnd = extEnd
		}
		if rec.ebr.PartStart < logEnd {
			segments = append(segments, logicalDiskSegment{
				start: rec.ebr.PartStart,
				end:   logEnd,
				name:  strings.TrimRight(string(rec.ebr.PartName[:]), "\x00"),
				kind:  "logical",
			})
		}
		if logEnd > lastEnd {
			lastEnd = logEnd
		}
	}
	if lastEnd < extEnd {
		segments = append(segments, logicalDiskSegment{start: lastEnd, end: extEnd, kind: "free"})
	}
	return segments
}

func htmlEscape(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	return s
}
