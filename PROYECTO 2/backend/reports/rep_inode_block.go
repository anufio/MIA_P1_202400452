package reports

// rep_inode_block.go
import (
	"fmt"
	"strings"

	"MIA_P1_202400452/disk"
	"MIA_P1_202400452/ext2"
)

func ReportINODE(diskPath string, partition disk.Partition, outPath string) error {
	sb, err := readEXT2SuperBlock(diskPath, partition)
	if err != nil {
		return err
	}

	inodeBitmap, err := ext2.ReadBitmap(diskPath, sb.SBmInodeStart, sb.SInodesCount)
	if err != nil {
		return fmt.Errorf("error al leer bitmap de inodos: %v", err)
	}

	var dot strings.Builder
	dot.WriteString("digraph G {\n")
	dot.WriteString("rankdir=LR;\n")
	dot.WriteString("node [shape=plaintext];\n")
	dot.WriteString("label=\"Reporte INODE - Ana Lucía Nufio Roblero - 202400452\";\n")
	dot.WriteString("labelloc=t;\n\n")

	for i := int32(0); i < sb.SInodesCount && int(i) < len(inodeBitmap); i++ {
		if inodeBitmap[i] != '1' {
			continue
		}

		inode, err := ext2.ReadInode(diskPath, sb, i)
		if err != nil {
			continue
		}

		writeInodeGraphNode(&dot, i, inode)
	}

	for i := int32(0); i < sb.SInodesCount && int(i) < len(inodeBitmap); i++ {
		if inodeBitmap[i] != '1' {
			continue
		}

		inode, err := ext2.ReadInode(diskPath, sb, i)
		if err != nil || inode.IType[0] != '0' {
			continue
		}

		writeFolderInodeEdges(&dot, diskPath, sb, i, inode, inodeBitmap)
	}

	dot.WriteString("}\n")

	return renderGraphviz(dot.String(), outPath)
}

func ReportBLOCK(diskPath string, partition disk.Partition, outPath string) error {
	sb, err := readEXT2SuperBlock(diskPath, partition)
	if err != nil {
		return err
	}

	blockBitmap, err := ext2.ReadBitmap(diskPath, sb.SBmBlockStart, sb.SBlocksCount)
	if err != nil {
		return fmt.Errorf("error al leer bitmap de bloques: %v", err)
	}

	blockKinds := collectBlockKinds(diskPath, sb, blockBitmap)

	var dot strings.Builder
	dot.WriteString("digraph G {\n")
	dot.WriteString("rankdir=LR;\n")
	dot.WriteString("node [shape=plaintext];\n")
	dot.WriteString("label=\"Reporte BLOCK - Ana Lucía Nufio Roblero - 202400452\";\n")
	dot.WriteString("labelloc=t;\n\n")

	for i := int32(0); i < sb.SBlocksCount && int(i) < len(blockBitmap); i++ {
		if blockBitmap[i] != '1' {
			continue
		}

		kind := blockKinds[i]
		if kind == "" {
			kind = "file"
		}

		writeBlockGraphNode(&dot, diskPath, sb, i, kind)
	}

	writeBlockEdges(&dot, diskPath, sb, blockBitmap, blockKinds)

	dot.WriteString("}\n")

	return renderGraphviz(dot.String(), outPath)
}

