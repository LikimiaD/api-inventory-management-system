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

type Product struct {
	SupplierID  *int64    `json:"supplier_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       *float64  `json:"price"`
	Quantity    *int64    `json:"quantity"`
	Category    string    `json:"category"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (s *Server) addProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer r.Body.Close()

	user, err := s.getUserFromToken(r)
	if err != nil || user == "" {
		s.respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var p Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		s.respondWithError(w, http.StatusBadRequest, "Troubles with parsing data")
		return
	}

	if p.SupplierID == nil || p.Name == "" || p.Description == "" || p.Price == nil || p.Quantity == nil || p.Category == "" {
		s.respondWithError(w, http.StatusBadRequest, "Not enough information to create")
		return
	}

	if *p.Price < 0 {
		s.respondWithError(w, http.StatusBadRequest, "The price of the product cannot be negative")
		return
	}

	if *p.Quantity < 0 {
		s.respondWithError(w, http.StatusBadRequest, "Product quantities cannot be negative")
		return
	}

	if exists, err := s.DB.IDProduct(p.Name); err != nil {
		s.respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Problem on the server side, please report the error time %s", time.Now()))
		return
	} else if exists {
		s.respondWithError(w, http.StatusBadRequest, fmt.Sprintf("product with name %s exist", p.Name))
		return
	}

	if exists, err := s.DB.CheckSupplierExists(*p.SupplierID); err != nil {
		s.Log.Error("server side -> addProduct() -> s.DB.CheckSupplierExists()", err)
		s.respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Problem on the server side, please report the error time %s", time.Now()))
		return
	} else if !exists {
		s.respondWithError(w, http.StatusBadRequest, fmt.Sprintf("supplier with ID %d does not exist", *p.SupplierID))
		return
	}

	id, err := s.DB.AddProduct(*p.SupplierID, p.Name, p.Description, *p.Price, *p.Quantity, p.Category)
	if err != nil {
		s.respondWithError(w, http.StatusInternalServerError, "Problem when adding a new product")
		return
	}

	s.respondWithNew(w, id)
	s.Log.Info(fmt.Sprintf("User %s created new product with name %s and id %d", user, p.Name, id))
}

func (s *Server) deleteProduct(w http.ResponseWriter, r *http.Request) {
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

	if err = s.DB.DeleteProduct(deleteStruct.ID); errors.Is(err, database.ErrNoProductFound) {
		s.respondWithError(w, http.StatusBadRequest, "No product found with the provided ID")
		return
	} else if err != nil {
		s.respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Problem on the server side, please report the error time %s", time.Now()))
		return
	}

	s.respondNoContent(w)
	s.Log.Info(fmt.Sprintf("User %s removed product with id %d", user, deleteStruct.ID))
}

