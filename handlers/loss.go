package handlers

import (
	"html/template"
	"net/http"
	"strconv"
)

// Функція обробки запиту для калькулятора втрат
func LossHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Отримуємо значення з форми
		omega, _ := strconv.ParseFloat(r.FormValue("omega"), 64)
		tb, _ := strconv.ParseFloat(r.FormValue("tb"), 64)
		pNom, _ := strconv.ParseFloat(r.FormValue("pNom"), 64)
		tm, _ := strconv.ParseFloat(r.FormValue("tm"), 64)
		kp, _ := strconv.ParseFloat(r.FormValue("kp"), 64)
		zPer0, _ := strconv.ParseFloat(r.FormValue("zPer0"), 64)
		zPlan, _ := strconv.ParseFloat(r.FormValue("zPlan"), 64)

		// Обчислення втрат
		mwAvar := omega * pNom * tb * tm
		mwPlan := kp * pNom * tm
		totalLosses := zPer0 + (mwAvar * zPer0) + (mwPlan * zPlan)

		// Відображаємо результати у шаблоні
		tmpl, _ := template.ParseFiles("templates/loss.html")
		tmpl.Execute(w, map[string]string{
			"MWnedAvar":   roundFloat(mwAvar, 4),
			"MWnedPlan":   roundFloat(mwPlan, 4),
			"TotalLosses": roundFloat(totalLosses, 4),
		})
		return
	}

	// Відображення сторінки з формою
	tmpl, _ := template.ParseFiles("templates/loss.html")
	tmpl.Execute(w, nil)
}
