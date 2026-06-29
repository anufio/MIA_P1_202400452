package services

//partition_space_service.go: contiene la lógica para manejar el espacio libre dentro de las particiones del sistema de archivos, incluyendo la identificación de espacios libres en particiones extendidas y la selección de espacios según criterios de ajuste (first fit, best fit, worst fit).

import (
	"MIA_P2_202400452/internal/disk"
)

type logicalSpace struct {
	start int32
	size  int32
}

func (s *DiskService) logicalFreeSpaces(extended disk.Partition, records []logicalRecord) []logicalSpace {
	var spaces []logicalSpace

	lastEnd := extended.PartStart
	extendedEnd := extended.PartStart + extended.PartS

	for _, record := range records {
		if record.position > lastEnd {
			spaces = append(spaces, logicalSpace{
				start: lastEnd,
				size:  record.position - lastEnd,
			})
		}

		recordEnd := record.ebr.PartStart + record.ebr.PartS
		if recordEnd > lastEnd {
			lastEnd = recordEnd
		}
	}

	if lastEnd < extendedEnd {
		spaces = append(spaces, logicalSpace{
			start: lastEnd,
			size:  extendedEnd - lastEnd,
		})
	}

	return spaces
}

func (s *DiskService) selectLogicalSpace(spaces []logicalSpace, requiredSize int32, fitByte byte) logicalSpace {
	selected := logicalSpace{start: -1}

	for _, space := range spaces {
		if space.size < requiredSize {
			continue
		}

		switch fitByte {
		case 'F':
			return space

		case 'B':
			if selected.start == -1 || space.size < selected.size {
				selected = space
			}

		case 'W':
			if selected.start == -1 || space.size > selected.size {
				selected = space
			}

		default:
			if selected.start == -1 {
				selected = space
			}
		}
	}

	return selected
}
