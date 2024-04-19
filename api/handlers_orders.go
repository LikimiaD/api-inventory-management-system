package api

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/likimiad/golang-restapi-inventory-managment-system/iternal/database"
	"net/http"
	"strconv"
	"time"
)

type OrderInformationInput struct {
	CustomerID *int64   `json:"customer_id"`
	Status     string   `json:"status"`
	ProductID  *int64   `json:"product_id"`
	Quantity   *int64   `json:"quantity"`
	Price      *float64 `json:"price"`
}

type OrderDetailInput struct {
	OrderID   *int64   `json:"order_id"`
	ProductID *int64   `json:"product_id"`
	Quantity  *int64   `json:"quantity"`
	Price     *float64 `json:"price"`
}

type OrderStatusUpdateInput struct {
	OrderID *int64  `json:"order_id"`
	Status  *string `json:"status"`
}

func (s *Server) addOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer r.Body.Close()

	user, err := s.getUserFromToken(r)
	if err != nil || user == "" {
		s.respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var o OrderInformationInput
	if err := json.NewDecoder(r.Body).Decode(&o); err != nil {
		s.respondWithError(w, http.StatusBadRequest, "Troubles with parsing data")
		return
	}

	if o.CustomerID == nil || o.Status == "" || o.ProductID == nil || o.Quantity == nil || o.Price == nil {
		s.respondWithError(w, http.StatusBadRequest, "Not enough information to create")
		return
	}

	if *o.Price < 0 {
		s.respondWithError(w, http.StatusBadRequest, "Price cannot be negative")
		return
	}

	if *o.Quantity < 0 {
		s.respondWithError(w, http.StatusBadRequest, "Quantities cannot be negative")
		return
	}

	if exists, err := s.DB.CheckCustomerExists(*o.CustomerID); err != nil {
		s.respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Problem on the server side, please report the error time %s", time.Now()))
		return
	} else if !exists {
		s.respondWithError(w, http.StatusBadRequest, fmt.Sprintf("customer with id %d  not exist", *o.CustomerID))
		return
	}

	if exists, err := s.DB.CheckProductExists(*o.ProductID); err != nil {
		s.respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Problem on the server side, please report the error time %s", time.Now()))
		return
	} else if !exists {
		s.respondWithError(w, http.StatusBadRequest, fmt.Sprintf("product with id %d not exist", *o.ProductID))
		return
	}

	if can, err := s.DB.CheckProductAvailability(*o.ProductID, *o.Quantity); err != nil {
		s.respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Problem on the server side, please report the error time %s", time.Now()))
		return
	} else if !can {
		s.respondWithError(w, http.StatusBadRequest, "Don't have that amount of product in stock")
		return
	}

	idOrder, idOrderDetail, err := s.DB.AddOrder(*o.CustomerID, *o.ProductID, *o.Quantity, *o.Price)
	if err != nil {
		s.respondWithError(w, http.StatusInternalServerError, "Problem when adding a new order")
		return
	}

	s.respondNewOrder(w, idOrder, idOrderDetail)
	s.Log.Info(fmt.Sprintf("User %s created new order with order id %d and order detail id %d", user, idOrder, idOrderDetail))
}

