package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

// Match phrases to indicator
func matchIndicator(input string) string {
	input = strings.ToLower(input)

	switch {
	// --- EDUCATION ---
	case strings.Contains(input, "tech") ||
		strings.Contains(input, "ict") ||
		strings.Contains(input, "computer science") ||
		strings.Contains(input, "it field") ||
		strings.Contains(input, "programming") ||
		strings.Contains(input, "technology") ||
		strings.Contains(input, "engineering"):
		return "Female share of graduates in Information and Communication Technologies programmes, tertiary (%)"

	case strings.Contains(input, "dropout") || strings.Contains(input, "out of school"):
		return "Children out of school, primary, female"

	case strings.Contains(input, "literacy") || strings.Contains(input, "read") || strings.Contains(input, "write"):
		return "Literacy rate, adult total (% of people ages 15 and above)"

	case strings.Contains(input, "college") || strings.Contains(input, "university") || strings.Contains(input, "graduate"):
		return "Gross graduation ratio, tertiary, female (%)"

	case strings.Contains(input, "stem") || strings.Contains(input, "science") || strings.Contains(input, "technology") || strings.Contains(input, "engineering"):
		return "Female share of graduates from Science, Technology, Engineering and Mathematics (STEM) programmes, tertiary (%)"

	case strings.Contains(input, "gender parity") || strings.Contains(input, "equal") || (strings.Contains(input, "boys") && strings.Contains(input, "girls")):
		return "School enrollment, primary and secondary (gross), gender parity index (GPI)"

	// --- HEALTH ---
	case strings.Contains(input, "health") || strings.Contains(input, "sick"):
		return "Maternal mortality ratio (modeled estimate, per 100,000 live births)"

	case strings.Contains(input, "mortality") && strings.Contains(input, "adolescent"):
		return "Mortality rate, adolescent female (per 1,000 female adolescents)"

	case strings.Contains(input, "maternal") || strings.Contains(input, "birth") || strings.Contains(input, "pregnancy"):
		return "Maternal mortality ratio (modeled estimate, per 100,000 live births)"

	case strings.Contains(input, "contraception") || strings.Contains(input, "birth control"):
		return "Contraceptive prevalence, any methods (% of women ages 15-49)"

	case strings.Contains(input, "fertility"):
		return "Adolescent fertility rate (births per 1,000 women ages 15-19)"

	case strings.Contains(input, "life expectancy") || strings.Contains(input, "live longer"):
		return "Life expectancy at birth, female (years)"

	case strings.Contains(input, "ict") || strings.Contains(input, "tech"):
		return "Female share of graduates from Information and Communication Technologies programmes, tertiary (%)"

	// --- DEFAULT ---
	default:
		return ""
	}
}

// Parse CSV and return trend
func getCountryTrend(filename, indicatorName, country string) map[string]float64 {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	header := records[0]
	idxCountry, idxIndicator, idxYear, idxValue := -1, -1, -1, -1

	for i, h := range header {
		switch h {
		case "Country Name":
			idxCountry = i
		case "Indicator Name":
			idxIndicator = i
		case "Year":
			idxYear = i
		case "Value":
			idxValue = i
		}
	}

	trend := make(map[string]float64)

	for _, row := range records[1:] {
		if row[idxCountry] != country || row[idxIndicator] != indicatorName {
			continue
		}
		year := row[idxYear]
		val, err := strconv.ParseFloat(row[idxValue], 64)
		if err != nil {
			continue
		}
		trend[year] = val
	}

	return trend
}

// Build natural prompt for Gemini
func buildPrompt(country, indicator string, trend map[string]float64, userContext string) string {
	fmt.Println(trend)
	safeIndicator := strings.ReplaceAll(indicator, "%", "%%")
	return fmt.Sprintf(`
%s

Use the following indicator to provide a response:
Indicator: %s

Data for %s:

Trend: %s

Based on the indicator and trend, generate a helpful and empathetic response including stats and data insights that we are passing as Trend.
`, userContext, safeIndicator, country, formatTrendData(trend))
}

// Sort & format map as string
func formatTrendData(trend map[string]float64) string {
	years := make([]string, 0, len(trend))
	for year := range trend {
		years = append(years, year)
	}
	sort.Strings(years)

	var sb strings.Builder
	for _, year := range years {
		sb.WriteString(fmt.Sprintf("%s: %.1f\n", year, trend[year]))
	}
	return sb.String()
}
