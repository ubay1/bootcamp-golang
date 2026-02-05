package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"kasir-api/models"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (repo *TransactionRepository) CreateTransaction(items []models.CheckoutItem) (*models.Transaction, error) {
	tx, err := repo.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	totalAmount := 0
	details := make([]models.TransactionDetail, 0)

	for _, item := range items {
		var productPrice, stock int
		var productName string

		err := tx.QueryRow("SELECT name, price, stock FROM products WHERE id = $1", item.ProductID).Scan(&productName, &productPrice, &stock)
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product id %d not found", item.ProductID)
		}
		if err != nil {
			return nil, err
		}

		subtotal := productPrice * item.Quantity
		totalAmount += subtotal

		_, err = tx.Exec("UPDATE products SET stock = stock - $1 WHERE id = $2", item.Quantity, item.ProductID)
		if err != nil {
			return nil, err
		}

		details = append(details, models.TransactionDetail{
			ProductID:   item.ProductID,
			ProductName: productName,
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		})
	}

	var transactionID int
	err = tx.QueryRow("INSERT INTO transactions (total_amount) VALUES ($1) RETURNING id", totalAmount).Scan(&transactionID)
	if err != nil {
		return nil, err
	}

	for i := range details {
		details[i].TransactionID = transactionID
		_, err = tx.Exec("INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal) VALUES ($1, $2, $3, $4)",
			transactionID, details[i].ProductID, details[i].Quantity, details[i].Subtotal)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &models.Transaction{
		ID:          transactionID,
		TotalAmount: totalAmount,
		Details:     details,
	}, nil
}

func (repo *TransactionRepository) GetSummaryToday() (*models.TransactionSummary, error) {
	var summary models.TransactionSummary
	err := repo.db.QueryRow("SELECT COUNT(*) AS total_transaksi, SUM(total_amount) AS total_revenue FROM transactions WHERE created_at::date = CURRENT_DATE").Scan(
		&summary.TotalTransactions, &summary.TotalRevenue)
	if err != nil {
		return nil, err
	}

	err = repo.db.QueryRow("SELECT p.name, SUM(td.quantity) AS qty_terjual FROM transaction_details td JOIN transactions t ON td.transaction_id = t.id JOIN products p ON td.product_id = p.id WHERE t.created_at::date = CURRENT_DATE GROUP BY p.name ORDER BY qty_terjual DESC LIMIT 1").Scan(
		&summary.ProductBestSeller.Name, &summary.ProductBestSeller.Quantity)
	if err != nil {
		return nil, err
	}

	return &summary, nil
}

func (repo *TransactionRepository) GetSummaryByDate(startDate string, endDate string) (*models.TransactionSummary, error) {
	var summary models.TransactionSummary
	err := repo.db.QueryRow("SELECT COUNT(*) AS total_transaksi, COALESCE(SUM(total_amount), 0) AS total_revenue FROM transactions WHERE created_at::date >= $1 AND created_at::date <= $2", startDate, endDate).Scan(
		&summary.TotalTransactions, &summary.TotalRevenue)
	if err != nil {
		return nil, err
	}

	if summary.TotalTransactions == 0 {
		return nil, errors.New("data tidak ditemukan")
	}

	err = repo.db.QueryRow("SELECT p.name, SUM(td.quantity) AS qty_terjual FROM transaction_details td JOIN transactions t ON td.transaction_id = t.id JOIN products p ON td.product_id = p.id WHERE t.created_at::date >= $1 AND t.created_at::date <= $2 GROUP BY p.name ORDER BY qty_terjual DESC LIMIT 1", startDate, endDate).Scan(
		&summary.ProductBestSeller.Name, &summary.ProductBestSeller.Quantity)
	if err != nil {
		return nil, err
	}

	return &summary, nil
}
