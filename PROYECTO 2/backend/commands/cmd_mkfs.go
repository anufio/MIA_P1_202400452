package commands

// cmd_mkfs.go: comando para formatear una partición montada con el sistema de archivos EXT2

import (
	"MIA_P1_202400452/ext2"
	"fmt"
	"strings"
)

func CmdMKFS(params map[string]string) string {
	id, ok := params["id"]
	if !ok {
		return "Error: falta el parámetro obligatorio -id"
	}

	fmtType := "full"
	if t, ok := params["type"]; ok {
		fmtType = strings.ToLower(t)
		if fmtType != "full" {
			return fmt.Sprintf("Error: tipo de formateo inválido '%s'. Use 'full'", params["type"])
		}
	}

	diskPath, partition, _, err := GetMountedPartitionInfo(id)
	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}
	partStart := partition.PartStart
	partSize := partition.PartS

	numInodes, numBlocks := ext2.CalculateStructures(partSize)
	if numInodes < 2 || numBlocks < 2 {
		return "Error: la partición es demasiado pequeña para formatear"
	}

	sbSize := int32(ext2.SizeSuperBlock())
	inodeSize := int32(ext2.SizeInode())
	blockSize := int32(ext2.SizeBlock())

	inodeBitmapSize := numInodes
	blockBitmapSize := numBlocks
	inodeTableSize := numInodes * inodeSize
	blockTableSize := numBlocks * blockSize

	sbStart := partStart
	bmInodeStart := sbStart + sbSize
	bmBlockStart := bmInodeStart + inodeBitmapSize
	inodeStart := bmBlockStart + blockBitmapSize
	blockStart := inodeStart + inodeTableSize
	if blockStart+blockTableSize > partStart+partSize {
		return "Error: la partición no tiene espacio suficiente para las estructuras EXT2"
	}

	sb := ext2.SuperBlock{
		SFilesystemType:  2,
		SInodesCount:     numInodes,
		SBlocksCount:     numBlocks,
		SFreeInodesCount: numInodes,
		SFreeBlocksCount: numBlocks,
		SMagic:           0xEF53,
		SInodeS:          inodeSize,
		SBlockS:          blockSize,
		SBmInodeStart:    bmInodeStart,
		SBmBlockStart:    bmBlockStart,
		SInodeStart:      inodeStart,
		SBlockStart:      blockStart,
	}
	sb.SMtime = ext2.GetCurrentTime()
	sb.SUmtime = ext2.GetCurrentTime()
	sb.SMntCount = 0

	if err := ext2.WriteSuperBlock(diskPath, sbStart, sb); err != nil {
		return fmt.Sprintf("Error escribiendo superbloque: %v", err)
	}

	if err := ext2.InitBitmap(diskPath, bmInodeStart, inodeBitmapSize); err != nil {
		return fmt.Sprintf("Error inicializando bitmap de inodos: %v", err)
	}
	if err := ext2.InitBitmap(diskPath, bmBlockStart, blockBitmapSize); err != nil {
		return fmt.Sprintf("Error inicializando bitmap de bloques: %v", err)
	}
	ext2.UpdateFirstFreeInode(diskPath, &sb)
	ext2.UpdateFirstFreeBlock(diskPath, &sb)

	rootInodeIdx, err := ext2.AllocateInode(diskPath, &sb)
	if err != nil {
		return fmt.Sprintf("Error asignando inodo raíz: %v", err)
	}
	if rootInodeIdx != 0 {
		return fmt.Sprintf("Error: el inodo raíz debe ser 0 y se obtuvo %d", rootInodeIdx)
	}
	rootInode := ext2.NewInode()
	rootInode.IType = [1]byte{'0'}
	rootInode.IPerm = [3]byte{'6', '6', '4'}
	rootInode.IUid = 1
	rootInode.IGid = 1
	rootInode.ICtime = ext2.GetCurrentTime()
	rootInode.IAtime = ext2.GetCurrentTime()
	rootInode.IMtime = ext2.GetCurrentTime()

	rootBlockIdx, err := ext2.AllocateBlock(diskPath, &sb)
	if err != nil {
		return fmt.Sprintf("Error asignando bloque para raíz: %v", err)
	}

	rootFolder := ext2.NewFolderBlock()
	copy(rootFolder.BContent[0].BName[:], ".")
	rootFolder.BContent[0].BInode = 0
	copy(rootFolder.BContent[1].BName[:], "..")
	rootFolder.BContent[1].BInode = 0
	if err := ext2.WriteFolderBlock(diskPath, sb, rootBlockIdx, rootFolder); err != nil {
		return fmt.Sprintf("Error escribiendo bloque raíz: %v", err)
	}
	rootInode.IBlock[0] = rootBlockIdx
	if err := ext2.WriteInode(diskPath, sb, 0, rootInode); err != nil {
		return fmt.Sprintf("Error escribiendo inodo raíz: %v", err)
	}

	usersInodeIdx, err := ext2.AllocateInode(diskPath, &sb)
	if err != nil {
		return fmt.Sprintf("Error asignando inodo users.txt: %v", err)
	}
	if usersInodeIdx != 1 {
		return fmt.Sprintf("Error: el inodo de users.txt debe ser 1 y se obtuvo %d", usersInodeIdx)
	}
	usersInode := ext2.NewInode()
	usersInode.IType = [1]byte{'1'}
	usersInode.IPerm = [3]byte{'6', '6', '4'}
	usersInode.IUid = 1
	usersInode.IGid = 1
	usersInode.ICtime = ext2.GetCurrentTime()
	usersInode.IAtime = ext2.GetCurrentTime()
	usersInode.IMtime = ext2.GetCurrentTime()

	content := "1,G,root\n1,U,root,root,123\n"
	if err := ext2.WriteFileContent(diskPath, &sb, &usersInode, []byte(content)); err != nil {
		return fmt.Sprintf("Error escribiendo contenido de users.txt: %v", err)
	}
	if err := ext2.WriteInode(diskPath, sb, 1, usersInode); err != nil {
		return fmt.Sprintf("Error escribiendo inodo de users.txt: %v", err)
	}

	if err := ext2.AddEntryToFolder(diskPath, &sb, 0, "users.txt", 1); err != nil {
		return fmt.Sprintf("Error agregando users.txt a la raíz: %v", err)
	}

	ext2.UpdateFirstFreeInode(diskPath, &sb)
	ext2.UpdateFirstFreeBlock(diskPath, &sb)
	if err := ext2.WriteSuperBlock(diskPath, sbStart, sb); err != nil {
		return fmt.Sprintf("Error actualizando superbloque: %v", err)
	}

	return fmt.Sprintf("Partición '%s' formateada como EXT2 exitosamente.\n Inodos: %d | Bloques: %d", id, numInodes, numBlocks)
}