func (s *Server) refundOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer r.Body.Close()

	user, err := s.getUserFromToken(r)
	if err != nil || user == "" {
		s.respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var input struct {
		OrderID *int64 `json:"order_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		s.respondWithError(w, http.StatusBadRequest, "Troubles with parsing data")
		return
	}

	if input.OrderID == nil {
		s.respondWithError(w, http.StatusBadRequest, "Order ID is required")
		return
	}

	if exists, err := s.DB.CheckOrderExists(*input.OrderID); err != nil {
		s.respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Problem on the server side, please report the error time %s", time.Now()))
		return
	} else if !exists {
		s.respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Order with ID %d does not exist", *input.OrderID))
		return
	}

	err = s.DB.RefundOrder(*input.OrderID)
	if err != nil {
		s.respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to process refund for order %d", *input.OrderID))
		return
	}

	s.respondWithStatus(w, http.StatusOK, fmt.Sprintf("Refund processed successfully for order %d", *input.OrderID))
	s.Log.Info(fmt.Sprintf("User %s make refund order with ID %d", user, *input.OrderID))
}

func (s *Server) updateOrderStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer r.Body.Close()

	user, err := s.getUserFromToken(r)
	if err != nil || user == "" {
		s.respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var input OrderStatusUpdateInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		s.respondWithError(w, http.StatusBadRequest, "Troubles with parsing data")
		return
	}

	if input.OrderID == nil || input.Status == nil {
		s.respondWithError(w, http.StatusBadRequest, "Order ID and Status are required")
		return
	}

	currentStatus, err := s.DB.GetOrderStatus(*input.OrderID)
	if err != nil {
		s.respondWithError(w, http.StatusInternalServerError, "Failed to retrieve current order status")
		return
	}

	if currentStatus == "refunded" || currentStatus == "refund" {
		s.respondWithError(w, http.StatusBadRequest, "Updating status from 'refund' is not safe and not allowed")
		return
	}

	if *input.Status == "refunded" || *input.Status == "refund" {
		s.respondWithError(w, http.StatusBadRequest, "Use a special command to process refunds")
		return
	}

	if err = s.DB.UpdateStatusOrder(*input.OrderID, *input.Status); err != nil {
		s.respondWithError(w, http.StatusInternalServerError, "Failed to update order status")
		return
	}

	s.respondWithStatus(w, http.StatusOK, fmt.Sprintf("Order status updated successfully for order %d", *input.OrderID))
}

func exportOrdersCSV(w http.ResponseWriter, orders []database.OrderInfo) error {
	cw := csv.NewWriter(w)
	defer cw.Flush()

	headers := []string{"OrderID", "Status", "CreatedAt"}
	if err := cw.Write(headers); err != nil {
		return err
	}

	for _, order := range orders {
		record := []string{
			fmt.Sprintf("%d", order.OrderID),
			order.Status,
			order.CreatedAt.Format(time.RFC3339),
		}
		if err := cw.Write(record); err != nil {
			return err
		}
	}
	return nil
}

func exportOrdersExcel(w http.ResponseWriter, orders []database.OrderInfo) error {
	f := excelize.NewFile()
	sheetName := "Orders"
	f.NewSheet(sheetName)
	f.SetActiveSheet(f.NewSheet(sheetName))

	headers := []string{"OrderID", "Status", "CreatedAt"}
	for i, h := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheetName, cell, h)
	}

	for i, order := range orders {
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", i+2), order.OrderID)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", i+2), order.Status)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", i+2), order.CreatedAt.Format(time.RFC3339))
	}

	if err := f.Write(w); err != nil {
		return err
	}
	return nil
}

func (s *Server) showCustomerOrders(w http.ResponseWriter, r *http.Request) {
	user, err := s.getUserFromToken(r)
	if err != nil || user == "" {
		s.respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	customerIDParam := r.URL.Query().Get("customer_id")
	if customerIDParam == "" {
		s.respondWithError(w, http.StatusBadRequest, "Customer ID is required")
		return
	}

	customerID, err := strconv.ParseInt(customerIDParam, 10, 64)
	if err != nil {
		s.respondWithError(w, http.StatusBadRequest, "Invalid customer ID format")
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

	orders, err := s.DB.ShowByCustomerOrders(customerID, limit)
	if err != nil {
		s.respondWithError(w, http.StatusInternalServerError, "Failed to retrieve orders")
		return
	}

	switch format {
	case "json":
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(orders); err != nil {
			s.respondWithError(w, http.StatusInternalServerError, "Error encoding response data")
			return
		}
	case "csv":
		w.Header().Set("Content-Type", "text/csv")
		if err := exportOrdersCSV(w, orders); err != nil {
			s.respondWithError(w, http.StatusInternalServerError, "Failed to generate CSV")
			return
		}
	case "excel":
		w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"orders-%d.xlsx\"", time.Now().Unix()))
		if err := exportOrdersExcel(w, orders); err != nil {
			s.respondWithError(w, http.StatusInternalServerError, "Failed to generate Excel file")
			return
		}
	default:
		s.respondWithError(w, http.StatusBadRequest, "Invalid format specified")
		return
	}
	s.Log.Info(fmt.Sprintf("User %s requested orders for customer %d in %s format", user, customerID, format))
}

func (s *Server) showOrdersByDate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer r.Body.Close()

	user, err := s.getUserFromToken(r)
	if err != nil || user == "" {
		s.respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	startDateParam := r.URL.Query().Get("startDate")
	endDateParam := r.URL.Query().Get("endDate")

	if startDateParam == "" || endDateParam == "" {
		s.respondWithError(w, http.StatusBadRequest, "Both startDate and endDate parameters are required")
		return
	}

	startDate, err := time.Parse(time.RFC3339, startDateParam)
	if err != nil {
		s.respondWithError(w, http.StatusBadRequest, "Invalid startDate format, use ISO8601 format")
		return
	}

	endDate, err := time.Parse(time.RFC3339, endDateParam)
	if err != nil {
		s.respondWithError(w, http.StatusBadRequest, "Invalid endDate format, use ISO8601 format")
		return
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

	orders, err := s.DB.ShowByDateOrders(startDate, endDate, limit)
	if err != nil {
		s.respondWithError(w, http.StatusInternalServerError, "Failed to retrieve orders")
		return
	}

	if err := json.NewEncoder(w).Encode(orders); err != nil {
		s.respondWithError(w, http.StatusInternalServerError, "Error encoding response data")
		return
	}

	s.Log.Info(fmt.Sprintf("User %s requested orders between %s and %s", user, startDate.Format(time.RFC3339), endDate.Format(time.RFC3339)))
}

func (s *Server) showOrdersByStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer r.Body.Close()

	user, err := s.getUserFromToken(r)
	if err != nil || user == "" {
		s.respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	status := r.URL.Query().Get("status")
	if status == "" {
		s.respondWithError(w, http.StatusBadRequest, "Status parameter is required")
		return
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

	orders, err := s.DB.ShowByStatusOrders(status, limit)
	if err != nil {
		s.respondWithError(w, http.StatusInternalServerError, "Failed to retrieve orders")
		return
	}

	if err := json.NewEncoder(w).Encode(orders); err != nil {
		s.respondWithError(w, http.StatusInternalServerError, "Error encoding response data")
		return
	}

	s.Log.Info(fmt.Sprintf("User %s requested orders with status '%s'", user, status))
}

func exportOrdersFullCSV(w http.ResponseWriter, orders []database.Order) error {
	cw := csv.NewWriter(w)
	defer cw.Flush()

	headers := []string{"OrderID", "CustomerID", "Status", "CreatedAt", "UpdatedAt"}
	if err := cw.Write(headers); err != nil {
		return err
	}

	for _, order := range orders {
		record := []string{
			fmt.Sprintf("%d", order.OrderID),
			fmt.Sprintf("%d", order.CustomerID),
			order.Status,
			order.CreatedAt.Format(time.RFC3339),
			order.UpdatedAt.Format(time.RFC3339),
		}
		if err := cw.Write(record); err != nil {
			return err
		}
	}
	return nil
}

func exportOrdersFullExcel(w http.ResponseWriter, orders []database.Order) error {
	f := excelize.NewFile()
	sheetName := "Orders"
	f.NewSheet(sheetName)
	f.SetActiveSheet(f.NewSheet(sheetName))

	headers := []string{"OrderID", "CustomerID", "Status", "CreatedAt", "UpdatedAt"}
	for i, h := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheetName, cell, h)
	}

	for i, order := range orders {
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", i+2), order.OrderID)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", i+2), order.CustomerID)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", i+2), order.Status)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", i+2), order.CreatedAt.Format(time.RFC3339))
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", i+2), order.UpdatedAt.Format(time.RFC3339))
	}

	if err := f.Write(w); err != nil {
		return err
	}
	return nil
}

func (s *Server) showCustomerOrdersFull(w http.ResponseWriter, r *http.Request) {
	user, err := s.getUserFromToken(r)
	if err != nil || user == "" {
		s.respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	customerIDParam := r.URL.Query().Get("customer_id")
	if customerIDParam == "" {
		s.respondWithError(w, http.StatusBadRequest, "Customer ID is required")
		return
	}

	customerID, err := strconv.ParseInt(customerIDParam, 10, 64)
	if err != nil {
		s.respondWithError(w, http.StatusBadRequest, "Invalid customer ID format")
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

	orders, err := s.DB.ShowCustomerOrders(customerID, limit)
	if err != nil {
		s.respondWithError(w, http.StatusInternalServerError, "Failed to retrieve orders")
		return
	}

	switch format {
	case "json":
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(orders); err != nil {
			s.respondWithError(w, http.StatusInternalServerError, "Error encoding response data")
			return
		}
	case "csv":
		w.Header().Set("Content-Type", "text/csv")
		if err := exportOrdersFullCSV(w, orders); err != nil {
			s.respondWithError(w, http.StatusInternalServerError, "Failed to generate CSV")
			return
		}
	case "excel":
		w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"orders_full-%d.xlsx\"", time.Now().Unix()))
		if err := exportOrdersFullExcel(w, orders); err != nil {
			s.respondWithError(w, http.StatusInternalServerError, "Failed to generate Excel file")
			return
		}
	default:
		s.respondWithError(w, http.StatusBadRequest, "Invalid format specified")
		return
	}

	s.Log.Info(fmt.Sprintf("User %s requested orders for customer %d in %s format", user, customerID, format))
}

func exportOrderDetailsCSV(w http.ResponseWriter, details []database.OrderDetail) error {
	cw := csv.NewWriter(w)
	defer cw.Flush()

	headers := []string{"OrderDetailID", "Quantity", "Price", "Name"}
	if err := cw.Write(headers); err != nil {
		return err
	}

	for _, detail := range details {
		record := []string{
			fmt.Sprintf("%d", detail.OrderDetailID),
			fmt.Sprintf("%d", detail.Quantity),
			fmt.Sprintf("%.2f", detail.Price),
			detail.Name,
		}
		if err := cw.Write(record); err != nil {
			return err
		}
	}
	return nil
}

func exportOrderDetailsExcel(w http.ResponseWriter, details []database.OrderDetail) error {
	f := excelize.NewFile()
	sheetName := "Order Details"
	f.NewSheet(sheetName)
	f.SetActiveSheet(f.NewSheet(sheetName))

	headers := []string{"OrderDetailID", "Quantity", "Price", "Name"}
	for i, h := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheetName, cell, h)
	}

	for i, detail := range details {
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", i+2), detail.OrderDetailID)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", i+2), detail.Quantity)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", i+2), detail.Price)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", i+2), detail.Name)
	}

	if err := f.Write(w); err != nil {
		return err
	}
	return nil
}

func (s *Server) showOrderDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer r.Body.Close()

	user, err := s.getUserFromToken(r)
	if err != nil || user == "" {
		s.respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	orderIDParam := r.URL.Query().Get("order_id")
	if orderIDParam == "" {
		s.respondWithError(w, http.StatusBadRequest, "Order ID is required")
		return
	}

	orderID, err := strconv.ParseInt(orderIDParam, 10, 64)
	if err != nil {
		s.respondWithError(w, http.StatusBadRequest, "Invalid order ID format")
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

	orderDetails, err := s.DB.ShowOrderDetails(orderID, limit)
	if err != nil {
		s.respondWithError(w, http.StatusInternalServerError, "Failed to retrieve order details")
		return
	}

	switch format {
	case "json":
		if err := json.NewEncoder(w).Encode(orderDetails); err != nil {
			s.respondWithError(w, http.StatusInternalServerError, "Error encoding response data")
			return
		}
	case "csv":
		w.Header().Set("Content-Type", "text/csv")
		if err := exportOrderDetailsCSV(w, orderDetails); err != nil {
			s.respondWithError(w, http.StatusInternalServerError, "Failed to generate CSV")
			return
		}
	case "excel":
		w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"order_details-%d.xlsx\"", time.Now().Unix()))
		if err := exportOrderDetailsExcel(w, orderDetails); err != nil {
			s.respondWithError(w, http.StatusInternalServerError, "Failed to generate Excel file")
			return
		}
	default:
		s.respondWithError(w, http.StatusBadRequest, "Invalid format specified")
		return
	}

	s.Log.Info(fmt.Sprintf("User %s requested order details for order %d in %s format", user, orderID, format))
}

func (s *Server) showOrdersFullByStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer r.Body.Close()

	user, err := s.getUserFromToken(r)
	if err != nil || user == "" {
		s.respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	status := r.URL.Query().Get("status")
	if status == "" {
		s.respondWithError(w, http.StatusBadRequest, "Status parameter is required")
		return
	}

	var limit *int
	if limitParam := r.URL.Query().Get("limit"); limitParam != "" {
		limitParsed, err := strconv.Atoi(limitParam)
		if err != nil {
			s.respondWithError(w, http.StatusBadRequest, "Invalid limit value")
			return
		}
		limit = &limitParsed
	}

	orders, err := s.DB.ShowByStatusFullOrders(status, limit)
	if err != nil {
		s.respondWithError(w, http.StatusInternalServerError, "Failed to retrieve orders")
		return
	}

	format := r.URL.Query().Get("format")
	if format == "" {
		format = "json"
	}

	switch format {
	case "json":
		if err := json.NewEncoder(w).Encode(orders); err != nil {
			s.respondWithError(w, http.StatusInternalServerError, "Error encoding response data")
			return
		}
	case "csv":
		w.Header().Set("Content-Type", "text/csv")
		if err := exportOrdersFullCSV(w, orders); err != nil {
			s.respondWithError(w, http.StatusInternalServerError, "Failed to generate CSV")
			return
		}
	case "excel":
		w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"orders-%d.xlsx\"", time.Now().Unix()))
		if err := exportOrdersFullExcel(w, orders); err != nil {
			s.respondWithError(w, http.StatusInternalServerError, "Failed to generate Excel file")
			return
		}
	default:
		s.respondWithError(w, http.StatusBadRequest, "Invalid format specified")
		return
	}
	s.Log.Info(fmt.Sprintf("User %s requested orders with status '%s' in %s format", user, status, format))
}