func ReportTREE(diskPath string, partition disk.Partition, outPath string) error {
	sb, err := readEXT2SuperBlock(diskPath, partition)
	if err != nil {
		return err
	}

	inodeBitmap, err := ext2.ReadBitmap(diskPath, sb.SBmInodeStart, sb.SInodesCount)
	if err != nil {
		return fmt.Errorf("error al leer bitmap de inodos: %v", err)
	}

	blockBitmap, err := ext2.ReadBitmap(diskPath, sb.SBmBlockStart, sb.SBlocksCount)
	if err != nil {
		return fmt.Errorf("error al leer bitmap de bloques: %v", err)
	}

	var dot strings.Builder
	dot.WriteString("digraph G {\n")
	dot.WriteString("rankdir=LR;\n")
	dot.WriteString("node [shape=plaintext];\n")
	dot.WriteString("label=\"Reporte TREE - Ana Lucía Nufio Roblero - 202400452\";\n")
	dot.WriteString("labelloc=t;\n\n")

	ctx := &treeContext{
		diskPath:      diskPath,
		sb:            sb,
		inodeBitmap:   inodeBitmap,
		blockBitmap:   blockBitmap,
		dot:           &dot,
		writtenInodes: make(map[int32]bool),
		writtenBlocks: make(map[int32]bool),
		visitedInodes: make(map[int32]bool),
		visitedBlocks: make(map[int32]bool),
	}

	if !ctx.inodeIsUsed(0) {
		return fmt.Errorf("no existe el inodo raíz en la partición")
	}

	ctx.writeTreeInode(0)

	dot.WriteString("}\n")

	return renderGraphviz(dot.String(), outPath)
}

func writeInodeGraphNode(dot *strings.Builder, index int32, inode ext2.Inode) {
	dot.WriteString(fmt.Sprintf("inode_%d [label=<\n", index))
	dot.WriteString("<TABLE BORDER=\"0\" CELLBORDER=\"1\" CELLSPACING=\"0\" CELLPADDING=\"4\">\n")
	dot.WriteString(fmt.Sprintf("<TR><TD COLSPAN=\"2\" BGCOLOR=\"lightblue\"><B>Inodo %d</B></TD></TR>\n", index))
	dot.WriteString(fmt.Sprintf("<TR><TD>i_uid</TD><TD>%d</TD></TR>\n", inode.IUid))
	dot.WriteString(fmt.Sprintf("<TR><TD>i_gid</TD><TD>%d</TD></TR>\n", inode.IGid))
	dot.WriteString(fmt.Sprintf("<TR><TD>i_s</TD><TD>%d</TD></TR>\n", inode.IS))
	dot.WriteString(fmt.Sprintf("<TR><TD>i_atime</TD><TD>%s</TD></TR>\n", htmlEscape(reportTrimBytes(inode.IAtime[:]))))
	dot.WriteString(fmt.Sprintf("<TR><TD>i_ctime</TD><TD>%s</TD></TR>\n", htmlEscape(reportTrimBytes(inode.ICtime[:]))))
	dot.WriteString(fmt.Sprintf("<TR><TD>i_mtime</TD><TD>%s</TD></TR>\n", htmlEscape(reportTrimBytes(inode.IMtime[:]))))

	for i := 0; i < len(inode.IBlock); i++ {
		dot.WriteString(fmt.Sprintf("<TR><TD>i_block_%d</TD><TD>%d</TD></TR>\n", i+1, inode.IBlock[i]))
	}

	dot.WriteString(fmt.Sprintf("<TR><TD>i_type</TD><TD>%s</TD></TR>\n", htmlEscape(reportTrimBytes(inode.IType[:]))))
	dot.WriteString(fmt.Sprintf("<TR><TD>i_perm</TD><TD>%s</TD></TR>\n", htmlEscape(reportTrimBytes(inode.IPerm[:]))))
	dot.WriteString("</TABLE>>];\n\n")
}

func writeFolderInodeEdges(dot *strings.Builder, diskPath string, sb ext2.SuperBlock, inodeIndex int32, inode ext2.Inode, inodeBitmap []byte) {
	for i := 0; i < 12 && i < len(inode.IBlock); i++ {
		blockIndex := inode.IBlock[i]
		if blockIndex == -1 {
			continue
		}

		folderBlock, err := ext2.ReadFolderBlock(diskPath, sb, blockIndex)
		if err != nil {
			continue
		}

		for _, content := range folderBlock.BContent {
			name := reportTrimBytes(content.BName[:])

			if content.BInode == -1 || name == "" || name == "." || name == ".." {
				continue
			}

			if !reportInodeUsed(content.BInode, inodeBitmap, sb.SInodesCount) {
				continue
			}

			dot.WriteString(fmt.Sprintf("inode_%d -> inode_%d [label=\"%s\"];\n",
				inodeIndex, content.BInode, dotEscape(name)))
		}
	}
}

