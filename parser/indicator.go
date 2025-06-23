package parser

import "strings"

func MatchIndicator(input string) string {
	input = strings.ToLower(input)
	switch {
	// --- EDUCATION & CAREER ---
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

	// --- DEFAULT ---
	default:
		return ""
	}
}

func ExtractTopicAndQuery(input string) (string, string) {
	switch {
	case strings.HasPrefix(input, "/health "):
		return "health", strings.TrimPrefix(input, "/health ")
	case strings.HasPrefix(input, "/career "):
		return "career", strings.TrimPrefix(input, "/career ")
	case strings.HasPrefix(input, "/education "):
		return "education", strings.TrimPrefix(input, "/education ")
	default:
		return "general", input
	}
}

func PromptForTopic(topic, query string) string {
	prompts := map[string]string{
		"health":    "You are a supportive assistant for women's health. User's concern: %s",
		"career":    "You are a career coach for women. User's concern: %s",
		"education": "You are an educational motivator for women. User's concern: %s",
	}
	return strings.Replace(prompts[topic], "%s", query, 1)
}
