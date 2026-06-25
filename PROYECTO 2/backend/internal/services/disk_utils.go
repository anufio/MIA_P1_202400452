package services

//disk_utils.go: contiene funciones auxiliares para manejar discos, incluyendo la codificación y decodificación de IDs de disco, validación de extensiones de disco, conversión de tamaños y manejo de nombres y IDs.

import (
	"encoding/base64"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"MIA_P2_202400452/internal/disk"
	"MIA_P2_202400452/internal/ext2"
)

func EncodeDiskID(path string) string {
	return base64.RawURLEncoding.EncodeToString([]byte(path))
}

func DecodeDiskID(id string) (string, error) {
	data, err := base64.RawURLEncoding.DecodeString(id)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" {
			return value
		}
	}

	return ""
}

func validDiskExtension(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".dsk" || ext == ".mia"
}

func cleanName(raw []byte) string {
	return strings.TrimSpace(strings.TrimRight(string(raw), "\x00"))
}

func equalName(a string, b string) bool {
	return strings.EqualFold(strings.TrimSpace(a), strings.TrimSpace(b))
}

func copyName16(name string) [16]byte {
	var arr [16]byte
	copy(arr[:], name)
	return arr
}

func copyID4(id string) [4]byte {
	var arr [4]byte
	copy(arr[:], strings.ToUpper(id))
	return arr
}

func byteToFit(value byte) string {
	switch value {
	case 'B':
		return "BF"
	case 'F':
		return "FF"
	case 'W':
		return "WF"
	default:
		return "FF"
	}
}

func positiveSizeToBytes(size int64, unit string) (int32, error) {
	if size <= 0 {
		return 0, fmt.Errorf("el tamaño debe ser mayor que cero")
	}

	unit = strings.ToUpper(strings.TrimSpace(unit))
	if unit == "" {
		unit = "M"
	}

	var bytes int64

	switch unit {
	case "B":
		bytes = size
	case "K":
		bytes = size * 1024
	case "M":
		bytes = size * 1024 * 1024
	default:
		return 0, fmt.Errorf("unidad inválida '%s'. Use B, K o M", unit)
	}

	if bytes > 2147483647 {
		return 0, fmt.Errorf("el tamaño excede el máximo permitido")
	}

	return int32(bytes), nil
}

func signedSizeToBytes(size int64, unit string) (int32, error) {
	if size == 0 {
		return 0, fmt.Errorf("el valor add debe ser diferente de cero")
	}

	sign := int64(1)

	if size < 0 {
		sign = -1
		size = -size
	}

	result, err := positiveSizeToBytes(size, unit)
	if err != nil {
		return 0, err
	}

	return int32(sign) * result, nil
}

func parseInt64(value string) int64 {
	value = strings.TrimSpace(value)

	if value == "" {
		return 0
	}

	result, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0
	}

	return result
}

func isFormatted(path string, start int32) bool {
	sb, err := ext2.ReadSuperBlock(path, start)
	if err != nil {
		return false
	}

	return sb.SMagic == 0xEF53
}

func zeroRange(path string, start int64, size int64) error {
	if size <= 0 {
		return nil
	}

	const chunkSize = 1024 * 1024
	zeroes := make([]byte, chunkSize)

	written := int64(0)

	for written < size {
		toWrite := chunkSize

		if written+int64(toWrite) > size {
			toWrite = int(size - written)
		}

		if err := disk.WriteBytes(path, start+written, zeroes[:toWrite]); err != nil {
			return err
		}

		written += int64(toWrite)
	}

	return nil
}