func writeBlockGraphNode(dot *strings.Builder, diskPath string, sb ext2.SuperBlock, blockIndex int32, kind string) {
	switch kind {
	case "folder":
		writeFolderBlockGraphNode(dot, diskPath, sb, blockIndex)

	case "pointer":
		writePointerBlockGraphNode(dot, diskPath, sb, blockIndex)

	default:
		writeFileBlockGraphNode(dot, diskPath, sb, blockIndex)
	}
}

func writeFolderBlockGraphNode(dot *strings.Builder, diskPath string, sb ext2.SuperBlock, blockIndex int32) {
	folderBlock, err := ext2.ReadFolderBlock(diskPath, sb, blockIndex)
	if err != nil {
		return
	}

	dot.WriteString(fmt.Sprintf("block_%d [label=<\n", blockIndex))
	dot.WriteString("<TABLE BORDER=\"0\" CELLBORDER=\"1\" CELLSPACING=\"0\" CELLPADDING=\"4\">\n")
	dot.WriteString(fmt.Sprintf("<TR><TD COLSPAN=\"2\" BGCOLOR=\"lightgoldenrod\"><B>Bloque Carpeta %d</B></TD></TR>\n", blockIndex))
	dot.WriteString("<TR><TD><B>b_name</B></TD><TD><B>b_inodo</B></TD></TR>\n")

	for _, content := range folderBlock.BContent {
		name := reportTrimBytes(content.BName[:])
		if name == "" {
			name = "-"
		}

		dot.WriteString(fmt.Sprintf("<TR><TD>%s</TD><TD>%d</TD></TR>\n", htmlEscape(name), content.BInode))
	}

	dot.WriteString("</TABLE>>];\n\n")
}

func writeFileBlockGraphNode(dot *strings.Builder, diskPath string, sb ext2.SuperBlock, blockIndex int32) {
	fileBlock, err := ext2.ReadFileBlock(diskPath, sb, blockIndex)
	if err != nil {
		return
	}

	content := reportTrimBytes(fileBlock.BContent[:])
	if len(content) > 80 {
		content = content[:80] + "..."
	}

	dot.WriteString(fmt.Sprintf("block_%d [label=<\n", blockIndex))
	dot.WriteString("<TABLE BORDER=\"0\" CELLBORDER=\"1\" CELLSPACING=\"0\" CELLPADDING=\"4\">\n")
	dot.WriteString(fmt.Sprintf("<TR><TD BGCOLOR=\"lightgreen\"><B>Bloque Archivo %d</B></TD></TR>\n", blockIndex))
	dot.WriteString(fmt.Sprintf("<TR><TD>%s</TD></TR>\n", htmlEscape(content)))
	dot.WriteString("</TABLE>>];\n\n")
}

func writePointerBlockGraphNode(dot *strings.Builder, diskPath string, sb ext2.SuperBlock, blockIndex int32) {
	pointerBlock, err := ext2.ReadPointerBlock(diskPath, sb, blockIndex)
	if err != nil {
		return
	}

	dot.WriteString(fmt.Sprintf("block_%d [label=<\n", blockIndex))
	dot.WriteString("<TABLE BORDER=\"0\" CELLBORDER=\"1\" CELLSPACING=\"0\" CELLPADDING=\"4\">\n")
	dot.WriteString(fmt.Sprintf("<TR><TD BGCOLOR=\"lightcoral\"><B>Bloque Apuntadores %d</B></TD></TR>\n", blockIndex))

	for i := 0; i < len(pointerBlock.BPointers); i += 4 {
		end := i + 4
		if end > len(pointerBlock.BPointers) {
			end = len(pointerBlock.BPointers)
		}

		values := make([]string, 0, end-i)
		for _, ptr := range pointerBlock.BPointers[i:end] {
			values = append(values, fmt.Sprintf("%d", ptr))
		}

		dot.WriteString(fmt.Sprintf("<TR><TD>%s</TD></TR>\n", htmlEscape(strings.Join(values, " | "))))
	}

	dot.WriteString("</TABLE>>];\n\n")
}

