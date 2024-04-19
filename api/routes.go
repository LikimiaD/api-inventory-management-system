package api

import (
	"net/http"
)

func (s *Server) routes() {
	s.Router.HandleFunc("/login", s.login).Methods("POST")

	s.Router.Handle("/add_supplier", s.isAuthorized(http.HandlerFunc(s.addSupplier))).Methods("POST")
	s.Router.Handle("/delete_supplier", s.isAuthorized(http.HandlerFunc(s.deleteSupplier))).Methods("POST")
	s.Router.Handle("/update_supplier", s.isAuthorized(http.HandlerFunc(s.updateSupplier))).Methods("POST")
	s.Router.Handle("/show_suppliers", s.isAuthorized(http.HandlerFunc(s.showSuppliers))).Methods("GET")

	s.Router.Handle("/add_customer", s.isAuthorized(http.HandlerFunc(s.addCustomer))).Methods("POST")
	s.Router.Handle("/delete_customer", s.isAuthorized(http.HandlerFunc(s.deleteCustomer))).Methods("POST")
	s.Router.Handle("/update_customer", s.isAuthorized(http.HandlerFunc(s.updateCustomer))).Methods("POST")
	s.Router.Handle("/show_customers", s.isAuthorized(http.HandlerFunc(s.showCustomers))).Methods("GET")

	s.Router.Handle("/add_product", s.isAuthorized(http.HandlerFunc(s.addProduct))).Methods("POST")
	s.Router.Handle("/delete_product", s.isAuthorized(http.HandlerFunc(s.deleteProduct))).Methods("POST")
	s.Router.Handle("/update_product", s.isAuthorized(http.HandlerFunc(s.updateProduct))).Methods("POST")
	s.Router.Handle("/show_products", s.isAuthorized(http.HandlerFunc(s.showProducts))).Methods("GET")
	s.Router.Handle("/show_in_stock_products", s.isAuthorized(http.HandlerFunc(s.showInStockProducts))).Methods("GET")
	s.Router.Handle("/show_category_products", s.isAuthorized(http.HandlerFunc(s.showCategoryProducts))).Methods("GET")
	s.Router.Handle("/show_price_products", s.isAuthorized(http.HandlerFunc(s.showPriceRangeProducts))).Methods("GET")
	s.Router.Handle("/show_purchase_request", s.isAuthorized(http.HandlerFunc(s.showPurchaseRequests))).Methods("GET")

	s.Router.Handle("/add_order", s.isAuthorized(http.HandlerFunc(s.addOrder))).Methods("POST")
	s.Router.Handle("/refund_order", s.isAuthorized(http.HandlerFunc(s.refundOrder))).Methods("POST")
	s.Router.Handle("/update_order", s.isAuthorized(http.HandlerFunc(s.updateOrderStatus))).Methods("POST")
	s.Router.Handle("/show_customer_orders", s.isAuthorized(http.HandlerFunc(s.showCustomerOrders))).Methods("GET")
	s.Router.Handle("/show_customer_orders_full", s.isAuthorized(http.HandlerFunc(s.showCustomerOrdersFull))).Methods("GET")
	s.Router.Handle("/show_orders_by_date", s.isAuthorized(http.HandlerFunc(s.showOrdersByDate))).Methods("GET")
	s.Router.Handle("/show_orders_by_status", s.isAuthorized(http.HandlerFunc(s.showOrdersByStatus))).Methods("GET")
	s.Router.Handle("/show_orders_by_status_full", s.isAuthorized(http.HandlerFunc(s.showOrdersFullByStatus))).Methods("GET")
	s.Router.Handle("/show_order_details", s.isAuthorized(http.HandlerFunc(s.showOrderDetails))).Methods("GET")

	s.Router.Handle("/sales_report", s.isAuthorized(http.HandlerFunc(s.showSalesReport))).Methods("GET")
	s.Router.Handle("/requirements_report", s.isAuthorized(http.HandlerFunc(s.showRequirementsReport))).Methods("GET")
}
