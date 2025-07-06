package helper

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

type StreamHelper struct {
	File        *excelize.File
	Writer      *excelize.StreamWriter
	Sheet       string
	StyleConfig StyleConfig
}

// HeaderConfig defines header configuration
type HeaderConfig struct {
	Cell   string
	Values []interface{}
	Style  *excelize.Style
}

// ColumnConfig defines column configuration
type ColumnConfig struct {
	Column int
	Width  float64
	Style  *excelize.Style
}

// ColumnStyleConfig defines column style configuration
type ColumnStyleConfig struct {
	ColumnRange
	Style string // Style name from StyleConfig
}

// ColumnRange defines a range of columns
type ColumnRange struct {
	Start int // Start column index (0-based)
	End   int // End column index (0-based)
}

// MergeConfig defines merge configuration
type MergeConfig struct {
	StartCell string
	EndCell   string
}

// MergeRange defines a simple merge range
type MergeRange struct {
	StartCell string
	EndCell   string
}

type StyleConfig map[string]int

func NewStreamHelper(file *excelize.File, sheet string) (*StreamHelper, error) {
	writer, err := file.NewStreamWriter(sheet)
	if err != nil {
		return nil, fmt.Errorf("failed to create stream writer: %w", err)
	}

	return &StreamHelper{
		File:   file,
		Writer: writer,
		Sheet:  sheet,
	}, nil
}

func (sh *StreamHelper) WriteHeader(config HeaderConfig) error {
	if err := sh.Writer.SetRow(config.Cell, config.Values); err != nil {
		return fmt.Errorf("failed to set header row: %w", err)
	}

	return nil
}

func (sh *StreamHelper) MergeCells(ranges []MergeRange) error {
	mustMerge := func(start, end string) error {
		if err := sh.Writer.MergeCell(start, end); err != nil {
			return fmt.Errorf("failed to merge %s:%s: %v", start, end, err)
		}
		return nil
	}

	for _, r := range ranges {
		if err := mustMerge(r.StartCell, r.EndCell); err != nil {
			return err
		}
	}
	return nil
}

func (sh *StreamHelper) SetColumnStyles(configs []ColumnStyleConfig, styleConfig map[string]int) error {
	mustSet := func(colStart, colEnd, style int) error {
		if err := sh.Writer.SetColStyle(colStart, colEnd, style); err != nil {
			return fmt.Errorf("failed to set column style %d:%d to style %d: %v", colStart, colEnd, style, err)
		}
		return nil
	}

	for _, config := range configs {
		styleID, exists := styleConfig[config.Style]
		if !exists {
			return fmt.Errorf("style ID for column %d not found in style config", config.Style)
		}
		if err := mustSet(config.Start, config.End, styleID); err != nil {
			return err
		}
	}

	return nil
}

func (sh *StreamHelper) SetColumnWidths(widths map[ColumnRange]float64) error {
	mustSet := func(colStart, colEnd int, width float64) error {
		if err := sh.Writer.SetColWidth(colStart, colEnd, width); err != nil {
			return fmt.Errorf("failed to set column width %d:%d to %.2f: %v", colStart, colEnd, width, err)
		}
		return nil
	}

	for colRange, width := range widths {
		if err := mustSet(colRange.Start, colRange.End, width); err != nil {
			return err
		}
	}
	return nil
}

// Flush commits the stream writer
func (sh *StreamHelper) Flush() error {
	if err := sh.Writer.Flush(); err != nil {
		return fmt.Errorf("failed to flush stream writer: %w", err)
	}
	return nil
}

// Close closes the stream writer
func (sh *StreamHelper) Close() error {
	if err := sh.Writer.Flush(); err != nil {
		return fmt.Errorf("failed to flush before close: %w", err)
	}
	return nil
}

func CreateStyles(f *excelize.File, styleDefs map[string]*excelize.Style) (map[string]int, error) {
	styles := make(map[string]int)
	for key, def := range styleDefs {
		styleID, err := f.NewStyle(def)
		if err != nil {
			return nil, err
		}
		styles[key] = styleID
	}
	return styles, nil
}
