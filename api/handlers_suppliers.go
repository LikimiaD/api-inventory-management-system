package api

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/likimiad/golang-restapi-inventory-managment-system/iternal/database"
	"io"
	"net/http"
	"strconv"
	"time"
)

type Supplier struct {
	Name         string `json:"name"`
	ContactName  string `json:"contact_name"`
	ContactEmail string `json:"contact_email"`
	ContactPhone string `json:"contact_phone"`
}

func (s *Server) addSupplier(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer r.Body.Close()

	user, err := s.getUserFromToken(r)
	if err != nil || user == "" {
		s.respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var sup Supplier
	if err := json.NewDecoder(r.Body).Decode(&sup); err != nil {
		s.respondWithError(w, http.StatusBadRequest, "Troubles with parsing data")
		return
	}

	if sup.Name == "" || sup.ContactName == "" || sup.ContactEmail == "" || sup.ContactPhone == "" {
		s.respondWithError(w, http.StatusBadRequest, "Not enough information to create")
		return
	}

	if id, err := s.DB.CheckEmailSupplier(sup.ContactEmail); err != nil {
		s.respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Problem on the server side, please report the error time %s", time.Now()))
		return
	} else if id != -1 {
		s.respondWithError(w, http.StatusBadRequest, "There is a user with this email")
		return
	}

	if id, err := s.DB.CheckPhoneSupplier(sup.ContactPhone); err != nil {
		s.respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Problem on the server side, please report the error time %s", time.Now()))
		return
	} else if id != -1 {
		s.respondWithError(w, http.StatusBadRequest, "There is a user with this number")
		return
	}

	id, err := s.DB.AddSupplier(sup.Name, sup.ContactName, sup.ContactEmail, sup.ContactPhone)
	if err != nil {
		s.respondWithError(w, http.StatusInternalServerError, "Problem when adding a new user")
		return
	}

	s.respondWithNew(w, id)
	s.Log.Info(fmt.Sprintf("User %s created new supplier with id %d", user, id))
}

func (s *Server) deleteSupplier(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer r.Body.Close()

	user, err := s.getUserFromToken(r)
	if err != nil || user == "" {
		s.respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var deleteStruct struct {
		ID int64 `json:"id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&deleteStruct); err != nil {
		s.respondWithError(w, http.StatusBadRequest, "Troubles with parsing data")
		return
	}

	if err = s.DB.DeleteSupplier(deleteStruct.ID); errors.Is(err, database.ErrNoSupplierFound) {
		s.respondWithError(w, http.StatusBadRequest, "No supplier found with the provided ID")
		return
	} else if err != nil {
		s.respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Problem on the server side, please report the error time %s", time.Now()))
		return
	}

	s.respondNoContent(w)
	s.Log.Info(fmt.Sprintf("User %s removed supplier with id %d", user, deleteStruct.ID))
}

func (s *Server) updateSupplier(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer r.Body.Close()

	user, err := s.getUserFromToken(r)
	if err != nil || user == "" {
		s.respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var updateStruct struct {
		ID           int64   `json:"id"`
		Name         *string `json:"name"`
		ContactName  *string `json:"contact_name"`
		ContactEmail *string `json:"contact_email"`
		ContactPhone *string `json:"contact_phone"`
	}

	if err := json.NewDecoder(r.Body).Decode(&updateStruct); err != nil {
		s.respondWithError(w, http.StatusBadRequest, "Troubles with parsing data")
		return
	}

	if err := s.DB.UpdateSupplier(updateStruct.ID, updateStruct.Name, updateStruct.ContactName, updateStruct.ContactEmail, updateStruct.ContactPhone); err != nil {
		if errors.Is(err, database.ErrNoSupplierFound) {
			s.respondWithError(w, http.StatusBadRequest, "No supplier found with the provided ID")
		} else {
			s.respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Problem on the server side, please report the error time %s", time.Now()))
		}
		return
	}

	s.respondNoContent(w)
	s.Log.Info(fmt.Sprintf("User %s updated information supplier with id %d", user, updateStruct.ID))
}

func (s *Server) exportSuppliersCSV(w io.Writer, suppliers []database.Supplier) error {
	cw := csv.NewWriter(w)
	defer cw.Flush()

	headers := []string{"ID", "Name", "Contact Name", "Contact Email", "Contact Phone"}
	if err := cw.Write(headers); err != nil {
		return err
	}

	for _, supplier := range suppliers {
		record := []string{
			fmt.Sprintf("%d", supplier.SupplierID),
			supplier.Name,
			supplier.ContactName,
			supplier.ContactEmail,
			supplier.ContactPhone,
		}
		if err := cw.Write(record); err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) exportSuppliersExcel(w io.Writer, suppliers []database.Supplier) error {
	f := excelize.NewFile()
	sheetName := fmt.Sprintf("Suppliers-%d", time.Now().Unix())
	f.NewSheet(sheetName)
	f.SetActiveSheet(f.NewSheet(sheetName))

	headers := []string{"ID", "Name", "Contact Name", "Contact Email", "Contact Phone"}
	for i, h := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheetName, cell, h)
	}

	for i, supplier := range suppliers {
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", i+2), supplier.SupplierID)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", i+2), supplier.Name)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", i+2), supplier.ContactName)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", i+2), supplier.ContactEmail)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", i+2), supplier.ContactPhone)
	}

	if err := f.Write(w); err != nil {
		return err
	}
	return nil
}

func (s *Server) showSuppliers(w http.ResponseWriter, r *http.Request) {
	user, err := s.getUserFromToken(r)
	if err != nil || user == "" {
		s.respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	format := r.URL.Query().Get("format")
	if format == "" {
		format = "json"
	}

	var limit *int
	if l := r.URL.Query().Get("limit"); l != "" {
		if lmt, err := strconv.Atoi(l); err == nil {
			limit = &lmt
		} else {
			s.respondWithError(w, http.StatusBadRequest, "Invalid limit value")
			return
		}
	}

	suppliers, err := s.DB.ShowSuppliers(limit)
	if err != nil {
		s.respondWithError(w, http.StatusInternalServerError, "Failed to retrieve suppliers")
		return
	}

	switch format {
	case "json":
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(suppliers); err != nil {
			s.respondWithError(w, http.StatusInternalServerError, "Error encoding response data")
			return
		}
	case "csv":
		w.Header().Set("Content-Type", "text/csv")
		if err := s.exportSuppliersCSV(w, suppliers); err != nil {
			s.respondWithError(w, http.StatusInternalServerError, "Failed to generate CSV")
			return
		}
	case "excel":
		w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=suppliers-%d.xlsx", time.Now().Unix()))
		if err := s.exportSuppliersExcel(w, suppliers); err != nil {
			s.respondWithError(w, http.StatusInternalServerError, "Failed to generate Excel file")
			return
		}
	default:
		s.respondWithError(w, http.StatusBadRequest, "Invalid format specified")
		return
	}
	s.Log.Info(fmt.Sprintf("User %s requested information on the suppliers in %s format", user, format))
}
