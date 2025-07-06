package main

import (
	"excel/stream/helper"
	"fmt"
	"log"
	"time"

	"github.com/xuri/excelize/v2"
)

func main() {
	fmt.Println("Starting Excel stream example...")

	startTime := time.Now()

	fmt.Printf("Start time: %s\n", startTime.Format("2006-01-02 15:04:05"))
	testHelper()

	endTime := time.Since(startTime)
	fmt.Printf("End time: %s\n", endTime)
}

func testHelper() {
	f := excelize.NewFile()

	sw, err := helper.NewStreamHelper(f, "Sheet1")
	if err != nil {
		log.Fatalf("failed to create stream helper: %v", err)
	}

	setFormat(f, sw)

	for i := 3; i <= 1000002; i++ {
		row := []interface{}{
			i - 2,
			fmt.Sprintf("Name%d", i-2),
			float64(i) * 1.5,
			fmt.Sprintf("2025-07-%02d", (i-2)%30+1),
			fmt.Sprintf("%02d:00", (i-2)%24),
			fmt.Sprintf("2025-07-%02d", (i-2)%30+1),
			fmt.Sprintf("%02d:30", (i-2)%24),
			fmt.Sprintf("2025-07-%02d", (i-2)%30+1),
			fmt.Sprintf("%02d:45", (i-2)%24),
			i - 2,
			fmt.Sprintf("ACC%04d", i-2),
		}
		cell, _ := excelize.CoordinatesToCellName(1, i)
		if err := sw.Writer.SetRow(cell, row, excelize.RowOpts{
			Height: 20,
		}); err != nil {
			log.Fatalf("failed to write row %d: %v", i, err)
		}
	}

	if err := sw.Flush(); err != nil {
		log.Fatalf("failed to flush stream writer: %v", err)
	}

	if err := f.SaveAs("stream_example.xlsx"); err != nil {
		log.Fatalf("failed to save file: %v", err)
	}

	fmt.Println("Excel file 'stream_example.xlsx' created successfully.")
	fmt.Println("Helper function called")
}

