package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/biter777/countries"
	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

var promptTemplates = map[string]string{
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

func extractTopicAndQuery(input string) (string, string) {
	if strings.HasPrefix(input, "/health ") {
		return "health", strings.TrimPrefix(input, "/health ")
	} else if strings.HasPrefix(input, "/career ") {
		return "career", strings.TrimPrefix(input, "/career ")
	} else if strings.HasPrefix(input, "/education ") {
		return "education", strings.TrimPrefix(input, "/education ")
	}
	return "general", input
}

func getCountryFromIP(ip string) string {
	if ip == "::1" || ip == "127.0.0.1" || ip == "" {
		// Fallback to India Cloudflare node for local testing
		ip = "103.21.244.1"
	}
	reqURL := fmt.Sprintf("https://ipinfo.io/%s/country", ip)
	resp, err := http.Get(reqURL)
	if err != nil || resp.StatusCode != 200 {
		return "Unknown"
	}
	defer resp.Body.Close()
	country, _ := io.ReadAll(resp.Body)
	code := strings.TrimSpace(string(country))
	country_name := countries.ByName(code)
	return country_name.Info().Name
}

func extractIPFromRequest(r *http.Request) string {
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		return strings.Split(forwarded, ",")[0]
	}
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

func startRestAPI(model *genai.GenerativeModel, restPort string) {
	http.HandleFunc("/api/chat", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Message string `json:"message"`
			IP      string `json:"ip"` // optional IP sent from frontend
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		ip := req.IP
		if ip == "" {
			ip = extractIPFromRequest(r)
		}
		country := getCountryFromIP(ip)

		fmt.Printf("User IP: %s | Country: %s\n", ip, country)

		topic, query := extractTopicAndQuery(req.Message)
		template := promptTemplates[topic]
		mentalWellnessSystemPrompt := fmt.Sprintf(template, query)

		fileName := fmt.Sprintf("world_bank_data/%s.csv", topic)

		indicator := matchIndicator(req.Message)
		if indicator != "" {
			trend := getCountryTrend(fileName, indicator, country)
			if len(trend) != 0 {
				mentalWellnessSystemPrompt = buildPrompt(country, indicator, trend, mentalWellnessSystemPrompt)
			}
		}

		resp, err := model.GenerateContent(
			r.Context(),
			genai.Text(mentalWellnessSystemPrompt),
		)
		if err != nil {
			http.Error(w, "AI error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		reply := ""
		for _, cand := range resp.Candidates {
			if len(cand.Content.Parts) > 0 {
				if part, ok := cand.Content.Parts[0].(genai.Text); ok {
					reply = string(part)
					break
				}
			}
		}

		json.NewEncoder(w).Encode(map[string]string{"reply": reply})
	})

	// Serve static frontend (HTML) if applicable
	http.Handle("/", http.FileServer(http.Dir("./static")))

	log.Println("REST server running at http://localhost:" + restPort)
	log.Fatal(http.ListenAndServe(":"+restPort, nil))
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Failed to load .env")
	}

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Fatal("GEMINI_API_KEY not set")
	}
	restPort := os.Getenv("REST_PORT")
	if restPort == "" {
		restPort = "8080"
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalf("Failed to create Gemini client: %v", err)
	}
	model := client.GenerativeModel("gemini-1.5-flash")

	startRestAPI(model, restPort)
}
