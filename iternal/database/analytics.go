package database

import "os"

type SalesReport struct {
	ProductID  int64   `json:"product_id"`
	Name       string  `json:"name"`
	TotalSales float64 `json:"total_sales"`
}

type ProductSalesAverage struct {
	ProductID int64   `json:"product_id"`
	AvgSold   float64 `json:"avg_sold"`
}

func (db *Database) FetchSalesReport() ([]SalesReport, error) {
	query, err := os.ReadFile(analyticsPath + "sales_report.sql")
	if err != nil {
		db.Log.Error("Database FetchSalesReport() -> Read SQL file", err)
		return nil, err
	}

	rows, err := db.Query(string(query))
	if err != nil {
		db.Log.Error("Database FetchSalesReport() -> tx.Query()", err)
		return nil, err
	}
	defer rows.Close()

	var reports []SalesReport
	for rows.Next() {
		var report SalesReport
		if err := rows.Scan(&report.ProductID, &report.Name, &report.TotalSales); err != nil {
			db.Log.Error("Database FetchSalesReport() -> parsing rows", err)
			return nil, err
		}
		reports = append(reports, report)
	}

	if err := rows.Err(); err != nil {
		db.Log.Error("Database FetchSalesReport() -> rows.Err()", err)
		return nil, err
	}
	return reports, nil
}

func (db *Database) FetchRequirementsReport() ([]ProductSalesAverage, error) {
	query, err := os.ReadFile(analyticsPath + "requirements_report.sql")
	if err != nil {
		db.Log.Error("Database FetchRequirementsReport() -> Read SQL file", err)
		return nil, err
	}

	rows, err := db.Query(string(query))
	if err != nil {
		db.Log.Error("Database FetchRequirementsReport() -> tx.Query()", err)
		return nil, err
	}
	defer rows.Close()

	var reports []ProductSalesAverage
	for rows.Next() {
		var report ProductSalesAverage
		if err := rows.Scan(&report.ProductID, &report.AvgSold); err != nil {
			db.Log.Error("Database FetchRequirementsReport() -> parsing rows", err)
			return nil, err
		}
		reports = append(reports, report)
	}

	if err := rows.Err(); err != nil {
		db.Log.Error("Database FetchRequirementsReport() -> rows.Err()", err)
		return nil, err
	}
	return reports, nil
}
