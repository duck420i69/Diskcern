package scanner

import (
	"bytes"
	"encoding/csv"
	"io"
	"os"
	"strconv"
	"strings"
	"unicode/utf16"
	"unicode/utf8"

	"github.com/diskcern/diskcern/internal/models"
)

func ImportWizTreeCSV(csvPath string) ([]models.FileRecord, error) {
	data, err := os.ReadFile(csvPath)
	if err != nil {
		return nil, err
	}

	if len(data) >= 2 && data[0] == 0xff && data[1] == 0xfe {
		u16s := make([]uint16, 0, len(data)/2)
		for i := 2; i < len(data); i += 2 {
			if i+1 >= len(data) { break }
			u16s = append(u16s, uint16(data[i]) | (uint16(data[i+1])<<8))
		}
		runes := utf16.Decode(u16s)
		var buf bytes.Buffer
		b := make([]byte, 4)
		for _, r := range runes {
			n := utf8.EncodeRune(b, r)
			buf.Write(b[:n])
		}
		data = buf.Bytes()
	}

	reader := csv.NewReader(bytes.NewReader(data))
	reader.LazyQuotes = true
	reader.FieldsPerRecord = -1
	
	header, err := reader.Read()
	if err != nil {
		return nil, err
	}

	pathIdx, sizeIdx, attrIdx := 0, 1, 4
	for i, h := range header {
		h = strings.ToLower(strings.TrimSpace(h))
		if strings.Contains(h, "file") { pathIdx = i }
		if strings.Contains(h, "size") { sizeIdx = i }
		if strings.Contains(h, "attr") { attrIdx = i }
	}

	var records []models.FileRecord
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil || len(record) <= sizeIdx || len(record) <= pathIdx {
			continue
		}

		path := record[pathIdx]
		sizeStr := strings.ReplaceAll(record[sizeIdx], ",", "")
		size, _ := strconv.ParseInt(sizeStr, 10, 64)
		
		isDir := true // Default to true for WizTree paths if attributes fail
		if len(record) > attrIdx {
			attr := record[attrIdx]
			if attr != "" {
				isDir = strings.Contains(strings.ToLower(attr), "d") || strings.Contains(strings.ToLower(attr), "dir")
			}
		} else {
			lastSlash := strings.LastIndex(path, "\\")
			if lastSlash != -1 && strings.Contains(path[lastSlash:], ".") {
				isDir = false // Has extension, probably a file
			}
		}

		records = append(records, models.FileRecord{
			Path:  path,
			Size:  size,
			IsDir: isDir,
		})
	}
	return records, nil
}
