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

type Customer struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
}

func (s *Server) addCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer r.Body.Close()

	user, err := s.getUserFromToken(r)
	if err != nil || user == "" {
		s.respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var cus Customer
	if err := json.NewDecoder(r.Body).Decode(&cus); err != nil {
		s.respondWithError(w, http.StatusBadRequest, "Troubles with parsing data")
		return
	}

	if cus.Name == "" || cus.Email == "" || cus.Phone == "" || cus.Address == "" {
		s.respondWithError(w, http.StatusBadRequest, "Not enough information to create")
		return
	}

	if id, err := s.DB.CheckEmailCustomer(cus.Email); err != nil {
		s.respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Problem on the server side, please report the error time %s", time.Now()))
		return
	} else if id != -1 {
		s.respondWithError(w, http.StatusBadRequest, "There is a user with this email")
		return
	}

	if id, err := s.DB.CheckPhoneCustomer(cus.Phone); err != nil {
		s.respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Problem on the server side, please report the error time %s", time.Now()))
		return
	} else if id != -1 {
		s.respondWithError(w, http.StatusBadRequest, "There is a user with this number")
		return
	}

	id, err := s.DB.AddCustomer(cus.Name, cus.Email, cus.Phone, cus.Address)
	if err != nil {
		s.respondWithError(w, http.StatusInternalServerError, "Problem when adding a new user")
		return
	}

	s.respondWithNew(w, id)
	s.Log.Info(fmt.Sprintf("User %s created new customer with id %d", user, id))
}

func (s *Server) deleteCustomer(w http.ResponseWriter, r *http.Request) {
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

	if err = s.DB.DeleteCustomer(deleteStruct.ID); errors.Is(err, database.ErrNoCustomerFound) {
		s.respondWithError(w, http.StatusBadRequest, "No customer found with the provided ID")
		return
	} else if err != nil {
		s.respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Problem on the server side, please report the error time %s", time.Now()))
		return
	}

	s.respondNoContent(w)
	s.Log.Info(fmt.Sprintf("User %s removed customer with id %d", user, deleteStruct.ID))
}

func (s *Server) updateCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer r.Body.Close()

	user, err := s.getUserFromToken(r)
	if err != nil || user == "" {
		s.respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var updateStruct struct {
		ID      int64   `json:"id"`
		Name    *string `json:"name"`
		Email   *string `json:"email"`
		Phone   *string `json:"phone"`
		Address *string `json:"address"`
	}

	if err := json.NewDecoder(r.Body).Decode(&updateStruct); err != nil {
		s.respondWithError(w, http.StatusBadRequest, "Troubles with parsing data")
		return
	}

	if err := s.DB.UpdateCustomer(updateStruct.ID, updateStruct.Name, updateStruct.Email, updateStruct.Phone, updateStruct.Address); err != nil {
		if errors.Is(err, database.ErrNoCustomerFound) {
			s.respondWithError(w, http.StatusBadRequest, "No customer found with the provided ID")
		} else {
			s.respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Problem on the server side, please report the error time %s", time.Now()))
		}
		return
	}

	s.respondNoContent(w)
	s.Log.Info(fmt.Sprintf("User %s updated information customer with id %d", user, updateStruct.ID))
}

func (s *Server) exportCustomersCSV(w io.Writer, customers []database.Customer) error {
	cw := csv.NewWriter(w)
	defer cw.Flush()

	headers := []string{"ID", "Name", "Email", "Phone", "Address"}
	if err := cw.Write(headers); err != nil {
		return err
	}

	for _, customer := range customers {
		record := []string{
			fmt.Sprintf("%d", customer.ID),
			customer.Name,
			customer.Email,
			customer.Phone,
			customer.Address,
		}
		if err := cw.Write(record); err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) exportCustomersExcel(w io.Writer, customers []database.Customer) error {
	f := excelize.NewFile()
	sheetName := fmt.Sprintf("Customers-%d", time.Now().Unix())
	f.NewSheet(sheetName)
	f.SetActiveSheet(f.NewSheet(sheetName))

	headers := []string{"ID", "Name", "Email", "Phone", "Address"}
	for i, h := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheetName, cell, h)
	}

	for i, customer := range customers {
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", i+2), customer.ID)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", i+2), customer.Name)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", i+2), customer.Email)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", i+2), customer.Phone)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", i+2), customer.Address)
	}

	if err := f.Write(w); err != nil {
		return err
	}
	return nil
}

func (s *Server) showCustomers(w http.ResponseWriter, r *http.Request) {
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

	customers, err := s.DB.ShowCustomers(limit)
	if err != nil {
		s.respondWithError(w, http.StatusInternalServerError, "Failed to retrieve customers")
		return
	}

	switch format {
	case "json":
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(customers); err != nil {
			s.respondWithError(w, http.StatusInternalServerError, "Error encoding response data")
			return
		}
	case "csv":
		w.Header().Set("Content-Type", "text/csv")
		if err := s.exportCustomersCSV(w, customers); err != nil {
			s.respondWithError(w, http.StatusInternalServerError, "Failed to generate CSV")
			return
		}
	case "excel":
		w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=customers-%d.xlsx", time.Now().Unix()))
		if err := s.exportCustomersExcel(w, customers); err != nil {
			s.respondWithError(w, http.StatusInternalServerError, "Failed to generate Excel file")
			return
		}
	default:
		s.respondWithError(w, http.StatusBadRequest, "Invalid format specified")
		return
	}
	s.Log.Info(fmt.Sprintf("User %s requested information on the customers in %s format", user, format))
}