func writeBlockEdges(dot *strings.Builder, diskPath string, sb ext2.SuperBlock, blockBitmap []byte, blockKinds map[int32]string) {
	for blockIndex, kind := range blockKinds {
		if kind == "pointer" {
			pointerBlock, err := ext2.ReadPointerBlock(diskPath, sb, blockIndex)
			if err != nil {
				continue
			}

			for _, ptr := range pointerBlock.BPointers {
				if reportBlockUsed(ptr, blockBitmap, sb.SBlocksCount) {
					dot.WriteString(fmt.Sprintf("block_%d -> block_%d;\n", blockIndex, ptr))
				}
			}
		}
	}

	inodeBitmap, err := ext2.ReadBitmap(diskPath, sb.SBmInodeStart, sb.SInodesCount)
	if err != nil {
		return
	}

	for inodeIndex := int32(0); inodeIndex < sb.SInodesCount && int(inodeIndex) < len(inodeBitmap); inodeIndex++ {
		if inodeBitmap[inodeIndex] != '1' {
			continue
		}

		inode, err := ext2.ReadInode(diskPath, sb, inodeIndex)
		if err != nil || inode.IType[0] != '0' {
			continue
		}

		for i := 0; i < 12 && i < len(inode.IBlock); i++ {
			parentBlock := inode.IBlock[i]
			if !reportBlockUsed(parentBlock, blockBitmap, sb.SBlocksCount) {
				continue
			}

			folderBlock, err := ext2.ReadFolderBlock(diskPath, sb, parentBlock)
			if err != nil {
				continue
			}

			for _, entry := range folderBlock.BContent {
				name := reportTrimBytes(entry.BName[:])

				if entry.BInode == -1 || name == "" || name == "." || name == ".." {
					continue
				}

				if !reportInodeUsed(entry.BInode, inodeBitmap, sb.SInodesCount) {
					continue
				}

				childInode, err := ext2.ReadInode(diskPath, sb, entry.BInode)
				if err != nil {
					continue
				}

				for j := 0; j < 12 && j < len(childInode.IBlock); j++ {
					childBlock := childInode.IBlock[j]

					if reportBlockUsed(childBlock, blockBitmap, sb.SBlocksCount) {
						dot.WriteString(fmt.Sprintf("block_%d -> block_%d [label=\"%s\"];\n",
							parentBlock, childBlock, dotEscape(name)))
					}
				}
			}
		}
	}
}

func collectBlockKinds(diskPath string, sb ext2.SuperBlock, blockBitmap []byte) map[int32]string {
	kinds := make(map[int32]string)

	inodeBitmap, err := ext2.ReadBitmap(diskPath, sb.SBmInodeStart, sb.SInodesCount)
	if err != nil {
		return kinds
	}

	for inodeIndex := int32(0); inodeIndex < sb.SInodesCount && int(inodeIndex) < len(inodeBitmap); inodeIndex++ {
		if inodeBitmap[inodeIndex] != '1' {
			continue
		}

		inode, err := ext2.ReadInode(diskPath, sb, inodeIndex)
		if err != nil {
			continue
		}

		ownerKind := "file"
		if inode.IType[0] == '0' {
			ownerKind = "folder"
		}

		for i := 0; i < 12 && i < len(inode.IBlock); i++ {
			blockIndex := inode.IBlock[i]
			if reportBlockUsed(blockIndex, blockBitmap, sb.SBlocksCount) {
				kinds[blockIndex] = ownerKind
			}
		}

		if len(inode.IBlock) > 12 && reportBlockUsed(inode.IBlock[12], blockBitmap, sb.SBlocksCount) {
			markPointerBlockKinds(diskPath, sb, blockBitmap, kinds, inode.IBlock[12], ownerKind, 1)
		}

		if len(inode.IBlock) > 13 && reportBlockUsed(inode.IBlock[13], blockBitmap, sb.SBlocksCount) {
			markPointerBlockKinds(diskPath, sb, blockBitmap, kinds, inode.IBlock[13], ownerKind, 2)
		}

		if len(inode.IBlock) > 14 && reportBlockUsed(inode.IBlock[14], blockBitmap, sb.SBlocksCount) {
			markPointerBlockKinds(diskPath, sb, blockBitmap, kinds, inode.IBlock[14], ownerKind, 3)
		}
	}

	return kinds
}