func (s *Server) updateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer r.Body.Close()

	user, err := s.getUserFromToken(r)
	if err != nil || user == "" {
		s.respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var updateStruct struct {
		ID          int64    `json:"id"`
		SupplierID  *int64   `json:"supplier_id"`
		Name        *string  `json:"name"`
		Description *string  `json:"description"`
		Price       *float64 `json:"price"`
		Quantity    *int64   `json:"quantity"`
		Category    *string  `json:"category"`
	}

	if err := json.NewDecoder(r.Body).Decode(&updateStruct); err != nil {
		s.respondWithError(w, http.StatusBadRequest, "Troubles with parsing data")
		return
	}

	if updateStruct.Price != nil && *updateStruct.Price < 0 {
		s.respondWithError(w, http.StatusBadRequest, "The price of the product cannot be negative")
		return
	}

	if updateStruct.Quantity != nil && *updateStruct.Quantity < 0 {
		s.respondWithError(w, http.StatusBadRequest, "Product quantities cannot be negative")
		return
	}

	if updateStruct.Name == nil {
		s.respondWithError(w, http.StatusBadRequest, "Name cannot be empty")
		return
	}

	if updateStruct.SupplierID == nil {
		s.respondWithError(w, http.StatusBadRequest, "Supplier ID cannot be empty")
		return
	}

	if exists, err := s.DB.IDProduct(*updateStruct.Name); err != nil {
		s.respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Problem on the server side, please report the error time %s", time.Now()))
		return
	} else if exists {
		s.respondWithError(w, http.StatusBadRequest, fmt.Sprintf("product with name %s exist", *updateStruct.Name))
		return
	}

	if exists, err := s.DB.CheckSupplierExists(*updateStruct.SupplierID); err != nil {
		s.Log.Error("server side -> addProduct() -> s.DB.CheckSupplierExists()", err)
		s.respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Problem on the server side, please report the error time %s", time.Now()))
		return
	} else if !exists {
		s.respondWithError(w, http.StatusBadRequest, fmt.Sprintf("supplier with ID %d does not exist", *updateStruct.SupplierID))
		return
	}

	if err := s.DB.UpdateProduct(updateStruct.ID, updateStruct.Name, updateStruct.SupplierID, updateStruct.Description, updateStruct.Price, updateStruct.Quantity, updateStruct.Category); err != nil {
		if errors.Is(err, database.ErrNoProductFound) {
			s.respondWithError(w, http.StatusBadRequest, "No product found with the provided ID")
		} else {
			s.respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Problem on the server side, please report the error time %s", time.Now()))
		}
		return
	}

	s.respondNoContent(w)
	s.Log.Info(fmt.Sprintf("User %s updated information of product with id %d", user, updateStruct.ID))
}

func (s *Server) exportProductsCSV(w io.Writer, products []database.Product) error {
	cw := csv.NewWriter(w)
	defer cw.Flush()

	headers := []string{"ProductID", "SupplierID", "Name", "Description", "Price", "Quantity", "Category", "CreatedAt", "UpdatedAt"}
	if err := cw.Write(headers); err != nil {
		return err
	}

	for _, product := range products {
		record := []string{
			fmt.Sprintf("%d", product.ProductID),
			fmt.Sprintf("%d", product.SupplierID),
			product.Name,
			product.Description,
			fmt.Sprintf("%.2f", product.Price),
			fmt.Sprintf("%d", product.Quantity),
			product.Category,
			product.CreatedAt.Format(time.RFC3339),
			product.UpdatedAt.Format(time.RFC3339),
		}
		if err := cw.Write(record); err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) exportProductsExcel(w io.Writer, products []database.Product) error {
	f := excelize.NewFile()
	sheetName := fmt.Sprintf("Products-%d", time.Now().Unix())
	f.NewSheet(sheetName)
	f.SetActiveSheet(f.NewSheet(sheetName))

	headers := []string{"ProductID", "SupplierID", "Name", "Description", "Price", "Quantity", "Category", "CreatedAt", "UpdatedAt"}
	for i, h := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheetName, cell, h)
	}

	for i, product := range products {
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", i+2), product.ProductID)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", i+2), product.SupplierID)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", i+2), product.Name)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", i+2), product.Description)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", i+2), fmt.Sprintf("%.2f", product.Price))
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", i+2), product.Quantity)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", i+2), product.Category)
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", i+2), product.CreatedAt.Format(time.RFC3339))
		f.SetCellValue(sheetName, fmt.Sprintf("I%d", i+2), product.UpdatedAt.Format(time.RFC3339))
	}

	if err := f.Write(w); err != nil {
		return err
	}
	return nil
}

func (s *Server) showInStockProducts(w http.ResponseWriter, r *http.Request) {
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

	products, err := s.DB.ShowNotEmptyQuantityProducts(limit)
	if err != nil {
		s.respondWithError(w, http.StatusInternalServerError, "Failed to retrieve products")
		return
	}

	switch format {
	case "json":
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(products); err != nil {
			s.respondWithError(w, http.StatusInternalServerError, "Error encoding response data")
			return
		}
	case "csv":
		w.Header().Set("Content-Type", "text/csv")
		if err := s.exportProductsCSV(w, products); err != nil {
			s.respondWithError(w, http.StatusInternalServerError, "Failed to generate CSV")
			return
		}
	case "excel":
		w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"products-%d.xlsx\"", time.Now().Unix()))
		if err := s.exportProductsExcel(w, products); err != nil {
			s.respondWithError(w, http.StatusInternalServerError, "Failed to generate Excel file")
			return
		}
	default:
		s.respondWithError(w, http.StatusBadRequest, "Invalid format specified")
		return
	}
	s.Log.Info(fmt.Sprintf("User %s requested information on the products in %s format", user, format))
}

