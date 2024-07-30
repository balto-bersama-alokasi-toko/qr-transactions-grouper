package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/schollz/progressbar/v3"
)

type merchantTrxDetail struct {
	Last1MonthTrx     int
	Last1MonthNominal float64
	TotalTrx          int
	TotalNominal      float64
	PostalCode        string
}

var merchantList = make(map[string]merchantTrxDetail)

func main() {
	// updateMerchantTransaction("merchant 1", true, 2000)
	// updateMerchantTransaction("merchant 1", false, 2000)
	// updateMerchantTransaction("merchant 1", true, 2000)
	// updateMerchantTransaction("merchant 1", true, 2000)
	// updateMerchantTransaction("merchant 2", true, 2000)
	// updateMerchantTransaction("merchant 2", true, 2000)

	// outputFile, err := os.Create("output.csv")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer outputFile.Close()

	// writer := csv.NewWriter(outputFile)
	// defer writer.Flush()
	// // write header
	// writer.Write([]string{"Merchant Name", "Total Transaction", "Total Nominal", "Last 1 Month Transaction", "Last 1 Month Nominal"})
	// for merchatName, merchantTrx := range merchantList {
	// 	row := []string{merchatName, strconv.Itoa(merchantTrx.TotalTrx), strconv.Itoa(merchantTrx.TotalNominal), strconv.Itoa(merchantTrx.Last1MonthTrx), strconv.Itoa(merchantTrx.Last1MonthNominal)}
	// 	err := writer.Write(row)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// }

	// totalData := 13000000
	// bar := progressbar.Default(int64(totalData))
	// for i := 0; i < totalData; i++ {
	// 	bar.Add(1)
	// }
	// floatNumber := 3.94
	// intNumber := int(floatNumber)
	// fmt.Println(intNumber)

	readTrxFile()
}

func updateMerchantTransaction(merchantName string, isLastMonth bool, trxAmount float64, postalCode string) {
	merchant, ok := merchantList[merchantName]
	if !ok {
		merchant = merchantTrxDetail{
			PostalCode: postalCode,
		}
	}
	if isLastMonth {
		merchant.Last1MonthNominal += trxAmount
		merchant.Last1MonthTrx += 1
	}

	merchant.TotalNominal += trxAmount
	merchant.TotalTrx += 1
	merchantList[merchantName] = merchant
}

func readTrxFile() {
	// read file
	file, err := os.Open("T_HSL_REK_30000_202407061935.csv")
	// Checks for the error
	if err != nil {
		log.Fatal("Error while reading the file", err)
	}

	// Closes the file
	defer file.Close()

	// The csv.NewReader() function is called in
	// which the object os.File passed as its parameter
	// and this creates a new csv.Reader that reads
	// from the file
	reader := csv.NewReader(file)

	// ReadAll reads all the records from the CSV file
	// and Returns them as slice of slices of string
	// and an error if any
	records, err := reader.ReadAll()

	// Checks for the error
	if err != nil {
		fmt.Println("Error reading records")
	}
	totalData := len(records)

	currentTime := time.Now()
	lastMonth := currentTime.AddDate(0, -1, 0)
	// Loop to iterate through
	// and print each of the string slice
	bar := progressbar.Default(int64(totalData))
	for _, record := range records {
		bar.Add(1)
		if record[21] != "00" {
			continue
		}
		transDate, err := time.Parse("20060102", record[12])
		if err != nil {
			fmt.Println("Error parsing data")
		}

		isLastMonth := false
		if transDate.After(lastMonth) {
			// transDate is within the last 1 month from today
			isLastMonth = true
		}

		trxAmount, err := strconv.ParseFloat(record[15], 64)
		if err != nil {
			fmt.Println("Error:", err)
		}

		updateMerchantTransaction(record[26], isLastMonth, trxAmount, record[45])

	}

	outputFile, err := os.Create("output.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()

	writer := csv.NewWriter(outputFile)
	defer writer.Flush()
	fmt.Println("")
	fmt.Println("")
	fmt.Println("Writing csv file..")
	// write header
	csvBar := progressbar.Default(int64(len(merchantList)))
	writer.Write([]string{"Merchant Name", "Postal Code", "Total Transaction", "Total Nominal", "Last 1 Month Transaction", "Last 1 Month Nominal"})
	for merchatName, merchantTrx := range merchantList {
		/**
		  minimal transaksi pertahun 10
		  minimal nominal transaksi 5000
		  minimal perbulan transaksi 1
		  minimal perbulan nominal 1000
		*/
		csvBar.Add(1)
		if merchantTrx.TotalTrx < 10 || merchantTrx.TotalNominal < 5000 || merchantTrx.Last1MonthTrx < 1 || merchantTrx.Last1MonthNominal < 1000 {
			continue
		}
		row := []string{merchatName, merchantTrx.PostalCode, strconv.Itoa(merchantTrx.TotalTrx), strconv.Itoa(int(merchantTrx.TotalNominal)), strconv.Itoa(merchantTrx.Last1MonthTrx), strconv.Itoa(int(merchantTrx.Last1MonthNominal))}
		err := writer.Write(row)
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println("")
	fmt.Println("")
	fmt.Println("Done!, Press any key to continue")
	fmt.Scanln()
}
