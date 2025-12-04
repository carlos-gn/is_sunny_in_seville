package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// WeatherData represents the current weather response
type WeatherData struct {
	Main struct {
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		Humidity  int     `json:"humidity"`
		Pressure  int     `json:"pressure"`
	} `json:"main"`
	Weather []struct {
		Main        string `json:"main"`
		Description string `json:"description"`
	} `json:"weather"`
	Wind struct {
		Speed float64 `json:"speed"`
	} `json:"wind"`
	Clouds struct {
		All int `json:"all"`
	} `json:"clouds"`
}

type (
	Input  struct{}
	Output struct {
		Message string `json:"message" jsonschema:"Message from the weather service"`
	}
)

const (
	sevilleLat     = 37.3886
	sevilleLon     = -5.9823
	openWeatherAPI = "https://api.openweathermap.org/data/2.5"
)

var apiKey string

func init() {
	apiKey = os.Getenv("OPENWEATHER_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENWEATHER_API_KEY environment variable not set")
	}
}

// getSevilleWeather fetches current weather data for Seville
func getSevilleWeather(ctx context.Context) (*WeatherData, error) {
	url := fmt.Sprintf("%s/weather?lat=%.4f&lon=%.4f&units=metric&appid=%s",
		openWeatherAPI, sevilleLat, sevilleLon, apiKey)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch weather: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var data WeatherData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &data, nil
}

// isSunnyInSeville checks if it's currently sunny in Seville
func isSunnyInSeville(ctx context.Context) (string, error) {
	data, err := getSevilleWeather(ctx)
	if err != nil {
		return "", err
	}

	if len(data.Weather) == 0 {
		return "", fmt.Errorf("no weather data available")
	}

	weather := data.Weather[0]
	isSunny := weather.Main == "Clear" || weather.Main == "Sunny"

	if isSunny {
		return fmt.Sprintf("Yes, it's sunny in Seville! ☀️ (%.1f°C)", data.Main.Temp), nil
	}
	return fmt.Sprintf("No, it's not sunny in Seville. Weather: %s (%.1f°C)", weather.Description, data.Main.Temp), nil
}

func main() {
	// Create a new MCP server
	s := mcp.NewServer(
		&mcp.Implementation{
			Name:    "seville-weather",
			Version: "1.0.0",
		},
		nil,
	)

	// Add the "is_sunny" tool
	sunnyTool := &mcp.Tool{
		Name:        "is_sunny",
		Description: "Check if it's currently sunny in Seville, Spain",
	}

	mcp.AddTool(s, sunnyTool, func(ctx context.Context, req *mcp.CallToolRequest, input Input) (*mcp.CallToolResult, Output, error) {
		result, err := isSunnyInSeville(ctx)
		if err != nil {
			return nil, Output{}, fmt.Errorf("failed to check weather: %w", err)
		}
		return nil, Output{Message: result}, nil
	})

	// Run the server over stdin/stdout
	if err := s.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
