package services

// report_service.go: genera reportes del sistema y los deja listos para verse desde la GUI.

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"MIA_P2_202400452/internal/disk"
	"MIA_P2_202400452/internal/reports"
)

type ReportService struct {
	reportRoot  string
	diskService *DiskService
}

func NewReportService(reportRoot string, diskService *DiskService) *ReportService {
	_ = os.MkdirAll(reportRoot, 0755)

	return &ReportService{
		reportRoot:  reportRoot,
		diskService: diskService,
	}
}

func (s *ReportService) Generate(input ReportInput) (ReportOutput, error) {
	reportType := reports.NormalizeReportType(strings.TrimSpace(input.Type))
	id := strings.ToUpper(strings.TrimSpace(input.ID))

	if reportType == "" {
		reportType = reports.NormalizeReportType(input.Name)
	}

	if reportType == "" {
		return ReportOutput{}, fmt.Errorf("debe indicar el tipo de reporte")
	}

	rawOutputPath := strings.TrimSpace(input.Path)
	inputDiskPath := strings.TrimSpace(input.DiskPath)
	if inputDiskPath == "" {
		inputDiskPath = strings.TrimSpace(input.FilePath)
	}

	// Compatibilidad: si un frontend viejo mandó la ruta del disco en path,
	// no se debe usar como archivo de salida.
	if inputDiskPath == "" && strings.EqualFold(filepath.Ext(rawOutputPath), ".dsk") {
		inputDiskPath = rawOutputPath
		rawOutputPath = ""
	}

	outputPath := rawOutputPath
	if outputPath == "" {
		outputPath = filepath.Join(
			s.reportRoot,
			fmt.Sprintf("%s_%d%s", reportType, time.Now().UnixNano(), defaultReportExtension(reportType)),
		)
	}

	outputPath = absolutePath(outputPath)

	var diskPath string
	var partition disk.Partition

	switch reportType {
	case "mbr", "disk":
		if id != "" {
			mounted, ok := s.diskService.GetMounted(id)
			if !ok {
				return ReportOutput{}, fmt.Errorf("la partición '%s' no está montada", id)
			}

			diskPath = mounted.DiskPath
		} else {
			diskPath = inputDiskPath
		}

	default:
		if id == "" {
			return ReportOutput{}, fmt.Errorf("debe indicar el id de la partición montada")
		}

		mounted, ok := s.diskService.GetMounted(id)
		if !ok {
			return ReportOutput{}, fmt.Errorf("la partición '%s' no está montada", id)
		}

		diskPath = mounted.DiskPath
		partition = mountedToPartition(mounted)
	}

	if diskPath == "" {
		return ReportOutput{}, fmt.Errorf("debe indicar ruta del disco o id de partición")
	}

	result, err := reports.GenerateReport(reports.ReportRequest{
		Type:      reportType,
		DiskPath:  diskPath,
		Path:      outputPath,
		FilePath:  input.FilePath,
		Partition: partition,
	})
	if err != nil {
		return ReportOutput{}, err
	}

	return ReportOutput{
		Type:   result.Type,
		Path:   result.Path,
		URL:    "/reports/" + filepath.Base(result.Path),
		Format: result.Format,
	}, nil
}

func defaultReportExtension(reportType string) string {
	switch reportType {
	case "tree", "disk", "mbr", "inode", "block", "bm_inode", "bm_block", "sb":
		return ".svg"
	case "file", "ls":
		return ".txt"
	default:
		return ".svg"
	}
}

func mountedToPartition(mounted MountedInfo) disk.Partition {
	partType := byte('P')

	if strings.ToUpper(mounted.Type) == "L" {
		partType = 'L'
	}

	partition := disk.Partition{
		PartStatus: [1]byte{'1'},
		PartType:   [1]byte{partType},
		PartFit:    [1]byte{'W'},
		PartStart:  mounted.Start,
		PartS:      mounted.Size,
		PartName:   copyName16(mounted.Name),
		PartId:     copyID4(mounted.ID),
	}

	return partition
}
