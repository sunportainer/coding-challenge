package service

import (
	"encoding/json"
	"fmt"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"io/ioutil"
	"os"
)

//Generate Golang Struct to map json entity
type BookKeeping struct {
	BalanceDate  string `json:"balance_date"`
	ConnectionID string `json:"connection_id"`
	Currency     string `json:"currency"`
	Data         []struct {
		AccountCategory   string  `json:"account_category"`
		AccountCode       string  `json:"account_code"`
		AccountCurrency   string  `json:"account_currency"`
		AccountIdentifier string  `json:"account_identifier"`
		AccountName       string  `json:"account_name"`
		AccountStatus     string  `json:"account_status"`
		AccountType       string  `json:"account_type"`
		AccountTypeBank   string  `json:"account_type_bank"`
		SystemAccount     string  `json:"system_account"`
		TotalValue        float64 `json:"total_value"`
		ValueType         string  `json:"value_type"`
	} `json:"data"`
	ObjectCategory       string `json:"object_category"`
	ObjectClass          string `json:"object_class"`
	ObjectCreationDate   string `json:"object_creation_date"`
	ObjectOriginCategory string `json:"object_origin_category"`
	ObjectOriginType     string `json:"object_origin_type"`
	ObjectType           string `json:"object_type"`
	User                 string `json:"user"`
}

func LoadFile() {

	jsonFile, err := os.Open("./data/data.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var books BookKeeping

	json.Unmarshal(byteValue, &books)

	var revenue, expense float64
	var grossProfitMarginTotal float64
	var workingCapitalRatio, liabilities float64

	for _, record := range books.Data {
		if record.AccountCategory == "revenue" {
			revenue = revenue + record.TotalValue
		}
		if record.AccountCategory == "expense" {
			expense = expense + record.TotalValue
		}
		if record.AccountType == "sales" && record.ValueType == "debit" {
			grossProfitMarginTotal = grossProfitMarginTotal + record.TotalValue
		}

		if record.AccountCategory == "assets" &&
			(record.AccountType == "current" || record.AccountType == "bank" || record.AccountType == "current_accounts_receivable") {
			if record.ValueType == "debit" {
				workingCapitalRatio = workingCapitalRatio + record.TotalValue
			}
			if record.ValueType == "credit" {
				workingCapitalRatio = workingCapitalRatio - record.TotalValue
			}
		}

		if record.AccountCategory == "liability" &&
			(record.AccountType == "current" || record.AccountType == "current_accounts_payable") {
			if record.ValueType == "debit" {
				liabilities = liabilities - record.TotalValue
			}
			if record.ValueType == "credit" {
				liabilities = liabilities + record.TotalValue
			}
		}
	}
	p := message.NewPrinter(language.English)
	fmt.Printf(p.Sprintf("Revenue: $%.0f\n", revenue))
	fmt.Printf(p.Sprintf("Revenue: $%.0f\n", expense))
	grossProfitMargin := grossProfitMarginTotal / revenue
	netProfitMargin := (revenue - expense) / revenue
	fmt.Printf("Gross Profit Margin:%.0f%%\n", grossProfitMargin*100)
	fmt.Printf("Net Profit Margin:%.0f%%\n", netProfitMargin*100)
	fmt.Printf("Working Capital Ratio:%.0f%%", workingCapitalRatio/liabilities*100)
}
