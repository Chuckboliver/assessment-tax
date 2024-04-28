package tax

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strconv"
)

type parser interface {
	parseCalculationRequest(reader io.Reader) ([]calculationRequest, error)
}

var _ parser = (*csvParser)(nil)

type csvParser struct{}

func newCSVParser() parser {
	return &csvParser{}
}

func (c *csvParser) parseCalculationRequest(reader io.Reader) ([]calculationRequest, error) {
	csvReader := csv.NewReader(reader)

	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	if (len(records) == 0) || (len(records[0]) == 0) {
		return nil, errors.New("empty csv file")
	}

	parsedResult := make([]calculationRequest, 0)

	headerRow := records[0]
	if err := validateHeaderRow(headerRow); err != nil {
		return nil, err
	}

	for i := 1; i < len(records); i++ {

		calculationRequest := calculationRequest{
			Allowances: make([]Allowance, 0),
		}

		for j, col := range records[i] {
			columnName := headerRow[j]
			switch columnName {
			case "totalIncome":
				value, err := strconv.ParseFloat(col, 64)
				if err != nil {
					return nil, fmt.Errorf("failed to parse totalIncome: %w", err)
				}

				calculationRequest.TotalIncome = value
			case "wht":
				value, err := strconv.ParseFloat(col, 64)
				if err != nil {
					return nil, fmt.Errorf("failed to parse wht: %w", err)
				}

				calculationRequest.Wht = value
			case "donation":
				value, err := strconv.ParseFloat(col, 64)
				if err != nil {
					return nil, fmt.Errorf("failed to parse donation: %w", err)
				}

				calculationRequest.Allowances = append(calculationRequest.Allowances, Allowance{
					AllowanceType: AllowanceDonation,
					Amount:        value,
				})
			}
		}

		parsedResult = append(parsedResult, calculationRequest)
	}

	return parsedResult, nil
}

func validateHeaderRow(headers []string) error {
	mustHaveColumns := map[string]struct{}{
		"totalIncome": {},
		"wht":         {},
		"donation":    {},
	}

	for _, header := range headers {
		if _, ok := mustHaveColumns[header]; !ok {
			return fmt.Errorf("unknown column: %s", header)
		}
	}

	return nil
}
