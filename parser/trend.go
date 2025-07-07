package parser

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

func GetCountryTrend(file, indicator, country string) map[string]float64 {
	f, _ := os.Open(file)
	defer f.Close()
	r := csv.NewReader(f)
	rows, _ := r.ReadAll()

	idx := map[string]int{}
	for i, h := range rows[0] {
		idx[h] = i //idx[indicatorName]=0, idx[countryName]=1, idx[year]=2, idx[value]=3
	}

	trend := map[string]float64{}
	for _, row := range rows[1:] {
		if row[idx["Country Name"]] == country && row[idx["Indicator Name"]] == indicator {
			val, _ := strconv.ParseFloat(row[idx["Value"]], 64)
			trend[row[idx["Year"]]] = val
		}
	}
	return trend //trend[2018]=1.2, trend[2019]=1.5, ...
}

func BuildTrendPrompt(country, indicator string, trend map[string]float64, base string) string {
	return fmt.Sprintf(`
	%s
	
	Use the following indicator to provide a response:
	
	Indicator: %s
	
	Data for %s:
	
	Trend: %s
	
	Based on the indicator and trend, generate a helpful and empathetic response including stats and data insights that we are passing as Trend in less than 200 words.

	Refer to the data as retrieved by you (the assistant), not provided by the user.

	Make the data look structured in a paragraph when giving output!

	Avoid excessive indentation or Markdown nesting.

	Use empathy and understanding in your response, acknowledging the user's situation and providing insights based on the data.

	`, base, indicator, country, formatTrend(trend))
}

func formatTrend(trend map[string]float64) string {
	var years []string
	for y := range trend {
		years = append(years, y)
	}
	sort.Strings(years)

	var b strings.Builder
	for _, y := range years {
		b.WriteString(fmt.Sprintf("%s: %.1f\n", y, trend[y]))
	}
	return b.String()
}
