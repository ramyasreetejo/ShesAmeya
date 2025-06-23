package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ramyasreetejo/ShesAmeya/model"
	"github.com/ramyasreetejo/ShesAmeya/parser"
	"github.com/ramyasreetejo/ShesAmeya/util"
)

type Request struct {
	Message string `json:"message"`
	IP      string `json:"ip"`
}

type Response struct {
	Reply string `json:"reply"`
}

func ChatHandler(client *model.GeminiClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req Request
		_ = json.NewDecoder(r.Body).Decode(&req)

		country := util.GetCountryFromIP(req.IP)
		topic, query := parser.ExtractTopicAndQuery(req.Message)
		prompt := parser.PromptForTopic(topic, query)

		file := fmt.Sprintf("world_bank_data/%s.csv", topic)
		if ind := parser.MatchIndicator(req.Message); ind != "" {
			if trend := parser.GetCountryTrend(file, ind, country); len(trend) > 0 {
				prompt = parser.BuildTrendPrompt(country, ind, trend, prompt)
			}
		}

		reply, err := client.Generate(prompt)
		if err != nil {
			http.Error(w, "AI error: "+err.Error(), 500)
			return
		}

		_ = json.NewEncoder(w).Encode(Response{Reply: reply})
	}
}
