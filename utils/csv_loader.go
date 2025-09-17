package utils

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"
	"transaction-api/config"
	"transaction-api/models"
)

func LoadCSVData(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("error opening CSV file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("error reading CSV: %v", err)
	}

	for i, record := range records {
		if len(record) < 6 {
			fmt.Printf("Skipping row %d: insufficient columns\n", i+1)
			continue
		}

		trxDate, err := time.Parse("2006-01-02 15:04:05", record[2])
		if err != nil {
			fmt.Printf("Error parsing date in row %d: %v\n", i+1, err)
			continue
		}

		status, err := strconv.Atoi(record[5])
		if err != nil {
			fmt.Printf("Error parsing status in row %d: %v\n", i+1, err)
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

		if err := config.DB.Create(&transaction).Error; err != nil {
			fmt.Printf("Error creating transaction in row %d: %v\n", i+2, err)
			continue
		}
	}
	
	fmt.Println("CSV data loaded successfully!")
	return nil
}