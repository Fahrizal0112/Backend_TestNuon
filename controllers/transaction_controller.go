package controllers

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"transaction-api/config"
	"transaction-api/models"

	"github.com/gin-gonic/gin"
)

func CreateTransaction(c *gin.Context) {
	var transaction models.Transaction

	if err := c.ShouldBindJSON(&transaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Create(&transaction).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error":"Failed to create transaction"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message":"Transaction created successfully",
		"data": transaction,
	})
}

func GetAllTransactions(c *gin.Context) {
	var transactions []models.Transaction

	query := config.DB.Model(&models.Transaction{})

	// Tambahkan parameter search untuk pencarian global
	if search := c.Query("search"); search != "" {
		query = query.Where(
			"msisdn ILIKE ? OR trx_id ILIKE ? OR item ILIKE ? OR voucher_code ILIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%",
		)
	}

	if msisdn := c.Query("msisdn"); msisdn != "" {
		query = query.Where("msisdn = ?", msisdn)
	}

	if status := c.Query("status"); status != "" {
		statusInt, err := strconv.Atoi(status)
		if err == nil {
			query = query.Where("status = ?", statusInt)
		}
	}
	if startDate := c.Query("start_date"); startDate != "" {
		if parsedDate, err := time.Parse("2006-01-02", startDate); err == nil {
			query = query.Where("trx_date >= ?", parsedDate)
		}
	}

	if endDate := c.Query("end_date"); endDate != "" {
		if parsedDate, err := time.Parse("2006-01-02", endDate); err == nil {
			query = query.Where("trx_date <= ?", parsedDate.Add(24*time.Hour-time.Second))
		}
	}

	if item := c.Query("item"); item != "" {
		query = query.Where("item ILIKE ?", "%"+item+"%")
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "1000"))
	offset := (page - 1) * limit

	// Perbaiki query count untuk menggunakan filter yang sama
	countQuery := config.DB.Model(&models.Transaction{})
	if search := c.Query("search"); search != "" {
		countQuery = countQuery.Where(
			"msisdn ILIKE ? OR trx_id ILIKE ? OR item ILIKE ? OR voucher_code ILIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%",
		)
	}
	if msisdn := c.Query("msisdn"); msisdn != "" {
		countQuery = countQuery.Where("msisdn = ?", msisdn)
	}
	if status := c.Query("status"); status != "" {
		statusInt, err := strconv.Atoi(status)
		if err == nil {
			countQuery = countQuery.Where("status = ?", statusInt)
		}
	}
	if startDate := c.Query("start_date"); startDate != "" {
		if parsedDate, err := time.Parse("2006-01-02", startDate); err == nil {
			countQuery = countQuery.Where("trx_date >= ?", parsedDate)
		}
	}
	if endDate := c.Query("end_date"); endDate != "" {
		if parsedDate, err := time.Parse("2006-01-02", endDate); err == nil {
			countQuery = countQuery.Where("trx_date <= ?", parsedDate.Add(24*time.Hour-time.Second))
		}
	}
	if item := c.Query("item"); item != "" {
		countQuery = countQuery.Where("item ILIKE ?", "%"+item+"%")
	}

	if err := query.Offset(offset).Limit(limit).Order("trx_date DESC").Find(&transactions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transactions"})
		return
	}

	var total int64
	countQuery.Count(&total)

	c.JSON(http.StatusOK, gin.H{
		"data": transactions,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

func GetTransactionByID(c *gin.Context) {
	id := c.Param("id")
	var transaction models.Transaction

	if err := config.DB.First(&transaction, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": transaction})
}

func UploadCSV(c *gin.Context) {
	file, header, err := c.Request.FormFile("csv_file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded or invalid file"})
		return
	}
	defer file.Close()


	if !strings.HasSuffix(strings.ToLower(header.Filename), ".csv") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File must be a CSV file"})
		return
	}

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read CSV file: " + err.Error()})
		return
	}

	if len(records) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "CSV file is empty"})
		return
	}
	var successCount, errorCount int

	var errors []string

	for i, record := range records {
		if len(record) < 6 {
			errorCount++
			errors = append(errors, fmt.Sprintf("Row %d: insufficient columns (expected 6, got %d)", i+1, len(record)))
			continue
		}

		trxDate, err := time.Parse("2006-01-02 15:04:05", record[2])
		if err != nil {
			// Try alternative date format
			trxDate, err = time.Parse("2006-01-02", record[2])
			if err != nil {
				errorCount++
				errors = append(errors, fmt.Sprintf("Row %d: invalid date format '%s'", i+1, record[2]))
				continue
			}
		}
		status, err := strconv.Atoi(record[5])
		if err != nil {
			errorCount++
			errors = append(errors, fmt.Sprintf("Row %d: invalid status '%s'", i+1, record[5]))
			continue
		}

		transaction := models.Transaction{
			MSISDN:      record[0],
			TrxID:       record[1],
			TrxDate:     trxDate,
			Item:        record[3],
			VoucherCode: record[4],
			Status:      status,
		}

		var existingTransaction models.Transaction
		result := config.DB.Where("trx_id = ?", transaction.TrxID).First(&existingTransaction)
		
		if result.Error == nil {
			// Transaction exists, update it
			if err := config.DB.Model(&existingTransaction).Updates(transaction).Error; err != nil {
				errorCount++
				errors = append(errors, fmt.Sprintf("Row %d: failed to update transaction '%s': %v", i+1, transaction.TrxID, err))
				continue
			}
		} else {
			// Transaction doesn't exist, create new one
			if err := config.DB.Create(&transaction).Error; err != nil {
				errorCount++
				errors = append(errors, fmt.Sprintf("Row %d: failed to create transaction '%s': %v", i+1, transaction.TrxID, err))
				continue
			}
		}
		successCount++

	}
	response := gin.H{
		"message":       "CSV upload completed",
		"filename":      header.Filename,
		"total_rows":    len(records),
		"success_count": successCount,
		"error_count":   errorCount,
	}
	if len(errors) > 0 {
		response["errors"] = errors
	}

	if errorCount > 0 {
		c.JSON(http.StatusPartialContent, response)
	} else {
		c.JSON(http.StatusOK, response)
	}
}

func ClearAllTransactions(c *gin.Context) {
	// Add confirmation parameter
	confirm := c.Query("confirm")
	if confirm != "yes" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "This action will delete all transactions. Add ?confirm=yes to proceed",
		})
		return
	}

	// Delete all transactions
	result := config.DB.Unscoped().Delete(&models.Transaction{}, "1 = 1")
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear transactions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "All transactions cleared successfully",
		"deleted_count": result.RowsAffected,
	})
}