func markPointerBlockKinds(diskPath string, sb ext2.SuperBlock, blockBitmap []byte, kinds map[int32]string, pointerBlockIndex int32, ownerKind string, level int) {
	if !reportBlockUsed(pointerBlockIndex, blockBitmap, sb.SBlocksCount) {
		return
	}

	kinds[pointerBlockIndex] = "pointer"

	pointerBlock, err := ext2.ReadPointerBlock(diskPath, sb, pointerBlockIndex)
	if err != nil {
		return
	}

	for _, ptr := range pointerBlock.BPointers {
		if !reportBlockUsed(ptr, blockBitmap, sb.SBlocksCount) {
			continue
		}

		if level <= 1 {
			kinds[ptr] = ownerKind
		} else {
			markPointerBlockKinds(diskPath, sb, blockBitmap, kinds, ptr, ownerKind, level-1)
		}
	}
}

type treeContext struct {
	diskPath      string
	sb            ext2.SuperBlock
	inodeBitmap   []byte
	blockBitmap   []byte
	dot           *strings.Builder
	writtenInodes map[int32]bool
	writtenBlocks map[int32]bool
	visitedInodes map[int32]bool
	visitedBlocks map[int32]bool
}

func (ctx *treeContext) inodeIsUsed(index int32) bool {
	return reportInodeUsed(index, ctx.inodeBitmap, ctx.sb.SInodesCount)
}

func (ctx *treeContext) blockIsUsed(index int32) bool {
	return reportBlockUsed(index, ctx.blockBitmap, ctx.sb.SBlocksCount)
}

func (ctx *treeContext) writeTreeInode(inodeIndex int32) {
	if !ctx.inodeIsUsed(inodeIndex) {
		return
	}

	inode, err := ext2.ReadInode(ctx.diskPath, ctx.sb, inodeIndex)
	if err != nil {
		return
	}

	if !ctx.writtenInodes[inodeIndex] {
		writeInodeGraphNode(ctx.dot, inodeIndex, inode)
		ctx.writtenInodes[inodeIndex] = true
	}

	if ctx.visitedInodes[inodeIndex] {
		return
	}

	ctx.visitedInodes[inodeIndex] = true

	ownerKind := "file"
	if inode.IType[0] == '0' {
		ownerKind = "folder"
	}

	for i := 0; i < 12 && i < len(inode.IBlock); i++ {
		blockIndex := inode.IBlock[i]
		if !ctx.blockIsUsed(blockIndex) {
			continue
		}

		ctx.dot.WriteString(fmt.Sprintf("inode_%d -> block_%d [label=\"i_block_%d\"];\n", inodeIndex, blockIndex, i+1))
		ctx.writeTreeDataBlock(blockIndex, inodeIndex, ownerKind)
	}

	if len(inode.IBlock) > 12 && ctx.blockIsUsed(inode.IBlock[12]) {
		ctx.dot.WriteString(fmt.Sprintf("inode_%d -> block_%d [label=\"simple\"];\n", inodeIndex, inode.IBlock[12]))
		ctx.writeTreePointerBlock(inode.IBlock[12], inodeIndex, ownerKind, 1)
	}

	if len(inode.IBlock) > 13 && ctx.blockIsUsed(inode.IBlock[13]) {
		ctx.dot.WriteString(fmt.Sprintf("inode_%d -> block_%d [label=\"doble\"];\n", inodeIndex, inode.IBlock[13]))
		ctx.writeTreePointerBlock(inode.IBlock[13], inodeIndex, ownerKind, 2)
	}

	if len(inode.IBlock) > 14 && ctx.blockIsUsed(inode.IBlock[14]) {
		ctx.dot.WriteString(fmt.Sprintf("inode_%d -> block_%d [label=\"triple\"];\n", inodeIndex, inode.IBlock[14]))
		ctx.writeTreePointerBlock(inode.IBlock[14], inodeIndex, ownerKind, 3)
	}
}

