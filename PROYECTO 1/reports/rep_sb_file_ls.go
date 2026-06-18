package reports

// rep_sb_file_ls.go

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"MIA_P1_202400452/disk"
	"MIA_P1_202400452/ext2"
)

func ReportSB(diskPath string, partition disk.Partition, outPath string) error {
	sb, err := readEXT2SuperBlock(diskPath, partition)
	if err != nil {
		return err
	}

	var dot strings.Builder
	dot.WriteString("digraph G {\n")
	dot.WriteString("node [shape=plaintext]\n")
	dot.WriteString("rankdir=TB\n")
	dot.WriteString("SB [\n")
	dot.WriteString("label=<\n")
	dot.WriteString("<TABLE BORDER=\"0\" CELLBORDER=\"1\" CELLSPACING=\"0\" CELLPADDING=\"4\">\n")
	dot.WriteString("<TR><TD COLSPAN=\"2\" BGCOLOR=\"#3d85c8\"><FONT COLOR=\"white\"><B>REPORTE DE SUPERBLOQUE</B></FONT></TD></TR>\n")
	dot.WriteString(fmt.Sprintf("<TR><TD>s_filesystem_type</TD><TD>%d</TD></TR>\n", sb.SFilesystemType))
	dot.WriteString(fmt.Sprintf("<TR><TD>s_inodes_count</TD><TD>%d</TD></TR>\n", sb.SInodesCount))
	dot.WriteString(fmt.Sprintf("<TR><TD>s_blocks_count</TD><TD>%d</TD></TR>\n", sb.SBlocksCount))
	dot.WriteString(fmt.Sprintf("<TR><TD>s_free_blocks_count</TD><TD>%d</TD></TR>\n", sb.SFreeBlocksCount))
	dot.WriteString(fmt.Sprintf("<TR><TD>s_free_inodes_count</TD><TD>%d</TD></TR>\n", sb.SFreeInodesCount))
	dot.WriteString(fmt.Sprintf("<TR><TD>s_mtime</TD><TD>%s</TD></TR>\n", strings.TrimRight(string(sb.SMtime[:]), "\x00")))
	dot.WriteString(fmt.Sprintf("<TR><TD>s_umtime</TD><TD>%s</TD></TR>\n", strings.TrimRight(string(sb.SUmtime[:]), "\x00")))
	dot.WriteString(fmt.Sprintf("<TR><TD>s_mnt_count</TD><TD>%d</TD></TR>\n", sb.SMntCount))
	dot.WriteString(fmt.Sprintf("<TR><TD>s_magic</TD><TD>0x%04X</TD></TR>\n", sb.SMagic))
	dot.WriteString(fmt.Sprintf("<TR><TD>s_inode_s</TD><TD>%d</TD></TR>\n", sb.SInodeS))
	dot.WriteString(fmt.Sprintf("<TR><TD>s_block_s</TD><TD>%d</TD></TR>\n", sb.SBlockS))
	dot.WriteString(fmt.Sprintf("<TR><TD>s_first_ino</TD><TD>%d</TD></TR>\n", sb.SFirstIno))
	dot.WriteString(fmt.Sprintf("<TR><TD>s_first_blo</TD><TD>%d</TD></TR>\n", sb.SFirstBlo))
	dot.WriteString(fmt.Sprintf("<TR><TD>s_bm_inode_start</TD><TD>%d</TD></TR>\n", sb.SBmInodeStart))
	dot.WriteString(fmt.Sprintf("<TR><TD>s_bm_block_start</TD><TD>%d</TD></TR>\n", sb.SBmBlockStart))
	dot.WriteString(fmt.Sprintf("<TR><TD>s_inode_start</TD><TD>%d</TD></TR>\n", sb.SInodeStart))
	dot.WriteString(fmt.Sprintf("<TR><TD>s_block_start</TD><TD>%d</TD></TR>\n", sb.SBlockStart))
	dot.WriteString("</TABLE>>\n")
	dot.WriteString("]\n")
	dot.WriteString("}\n")

	return renderGraphviz(dot.String(), outPath)
}

func ReportFILE(diskPath string, partition disk.Partition, outPath string, filePath string) error {
	sb, err := readEXT2SuperBlock(diskPath, partition)
	if err != nil {
		return err
	}

	inodeIdx, err := ext2.FindInodeByPath(diskPath, sb, filePath)
	if err != nil {
		return fmt.Errorf("el archivo '%s' no existe: %v", filePath, err)
	}

	inode, err := ext2.ReadInode(diskPath, sb, inodeIdx)
	if err != nil {
		return fmt.Errorf("error al leer inodo: %v", err)
	}

	if inode.IType[0] != '1' {
		return fmt.Errorf("'%s' no es un archivo", filePath)
	}

	content, err := ext2.GetFileContent(diskPath, sb, inodeIdx)
	if err != nil {
		return fmt.Errorf("error al leer contenido: %v", err)
	}

	if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
		return fmt.Errorf("error creando directorio: %v", err)
	}

	baseName := filepath.Base(filePath)
	output := fmt.Sprintf("%s:\n%s", baseName, content)

	if err := os.WriteFile(outPath, []byte(output), 0644); err != nil {
		return fmt.Errorf("error escribiendo archivo: %v", err)
	}

	return nil
}