func setFormat(f *excelize.File, sw *helper.StreamHelper) {
	styleDefs := map[string]*excelize.Style{
		"header": {
			Font:      &excelize.Font{Bold: true, Color: "#FFFFFF"},
			Fill:      excelize.Fill{Type: "pattern", Color: []string{"#4F81BD"}, Pattern: 1},
			Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		},
		"idColumn": {
			Alignment: &excelize.Alignment{Horizontal: "center"},
			Font:      &excelize.Font{Bold: true},
		},
		"nameColumn": {
			Alignment: &excelize.Alignment{Horizontal: "left"},
			Font:      &excelize.Font{Color: "#0066CC"},
		},
		"scoreColumn": {
			Alignment: &excelize.Alignment{Horizontal: "right"},
			NumFmt:    2,
			Font:      &excelize.Font{Bold: true},
			Fill:      excelize.Fill{Type: "pattern", Color: []string{"#E6F3FF"}, Pattern: 1},
		},
		"dateColumn": {
			Alignment: &excelize.Alignment{Horizontal: "center"},
			NumFmt:    14,
			Fill:      excelize.Fill{Type: "pattern", Color: []string{"#FFF2E6"}, Pattern: 1},
		},
		"timeColumn": {
			Alignment: &excelize.Alignment{Horizontal: "center"},
			NumFmt:    21,
			Fill:      excelize.Fill{Type: "pattern", Color: []string{"#F0F8FF"}, Pattern: 1},
		},
	}
	styleIDs, err := helper.CreateStyles(f, styleDefs)
	if err != nil {
		log.Fatalf("failed to create styles: %v", err)
	}

	if err := sw.MergeCells([]helper.MergeRange{
		{StartCell: "A1", EndCell: "A2"},
		{StartCell: "B1", EndCell: "B2"},
		{StartCell: "C1", EndCell: "C2"},
		{StartCell: "D1", EndCell: "E1"},
		{StartCell: "F1", EndCell: "G1"},
		{StartCell: "H1", EndCell: "I1"},
		{StartCell: "J1", EndCell: "J2"},
		{StartCell: "K1", EndCell: "K2"},
	}); err != nil {
		log.Fatalf("failed to merge header cells: %v", err)
	}

	if err := sw.SetColumnWidths(map[helper.ColumnRange]float64{
		{Start: 1, End: 1}:   20.0,  // ID
		{Start: 2, End: 2}:   30.0,  // Name
		{Start: 3, End: 3}:   15.0,  // Score
		{Start: 4, End: 4}:   12.00, // Open Date (Date)
		{Start: 5, End: 5}:   12.0,  // Open Date (Time)
		{Start: 6, End: 6}:   15.0,  // Request Date (Date)
		{Start: 7, End: 7}:   12.0,  // Request Date (Time)
		{Start: 8, End: 8}:   15.0,  // Accept Date (Date)
		{Start: 9, End: 9}:   12.0,  // Accept Date (Time)
		{Start: 10, End: 10}: 20.0,  // Number
		{Start: 11, End: 11}: 25.0,  // Account No.
	}); err != nil {
		log.Fatalf("failed to set column widths: %v", err)
	}

	if err := sw.SetColumnStyles([]helper.ColumnStyleConfig{
		{ColumnRange: helper.ColumnRange{Start: 1, End: 1}, Style: "idColumn"},     // ID
		{ColumnRange: helper.ColumnRange{Start: 2, End: 2}, Style: "nameColumn"},   // Name
		{ColumnRange: helper.ColumnRange{Start: 3, End: 3}, Style: "scoreColumn"},  // Score
		{ColumnRange: helper.ColumnRange{Start: 4, End: 4}, Style: "dateColumn"},   // Open Date (Date)
		{ColumnRange: helper.ColumnRange{Start: 5, End: 5}, Style: "timeColumn"},   // Open Date (Time)
		{ColumnRange: helper.ColumnRange{Start: 6, End: 6}, Style: "dateColumn"},   // Request Date (Date)
		{ColumnRange: helper.ColumnRange{Start: 7, End: 7}, Style: "timeColumn"},   // Request Date (Time)
		{ColumnRange: helper.ColumnRange{Start: 8, End: 8}, Style: "dateColumn"},   // Accept Date (Date)
		{ColumnRange: helper.ColumnRange{Start: 9, End: 9}, Style: "timeColumn"},   // Accept Date (Time)
		{ColumnRange: helper.ColumnRange{Start: 10, End: 10}, Style: "idColumn"},   // Number
		{ColumnRange: helper.ColumnRange{Start: 11, End: 11}, Style: "nameColumn"}, // Account No.
	}, styleIDs); err != nil {
		log.Fatalf("failed to set column styles: %v", err)
	}

	if err := sw.WriteHeader(helper.HeaderConfig{
		Cell: "A1",
		Values: []interface{}{
			excelize.Cell{Value: "ID", StyleID: styleIDs["header"]},
			excelize.Cell{Value: "Name", StyleID: styleIDs["header"]},
			excelize.Cell{Value: "Score", StyleID: styleIDs["header"]},
			excelize.Cell{Value: "Open Date", StyleID: styleIDs["header"]},
			nil,
			excelize.Cell{Value: "Request Date", StyleID: styleIDs["header"]},
			nil,
			excelize.Cell{Value: "Accept Date", StyleID: styleIDs["header"]},
			nil,
			excelize.Cell{Value: "Number", StyleID: styleIDs["header"]},
			excelize.Cell{Value: "Account No.", StyleID: styleIDs["header"]},
		},
	}); err != nil {
		log.Fatalf("failed to set header row 1: %v", err)
	}

	if err := sw.WriteHeader(helper.HeaderConfig{
		Cell: "A2",
		Values: []interface{}{
			nil, nil, nil,
			excelize.Cell{Value: "Date", StyleID: styleIDs["header"]},
			excelize.Cell{Value: "Time", StyleID: styleIDs["header"]},
			excelize.Cell{Value: "Date", StyleID: styleIDs["header"]},
			excelize.Cell{Value: "Time", StyleID: styleIDs["header"]},
			excelize.Cell{Value: "Date", StyleID: styleIDs["header"]},
			excelize.Cell{Value: "Time", StyleID: styleIDs["header"]},
		},
	}); err != nil {
		log.Fatalf("failed to set header row 2: %v", err)
	}
}
