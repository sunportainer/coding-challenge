package service

import (
	"encoding/json"
	"github.com/pkg/errors"
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

func CalculateStatus(file string) error {
	if _, err := os.Stat(file); err != nil {
		if os.IsNotExist(err) {
			return errors.Wrap(err, "message: failed to locate data file")
		}
		return err
	}
	jsonFile, err := os.Open(file)
	if err != nil {
		return errors.Wrap(err, "message: failed to open data file")
	}
	defer jsonFile.Close()
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return errors.Wrap(err, "message: failed to read data file")
	}
	var books BookKeeping
	err = json.Unmarshal(byteValue, &books)
	if err != nil {
		return errors.Wrap(err, "message: failed to parse data file")
	}
	var revenue, expense, grossProfitMarginTotal float64
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
			if record.ValueType == "credit" {
				liabilities = liabilities + record.TotalValue
			}
			if record.ValueType == "debit" {
				liabilities = liabilities - record.TotalValue
			}
		}
	}
	p := message.NewPrinter(language.English)
	p.Printf("Revenue: $%.0f\n", revenue)
	p.Printf("Expense: $%.0f\n", expense)
	if revenue == 0 {
		p.Printf("Gross Profit Margin: 0%")
		p.Printf("Net Profit Margin: 0%")
	} else {
		grossProfitMargin := grossProfitMarginTotal / revenue
		p.Printf("Gross Profit Margin: %.0f%%\n", grossProfitMargin*100)
		netProfitMargin := (revenue - expense) / revenue
		p.Printf("Net Profit Margin: %.0f%%\n", netProfitMargin*100)
	}
	if liabilities == 0 {
		p.Printf("Working Capital Ratio: 0%")
	} else {
		p.Printf("Working Capital Ratio: %.0f%%\n", workingCapitalRatio/liabilities*100)
	}

	return nil
}