func ReportLS(diskPath string, partition disk.Partition, outPath string, pathFileLs string) error {
	sb, err := readEXT2SuperBlock(diskPath, partition)
	if err != nil {
		return err
	}

	var dirInodeIdx int32
	if pathFileLs == "" || pathFileLs == "/" {
		dirInodeIdx = 0
	} else {
		dirInodeIdx, err = ext2.FindInodeByPath(diskPath, sb, pathFileLs)
		if err != nil {
			return fmt.Errorf("el directorio '%s' no existe: %v", pathFileLs, err)
		}
	}

	dirInode, err := ext2.ReadInode(diskPath, sb, dirInodeIdx)
	if err != nil {
		return fmt.Errorf("error al leer inodo del directorio: %v", err)
	}

	if dirInode.IType[0] != '0' {
		return fmt.Errorf("'%s' no es un directorio", pathFileLs)
	}

	var dot strings.Builder
	userNames, groupNames := readUserGroupNames(diskPath, sb)
	dot.WriteString("digraph G {\n")
	dot.WriteString("node [shape=plaintext]\n")
	dot.WriteString("rankdir=TB\n")
	dot.WriteString("LS [\n")
	dot.WriteString("label=<\n")
	dot.WriteString("<TABLE BORDER=\"0\" CELLBORDER=\"1\" CELLSPACING=\"0\" CELLPADDING=\"4\">\n")
	dot.WriteString("<TR><TD COLSPAN=\"9\" BGCOLOR=\"#3d85c8\"><FONT COLOR=\"white\"><B>REPORTE LS</B></FONT></TD></TR>\n")
	dot.WriteString("<TR BGCOLOR=\"#6d9eeb\"><TD><B>Permisos</B></TD><TD><B>Owner</B></TD><TD><B>Grupo</B></TD><TD><B>Size</B></TD><TD><B>Mod Fecha</B></TD><TD><B>Mod Hora</B></TD><TD><B>Creación</B></TD><TD><B>Tipo</B></TD><TD><B>Name</B></TD></TR>\n")

	entries := getDirEntries(diskPath, sb, dirInode)

	for _, entry := range entries {
		if entry.BInode == -1 {
			continue
		}

		entryInode, err := ext2.ReadInode(diskPath, sb, entry.BInode)
		if err != nil {
			continue
		}

		name := strings.TrimRight(string(entry.BName[:]), "\x00")
		if name == "." || name == ".." {
			continue
		}

		permStr := formatPermissions(entryInode.IPerm)
		owner := htmlEscape(nameOrID(userNames, entryInode.IUid))
		group := htmlEscape(nameOrID(groupNames, entryInode.IGid))
		size := fmt.Sprintf("%d", entryInode.IS)

		mtime := strings.TrimRight(string(entryInode.IMtime[:]), "\x00")
		date, timeStr := splitDateTime(mtime)
		ctime := strings.TrimRight(string(entryInode.ICtime[:]), "\x00")

		var tipo string
		if entryInode.IType[0] == '0' {
			tipo = "Carpeta"
		} else {
			tipo = "Archivo"
		}

		dot.WriteString(fmt.Sprintf("<TR><TD>%s</TD><TD>%s</TD><TD>%s</TD><TD>%s</TD><TD>%s</TD><TD>%s</TD><TD>%s</TD><TD>%s</TD><TD>%s</TD></TR>\n",
			permStr, owner, group, size, htmlEscape(date), htmlEscape(timeStr), htmlEscape(ctime), tipo, htmlEscape(name)))
	}

	dot.WriteString("</TABLE>>\n")
	dot.WriteString("]\n")
	dot.WriteString("}\n")

	return renderGraphviz(dot.String(), outPath)
}

func getDirEntries(diskPath string, sb ext2.SuperBlock, inode ext2.Inode) []ext2.Content {
	var entries []ext2.Content
	for i := 0; i < 12; i++ {
		if inode.IBlock[i] == -1 {
			continue
		}
		fb, err := ext2.ReadFolderBlock(diskPath, sb, inode.IBlock[i])
		if err != nil {
			continue
		}
		for _, c := range fb.BContent {
			entries = append(entries, c)
		}
	}
	return entries
}

func formatPermissions(perm [3]byte) string {
	permStr := strings.TrimRight(string(perm[:]), "\x00")
	if len(permStr) < 3 {
		return "---------"
	}

	var result strings.Builder
	result.WriteString("-")

	for i := 0; i < 3; i++ {
		p := permStr[i] - '0'
		if p&4 != 0 {
			result.WriteString("r")
		} else {
			result.WriteString("-")
		}
		if p&2 != 0 {
			result.WriteString("w")
		} else {
			result.WriteString("-")
		}
		if p&1 != 0 {
			result.WriteString("x")
		} else {
			result.WriteString("-")
		}
	}

	return result.String()
}

func splitDateTime(datetime string) (string, string) {
	parts := strings.Split(datetime, " ")
	if len(parts) >= 2 {
		return parts[0], parts[1]
	}
	return datetime, ""
}