func (ctx *treeContext) writeTreeDataBlock(blockIndex int32, ownerInodeIndex int32, ownerKind string) {
	if !ctx.blockIsUsed(blockIndex) {
		return
	}

	if !ctx.writtenBlocks[blockIndex] {
		if ownerKind == "folder" {
			writeFolderBlockGraphNode(ctx.dot, ctx.diskPath, ctx.sb, blockIndex)
		} else {
			writeFileBlockGraphNode(ctx.dot, ctx.diskPath, ctx.sb, blockIndex)
		}

		ctx.writtenBlocks[blockIndex] = true
	}

	if ctx.visitedBlocks[blockIndex] {
		return
	}

	ctx.visitedBlocks[blockIndex] = true

	if ownerKind != "folder" {
		return
	}

	folderBlock, err := ext2.ReadFolderBlock(ctx.diskPath, ctx.sb, blockIndex)
	if err != nil {
		return
	}

	for _, entry := range folderBlock.BContent {
		name := reportTrimBytes(entry.BName[:])

		if entry.BInode == -1 || name == "" || name == "." || name == ".." || entry.BInode == ownerInodeIndex {
			continue
		}

		if !ctx.inodeIsUsed(entry.BInode) {
			continue
		}

		ctx.dot.WriteString(fmt.Sprintf("block_%d -> inode_%d [label=\"%s\"];\n",
			blockIndex, entry.BInode, dotEscape(name)))

		ctx.writeTreeInode(entry.BInode)
	}
}

func (ctx *treeContext) writeTreePointerBlock(blockIndex int32, ownerInodeIndex int32, ownerKind string, level int) {
	if !ctx.blockIsUsed(blockIndex) {
		return
	}

	if !ctx.writtenBlocks[blockIndex] {
		writePointerBlockGraphNode(ctx.dot, ctx.diskPath, ctx.sb, blockIndex)
		ctx.writtenBlocks[blockIndex] = true
	}

	if ctx.visitedBlocks[blockIndex] {
		return
	}

	ctx.visitedBlocks[blockIndex] = true

	pointerBlock, err := ext2.ReadPointerBlock(ctx.diskPath, ctx.sb, blockIndex)
	if err != nil {
		return
	}

	for _, ptr := range pointerBlock.BPointers {
		if !ctx.blockIsUsed(ptr) {
			continue
		}

		ctx.dot.WriteString(fmt.Sprintf("block_%d -> block_%d;\n", blockIndex, ptr))

		if level <= 1 {
			ctx.writeTreeDataBlock(ptr, ownerInodeIndex, ownerKind)
		} else {
			ctx.writeTreePointerBlock(ptr, ownerInodeIndex, ownerKind, level-1)
		}
	}
}

func reportInodeUsed(index int32, bitmap []byte, total int32) bool {
	return index >= 0 && index < total && int(index) < len(bitmap) && bitmap[index] == '1'
}

func reportBlockUsed(index int32, bitmap []byte, total int32) bool {
	return index >= 0 && index < total && int(index) < len(bitmap) && bitmap[index] == '1'
}

func reportTrimBytes(data []byte) string {
	return strings.TrimRight(string(data), "\x00")
}

func dotEscape(value string) string {
	value = strings.ReplaceAll(value, "\\", "\\\\")
	value = strings.ReplaceAll(value, "\"", "\\\"")
	value = strings.ReplaceAll(value, "\n", "\\n")
	value = strings.ReplaceAll(value, "\r", "")
	return value
}