func (s *Server) showCategoryProducts(w http.ResponseWriter, r *http.Request) {
	user, err := s.getUserFromToken(r)
	if err != nil || user == "" {
		s.respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	category := r.URL.Query().Get("category")
	if category == "" {
		s.respondWithError(w, http.StatusBadRequest, "Category parameter is required")
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

	products, err := s.DB.ShowByCategoryProducts(category, limit)
	if err != nil {
		s.respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to retrieve products for category '%s'", category))
		return
	}

	switch format {
	case "json":
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(products); err != nil {
			s.respondWithError(w, http.StatusInternalServerError, "Error encoding response data")
			return
		}
	case "csv":
		w.Header().Set("Content-Type", "text/csv")
		if err := s.exportProductsCSV(w, products); err != nil {
			s.respondWithError(w, http.StatusInternalServerError, "Failed to generate CSV")
			return
		}
	case "excel":
		w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"category_%s_products-%d.xlsx\"", category, time.Now().Unix()))
		if err := s.exportProductsExcel(w, products); err != nil {
			s.respondWithError(w, http.StatusInternalServerError, "Failed to generate Excel file")
			return
		}
	default:
		s.respondWithError(w, http.StatusBadRequest, "Invalid format specified")
		return
	}
	s.Log.Info(fmt.Sprintf("User %s requested information on category '%s' products in %s format", user, category, format))
}

func (s *Server) showPriceRangeProducts(w http.ResponseWriter, r *http.Request) {
	user, err := s.getUserFromToken(r)
	if err != nil || user == "" {
		s.respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	minPrice := r.URL.Query().Get("min")
	maxPrice := r.URL.Query().Get("max")
	if minPrice == "" || maxPrice == "" {
		s.respondWithError(w, http.StatusBadRequest, "Both 'min' and 'max' price parameters are required")
		return
	}

	min, err := strconv.ParseInt(minPrice, 10, 64)
	if err != nil {
		s.respondWithError(w, http.StatusBadRequest, "Invalid 'min' price value")
		return
	}

	max, err := strconv.ParseInt(maxPrice, 10, 64)
	if err != nil {
		s.respondWithError(w, http.StatusBadRequest, "Invalid 'max' price value")
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

	products, err := s.DB.ShowBetweenPriceProducts(min, max, limit)
	if err != nil {
		s.respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to retrieve products between prices %d and %d", min, max))
		return
	}

	switch format {
	case "json":
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(products); err != nil {
			s.respondWithError(w, http.StatusInternalServerError, "Error encoding response data")
			return
		}
	case "csv":
		w.Header().Set("Content-Type", "text/csv")
		if err := s.exportProductsCSV(w, products); err != nil {
			s.respondWithError(w, http.StatusInternalServerError, "Failed to generate CSV")
			return
		}
	case "excel":
		w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"products_price_range_%d_%d-%d.xlsx\"", min, max, time.Now().Unix()))
		if err := s.exportProductsExcel(w, products); err != nil {
			s.respondWithError(w, http.StatusInternalServerError, "Failed to generate Excel file")
			return
		}
	default:
		s.respondWithError(w, http.StatusBadRequest, "Invalid format specified")
		return
	}
	s.Log.Info(fmt.Sprintf("User %s requested information on products priced between %d and %d in %s format", user, min, max, format))
}

