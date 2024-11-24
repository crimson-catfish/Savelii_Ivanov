package operproc

import (
	"flag"
	"fmt"
	"math"
	"os"
	"strconv"
	"time"
)

type CompanyOperations struct {
	Name                 string `json:"company"`
	ValidOperationsCount int    `json:"valid_operations_count"`
	Balance              int    `json:"balance"`
	InvalidOperationIDs  []any  `json:"invalid_operations"`
}

func GetFilePath(envVar string) (string, error) {
	var file string

	flag.StringVar(&file, "file", "", "File to validate")
	flag.Parse()
	if file != "" {
		return file, nil
	}

	file = os.Getenv(envVar)
	if file != "" {
		return file, nil
	}

	fmt.Print("Enter file to validate: ")
	_, err := fmt.Scanln(&file)
	if err != nil {
		return "", err
	}

	return file, nil
}

func ProcessOperations(operations []map[string]any) ([]CompanyOperations, error) {
	companies := map[string]CompanyOperations{}

	for _, op := range operations {
		name := parseStringField(op, "company")
		if name == "" {
			continue
		}

		if _, exists := companies[name]; !exists {
			companies[name] = CompanyOperations{Name: name, InvalidOperationIDs: []any{}}
		}

		if !processOperation(op, companies, name) {
			continue
		}
	}

	return compileCompanies(companies), nil
}

func parseStringField(data map[string]any, key string) string {
	if value, ok := data[key].(string); ok {
		return value
	}
	return ""
}

func parseEmbeddedField(data, embedded map[string]any, key string) any {
	if value, exists := data[key]; exists {
		return value
	}
	return embedded[key]
}

func processOperation(op map[string]any, companies map[string]CompanyOperations, name string) bool {
	embedded := extractEmbedded(op, "operation")
	createdAt := parseEmbeddedField(op, embedded, "created_at")
	if !isValidDate(createdAt) {
		return false
	}

	id := parseEmbeddedField(op, embedded, "id")
	if !isValidID(id) {
		return false
	}

	value, validValue := parseValue(op, embedded)
	if !validValue {
		WriteInvalidOperation(companies, name, id)
		return false
	}

	opType := parseStringField(op, "type")
	if opType == "" {
		opType = parseStringField(embedded, "type")
	}
	processOperationType(companies, name, opType, value, id)
	return true
}

func extractEmbedded(data map[string]any, key string) map[string]any {
	if embedded, ok := data[key].(map[string]any); ok {
		return embedded
	}
	return nil
}

func isValidDate(value any) bool {
	strValue, ok := value.(string)
	if !ok {
		return false
	}
	_, err := time.Parse(time.RFC3339, strValue)
	return err == nil
}

func isValidID(id any) bool {
	switch v := id.(type) {
	case string:
		return true
	case float64:
		return v == math.Trunc(v)
	default:
		return false
	}
}

func parseValue(data, embedded map[string]any) (int, bool) {
	value := parseEmbeddedField(data, embedded, "value")
	switch v := value.(type) {
	case string:
		parsed, err := strconv.Atoi(v)
		if err != nil {
			return 0, false
		}
		return parsed, true
	case float64:
		if v == math.Trunc(v) {
			return int(v), true
		}
	}
	return 0, false
}

func processOperationType(companies map[string]CompanyOperations, name, opType string, value int, id any) {
	switch opType {
	case "income", "+":
		WriteValidOperation(companies, name, value)
	case "outcome", "-":
		WriteValidOperation(companies, name, -value)
	default:
		WriteInvalidOperation(companies, name, id)
	}
}

func compileCompanies(companies map[string]CompanyOperations) []CompanyOperations {
	result := make([]CompanyOperations, 0, len(companies))
	for _, company := range companies {
		result = append(result, company)
	}
	return result
}

func WriteInvalidOperation(companies map[string]CompanyOperations, name string, id any) {
	company := companies[name]
	company.InvalidOperationIDs = append(company.InvalidOperationIDs, id)
	companies[name] = company
}

func WriteValidOperation(companies map[string]CompanyOperations, name string, value int) {
	company := companies[name]
	company.ValidOperationsCount++
	company.Balance += value
	companies[name] = company
}
