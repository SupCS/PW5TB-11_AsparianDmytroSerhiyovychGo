package handlers

import (
	"html/template"
	"math"
	"net/http"
	"strconv"
)

// Структура параметрів надійності
type ReliabilityParameters struct {
	FailureRate     float64
	AvgRepairTime   int
	Frequency       float64
	AvgRecoveryTime int
}

// Дані для розрахунків
var reliabilityData = map[string]ReliabilityParameters{
	"T-110 kV":                     {0.015, 100, 1.0, 43},
	"T-35 kV":                      {0.02, 80, 1.0, 28},
	"T-10 kV (Cable Network)":      {0.005, 60, 0.5, 10},
	"T-10 kV (Overhead Network)":   {0.05, 60, 0.5, 10},
	"B-110 kV (Gas-Insulated)":     {0.01, 30, 0.1, 30},
	"B-10 kV (Oil)":                {0.02, 15, 0.33, 15},
	"B-10 kV (Vacuum)":             {0.05, 15, 0.33, 15},
	"Busbars 10 kV per Connection": {0.03, 2, 0.33, 15},
	"AV-0.38 kV":                   {0.05, 20, 1.0, 15},
	"ED 6,10 kV":                   {0.1, 50, 0.5, 0},
	"ED 0.38 kV":                   {0.1, 50, 0.5, 0},
	"PL-110 kV":                    {0.007, 10, 0.167, 35},
	"PL-35 kV":                     {0.02, 8, 0.167, 35},
	"PL-10 kV":                     {0.02, 10, 0.167, 35},
	"CL-10 kV (Trench)":            {0.03, 44, 1.0, 9},
	"CL-10 kV (Cable Channel)":     {0.005, 18, 1.0, 9},
}

// Функція для обробки запиту на розрахунок надійності
func ReliabilityHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Зчитуємо дані введення
		totalFailureRate := 0.0
		weightedRecoveryTime := 0.0

		for key := range reliabilityData {
			quantity, _ := strconv.Atoi(r.FormValue(key))
			params := reliabilityData[key]

			if quantity > 0 {
				totalFailureRate += float64(quantity) * params.FailureRate
				weightedRecoveryTime += float64(quantity) * params.FailureRate * float64(params.AvgRepairTime)
			}
		}

		// Обчислення показників
		averageRecovery := weightedRecoveryTime / totalFailureRate
		accidentalDowntime := averageRecovery * totalFailureRate / 8760
		plannedDowntime := 1.2 * 43 / 8760
		dualFailureRate := 2 * totalFailureRate * (accidentalDowntime + plannedDowntime)
		finalRate := dualFailureRate + 0.02

		// Відображаємо результати у шаблоні
		tmpl, _ := template.ParseFiles("templates/reliability.html")
		tmpl.Execute(w, map[string]string{
			"TotalFailureRate":   roundFloat(totalFailureRate, 4),
			"AverageRecovery":    roundFloat(averageRecovery, 4),
			"AccidentalDowntime": roundFloat(accidentalDowntime, 4),
			"PlannedDowntime":    roundFloat(plannedDowntime, 4),
			"DualFailureRate":    roundFloat(dualFailureRate, 4),
			"FinalFailureRate":   roundFloat(finalRate, 4),
		})
		return
	}

	// Відображення сторінки з формою
	tmpl, _ := template.ParseFiles("templates/reliability.html")
	tmpl.Execute(w, nil)
}

// Функція округлення числа до n знаків після коми
func roundFloat(value float64, precision int) string {
	return strconv.FormatFloat(math.Round(value*math.Pow10(precision))/math.Pow10(precision), 'f', precision, 64)
}
