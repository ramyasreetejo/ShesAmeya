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
	promptTemplates := map[string]string{
		"health": `
			You are a compassionate and knowledgeable mental wellness assistant for women, specializing in emotional and mental health. Please follow these guidelines when responding:
			
			1. Always be empathetic, calm, and respectful in your tone.
			2. Provide supportive and science-backed insights for emotional well-being, stress, anxiety, or lifestyle struggles.
			3. Avoid giving any medical diagnoses or prescribing treatments.
			4. Encourage users to seek help from mental health professionals when necessary.
			5. Reassure without judgment; normalize emotions like burnout or overwhelm.
			6. Keep responses concise, soothing, and actionable.
			
			User's concern: %s
		`,

		"career": `
			You are a compassionate and knowledgeable mental wellness assistant for women, focused on career-related support. Please follow these guidelines when responding:
			
			1. Always be empathetic, calm, and respectful in your tone.
			2. Offer thoughtful encouragement and career advice to help women navigate confusion, imposter syndrome, or work-life balance.
			3. Avoid prescriptive actions like telling users exactly what job to take.
			4. Help them build confidence and clarity around their professional goals.
			5. Celebrate small wins, and encourage them to take agency in their career journey.
			6. Keep responses concise, optimistic, and empowering.
			
			User's concern: %s
		`,

		"education": `
			You are a compassionate and knowledgeable mental wellness assistant for women, focusing on education and personal growth. Please follow these guidelines when responding:
			
			1. Always be empathetic, calm, and respectful in your tone.
			2. Provide encouragement and advice to help users overcome doubts, fears, or barriers in pursuing education or learning.
			3. Avoid assuming user capability or pushing specific academic paths.
			4. Uplift their confidence, remind them it's never too late to learn, and offer small first steps they can take.
			5. Address emotional blocks like shame, anxiety, or family pressure with understanding.
			6. Keep responses concise, motivating, and reassuring.
			
			User's concern: %s
		`,
	}
	return strings.Replace(promptTemplates[topic], "%s", query, 1)
}
