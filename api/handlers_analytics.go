package api

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/likimiad/golang-restapi-inventory-managment-system/iternal/database"
	"net/http"
	"time"
)

func (s *Server) showSalesReport(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer r.Body.Close()

	user, err := s.getUserFromToken(r)
	if err != nil || user == "" {
		s.respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	reports, err := s.DB.FetchSalesReport()
	if err != nil {
		s.respondWithError(w, http.StatusInternalServerError, "Failed to retrieve sales report")
		return
	}

	format := r.URL.Query().Get("format")
	if format == "" {
		format = "json"
	}

	switch format {
	case "json":
		if err := json.NewEncoder(w).Encode(reports); err != nil {
			s.respondWithError(w, http.StatusInternalServerError, "Error encoding response data")
			return
		}
	case "csv":
		w.Header().Set("Content-Type", "text/csv")
		if err := exportSalesReportCSV(w, reports); err != nil {
			s.respondWithError(w, http.StatusInternalServerError, "Failed to generate CSV")
			return
		}
	case "excel":
		w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"sales_report-%d.xlsx\"", time.Now().Unix()))
		if err := exportSalesReportExcel(w, reports); err != nil {
			s.respondWithError(w, http.StatusInternalServerError, "Failed to generate Excel file")
			return
		}
	default:
		s.respondWithError(w, http.StatusBadRequest, "Invalid format specified")
		return
	}
	s.Log.Info(fmt.Sprintf("User %s accessed the sales report in %s format", user, format))
}

func exportSalesReportCSV(w http.ResponseWriter, reports []database.SalesReport) error {
	cw := csv.NewWriter(w)
	defer cw.Flush()

	headers := []string{"OrderID", "TotalSales"}
	if err := cw.Write(headers); err != nil {
		return err
	}

	for _, report := range reports {
		record := []string{
			fmt.Sprintf("%d", report.ProductID),
			report.Name,
			fmt.Sprintf("%.2f", report.TotalSales),
		}
		if err := cw.Write(record); err != nil {
			return err
		}
	}
	return nil
}

func exportSalesReportExcel(w http.ResponseWriter, reports []database.SalesReport) error {
	f := excelize.NewFile()
	sheetName := "Sales Report"
	f.NewSheet(sheetName)
	f.SetActiveSheet(f.NewSheet(sheetName))

	headers := []string{"OrderID", "TotalSales"}
	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheetName, cell, header)
	}

	for i, report := range reports {
		row := i + 2
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), report.ProductID)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), report.Name)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), report.TotalSales)
	}

	if err := f.Write(w); err != nil {
		return err
	}
	return nil
}

func exportRequirementsReportCSV(w http.ResponseWriter, reports []database.ProductSalesAverage) error {
	cw := csv.NewWriter(w)
	defer cw.Flush()

	headers := []string{"ProductID", "AvgSold"}
	if err := cw.Write(headers); err != nil {
		return err
	}

	for _, report := range reports {
		record := []string{
			fmt.Sprintf("%d", report.ProductID),
			fmt.Sprintf("%.2f", report.AvgSold),
		}
		if err := cw.Write(record); err != nil {
			return err
		}
	}
	return nil
}

func exportRequirementsReportExcel(w http.ResponseWriter, reports []database.ProductSalesAverage) error {
	f := excelize.NewFile()
	sheetName := "Average Sales"
	f.NewSheet(sheetName)
	f.SetActiveSheet(f.NewSheet(sheetName))

	headers := []string{"ProductID", "AvgSold"}
	for i, h := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheetName, cell, h)
	}

	for i, report := range reports {
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", i+2), report.ProductID)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", i+2), fmt.Sprintf("%.2f", report.AvgSold))
	}

	if err := f.Write(w); err != nil {
		return err
	}
	return nil
}

func (s *Server) showRequirementsReport(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer r.Body.Close()

	user, err := s.getUserFromToken(r)
	if err != nil || user == "" {
		s.respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	reports, err := s.DB.FetchRequirementsReport()
	if err != nil {
		s.respondWithError(w, http.StatusInternalServerError, "Failed to retrieve requirements report")
		return
	}

	format := r.URL.Query().Get("format")
	if format == "" {
		format = "json"
	}

	switch format {
	case "json":
		if err := json.NewEncoder(w).Encode(reports); err != nil {
			s.respondWithError(w, http.StatusInternalServerError, "Error encoding response data")
			return
		}
	case "csv":
		w.Header().Set("Content-Type", "text/csv")
		if err := exportRequirementsReportCSV(w, reports); err != nil {
			s.respondWithError(w, http.StatusInternalServerError, "Failed to generate CSV")
			return
		}
	case "excel":
		w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"requirements_report-%d.xlsx\"", time.Now().Unix()))
		if err := exportRequirementsReportExcel(w, reports); err != nil {
			s.respondWithError(w, http.StatusInternalServerError, "Failed to generate Excel file")
			return
		}
	default:
		s.respondWithError(w, http.StatusBadRequest, "Invalid format specified")
		return
	}

	s.Log.Info(fmt.Sprintf("User %s accessed the requirements report in %s format", user, format))
}
