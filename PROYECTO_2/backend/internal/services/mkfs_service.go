package services

//mkfs_service.go: contiene la lógica para manejar el formateo de particiones dentro del sistema de archivos, incluyendo la creación de estructuras EXT2, la inicialización de bitmaps y la creación de archivos y carpetas raíz.

import (
	"fmt"
	"strings"

	"MIA_P2_202400452/internal/ext2"
)

func (s *DiskService) FormatPartition(input MkfsInput) (PartitionInfo, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	id := strings.ToUpper(strings.TrimSpace(input.ID))
	formatType := strings.ToLower(strings.TrimSpace(input.Type))

	if id == "" {
		return PartitionInfo{}, fmt.Errorf("debe indicar el id de la partición montada")
	}

	if formatType == "" {
		formatType = "full"
	}

	if formatType != "full" {
		return PartitionInfo{}, fmt.Errorf("tipo de formateo inválido. Use full")
	}

	mounted, ok := s.mounted[id]
	if !ok {
		return PartitionInfo{}, fmt.Errorf("el id '%s' no está montado", id)
	}

	partStart := mounted.Start
	partSize := mounted.Size

	numInodes, numBlocks := ext2.CalculateStructures(partSize)
	if numInodes < 2 || numBlocks < 2 {
		return PartitionInfo{}, fmt.Errorf("la partición es demasiado pequeña para formatear")
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
		return PartitionInfo{}, fmt.Errorf("la partición no tiene espacio suficiente para las estructuras EXT2")
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

	if err := ext2.WriteSuperBlock(mounted.DiskPath, sbStart, sb); err != nil {
		return PartitionInfo{}, fmt.Errorf("error escribiendo superbloque: %v", err)
	}

	if err := ext2.InitBitmap(mounted.DiskPath, bmInodeStart, inodeBitmapSize); err != nil {
		return PartitionInfo{}, fmt.Errorf("error inicializando bitmap de inodos: %v", err)
	}

	if err := ext2.InitBitmap(mounted.DiskPath, bmBlockStart, blockBitmapSize); err != nil {
		return PartitionInfo{}, fmt.Errorf("error inicializando bitmap de bloques: %v", err)
	}

	ext2.UpdateFirstFreeInode(mounted.DiskPath, &sb)
	ext2.UpdateFirstFreeBlock(mounted.DiskPath, &sb)

	rootInodeIdx, err := ext2.AllocateInode(mounted.DiskPath, &sb)
	if err != nil {
		return PartitionInfo{}, fmt.Errorf("error asignando inodo raíz: %v", err)
	}

	if rootInodeIdx != 0 {
		return PartitionInfo{}, fmt.Errorf("el inodo raíz debe ser 0 y se obtuvo %d", rootInodeIdx)
	}

	rootInode := ext2.NewInode()
	rootInode.IType = [1]byte{'0'}
	rootInode.IPerm = [3]byte{'6', '6', '4'}
	rootInode.IUid = 1
	rootInode.IGid = 1
	rootInode.ICtime = ext2.GetCurrentTime()
	rootInode.IAtime = ext2.GetCurrentTime()
	rootInode.IMtime = ext2.GetCurrentTime()

	rootBlockIdx, err := ext2.AllocateBlock(mounted.DiskPath, &sb)
	if err != nil {
		return PartitionInfo{}, fmt.Errorf("error asignando bloque raíz: %v", err)
	}

	rootFolder := ext2.NewFolderBlock()
	copy(rootFolder.BContent[0].BName[:], ".")
	rootFolder.BContent[0].BInode = 0
	copy(rootFolder.BContent[1].BName[:], ".")
	rootFolder.BContent[1].BInode = 0

	if err := ext2.WriteFolderBlock(mounted.DiskPath, sb, rootBlockIdx, rootFolder); err != nil {
		return PartitionInfo{}, fmt.Errorf("error escribiendo bloque raíz: %v", err)
	}

	rootInode.IBlock[0] = rootBlockIdx

	if err := ext2.WriteInode(mounted.DiskPath, sb, 0, rootInode); err != nil {
		return PartitionInfo{}, fmt.Errorf("error escribiendo inodo raíz: %v", err)
	}

	usersInodeIdx, err := ext2.AllocateInode(mounted.DiskPath, &sb)
	if err != nil {
		return PartitionInfo{}, fmt.Errorf("error asignando inodo users.txt: %v", err)
	}

	if usersInodeIdx != 1 {
		return PartitionInfo{}, fmt.Errorf("el inodo de users.txt debe ser 1 y se obtuvo %d", usersInodeIdx)
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

	if err := ext2.WriteFileContent(mounted.DiskPath, &sb, &usersInode, []byte(content)); err != nil {
		return PartitionInfo{}, fmt.Errorf("error escribiendo users.txt: %v", err)
	}

	if err := ext2.WriteInode(mounted.DiskPath, sb, 1, usersInode); err != nil {
		return PartitionInfo{}, fmt.Errorf("error escribiendo inodo users.txt: %v", err)
	}

	if err := ext2.AddEntryToFolder(mounted.DiskPath, &sb, 0, "users.txt", 1); err != nil {
		return PartitionInfo{}, fmt.Errorf("error agregando users.txt a raíz: %v", err)
	}

	ext2.UpdateFirstFreeInode(mounted.DiskPath, &sb)
	ext2.UpdateFirstFreeBlock(mounted.DiskPath, &sb)

	if err := ext2.WriteSuperBlock(mounted.DiskPath, partStart, sb); err != nil {
		return PartitionInfo{}, fmt.Errorf("error actualizando superbloque: %v", err)
	}

	return PartitionInfo{
		ID:        id,
		DiskID:    EncodeDiskID(mounted.DiskPath),
		DiskPath:  mounted.DiskPath,
		Name:      mounted.Name,
		Type:      mounted.Type,
		Start:     mounted.Start,
		Size:      mounted.Size,
		Mounted:   true,
		Formatted: true,
	}, nil
}