func (s *Server) showProducts(w http.ResponseWriter, r *http.Request) {
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

	products, err := s.DB.ShowProducts(limit)
	if err != nil {
		s.respondWithError(w, http.StatusInternalServerError, "Failed to retrieve products")
		return
	}

	switch format {
	case "json":
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(products); err != nil {
			s.respondWithError(w, http.StatusInternalServerError, "Error encoding response data")
			return
		}
	case "csv":
		w.Header().Set("Content-Type", "text/csv")
		if err := s.exportProductsCSV(w, products); err != nil {
			s.respondWithError(w, http.StatusInternalServerError, "Failed to generate CSV")
			return
		}
	case "excel":
		w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=products-%d.xlsx", time.Now().Unix()))
		if err := s.exportProductsExcel(w, products); err != nil {
			s.respondWithError(w, http.StatusInternalServerError, "Failed to generate Excel file")
			return
		}
	default:
		s.respondWithError(w, http.StatusBadRequest, "Invalid format specified")
		return
	}
	s.Log.Info(fmt.Sprintf("User %s requested information on products in %s format", user, format))
}

func exportPurchaseRequestsCSV(w io.Writer, requests []database.PurchaseRequest) error {
	cw := csv.NewWriter(w)
	defer cw.Flush()

	headers := []string{"ProductID", "Name", "SupplierID", "ContactEmail"}
	if err := cw.Write(headers); err != nil {
		return err
	}

	for _, req := range requests {
		record := []string{
			fmt.Sprintf("%d", req.ProductID),
			req.Name,
			fmt.Sprintf("%d", req.SupplierID),
			req.ContactEmail,
		}
		if err := cw.Write(record); err != nil {
			return err
		}
	}
	return nil
}

func exportPurchaseRequestsExcel(w io.Writer, requests []database.PurchaseRequest) error {
	f := excelize.NewFile()
	sheetName := "Purchase Requests"
	f.NewSheet(sheetName)
	f.SetActiveSheet(f.NewSheet(sheetName))

	headers := []string{"ProductID", "Name", "SupplierID", "ContactEmail"}
	for i, h := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheetName, cell, h)
	}

	for i, req := range requests {
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", i+2), req.ProductID)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", i+2), req.Name)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", i+2), req.SupplierID)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", i+2), req.ContactEmail)
	}

	if err := f.Write(w); err != nil {
		return err
	}
	return nil
}

func (s *Server) showPurchaseRequests(w http.ResponseWriter, r *http.Request) {
	user, err := s.getUserFromToken(r)
	if err != nil || user == "" {
		s.respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	maxQuantity := r.URL.Query().Get("maxQuantity")
	if maxQuantity == "" {
		s.respondWithError(w, http.StatusBadRequest, "maxQuantity parameter is required")
		return
	}

	maxQty, err := strconv.ParseInt(maxQuantity, 10, 64)
	if err != nil {
		s.respondWithError(w, http.StatusBadRequest, "Invalid maxQuantity value")
		return
	}

	format := r.URL.Query().Get("format")
	if format == "" {
		format = "json"
	}

	purchaseRequests, err := s.DB.PurchaseRequestProducts(maxQty)
	if err != nil {
		s.respondWithError(w, http.StatusInternalServerError, "Failed to retrieve purchase requests")
		return
	}

	switch format {
	case "json":
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(purchaseRequests); err != nil {
			s.respondWithError(w, http.StatusInternalServerError, "Error encoding response data")
			return
		}
	case "csv":
		w.Header().Set("Content-Type", "text/csv")
		if err := exportPurchaseRequestsCSV(w, purchaseRequests); err != nil {
			s.respondWithError(w, http.StatusInternalServerError, "Failed to generate CSV")
			return
		}
	case "excel":
		w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"purchase_requests-%d.xlsx\"", time.Now().Unix()))
		if err := exportPurchaseRequestsExcel(w, purchaseRequests); err != nil {
			s.respondWithError(w, http.StatusInternalServerError, "Failed to generate Excel file")
			return
		}
	default:
		s.respondWithError(w, http.StatusBadRequest, "Invalid format specified")
		return
	}
	s.Log.Info(fmt.Sprintf("User %s requested purchase requests in %s format", user, format))
}
