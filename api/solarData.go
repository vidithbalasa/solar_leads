package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "net/url"
    "os"
)

func getLatLong(address string) (map[string]float64, error) {
    apiKey := os.Getenv("GOOGLE_MAPS_API_KEY")
    encodedAddress := url.QueryEscape(address)
    url := "https://maps.googleapis.com/maps/api/geocode/json?address=" + encodedAddress + "&key=" + apiKey

    resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var result map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, err
    }

    if result["status"] == "OK" {
        results := result["results"].([]interface{})
        if len(results) > 0 {
            location := results[0].(map[string]interface{})["geometry"].(map[string]interface{})["location"].(map[string]interface{})
            lat := location["lat"].(float64)
            lng := location["lng"].(float64)
            return map[string]float64{"latitude": lat, "longitude": lng}, nil
        }
    }

    return nil, fmt.Errorf("no valid results found")
}

func getSolarData(address string) (interface{}, error) {
    location, err := getLatLong(address)
    if err != nil {
        return nil, fmt.Errorf("failed to get location: %v", err)
    }

    apiKey := os.Getenv("GOOGLE_MAPS_API_KEY")
    baseURL := "https://solar.googleapis.com/v1/buildingInsights:findClosest"
    params := url.Values{}
    params.Add("key", apiKey)
    params.Add("location.latitude", fmt.Sprintf("%f", location["latitude"]))
    params.Add("location.longitude", fmt.Sprintf("%f", location["longitude"]))
    params.Add("requiredQuality", "HIGH")

    reqURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())
    resp, err := http.Get(reqURL)
    if err != nil {
        return nil, fmt.Errorf("failed to get solar data: %v", err)
    }
    defer resp.Body.Close()

    var result map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, fmt.Errorf("failed to decode solar data: %v", err)
    }

    // Extracting specific solar data fields
    if solarPotential, ok := result["solarPotential"].(map[string]interface{}); ok {
        maxPanelCount := solarPotential["maxArrayPanelsCount"]
        maxArea := solarPotential["maxArrayAreaMeters2"]
        maxSunshineHours := solarPotential["maxSunshineHoursPerYear"]
        carbonOffset := solarPotential["carbonOffsetFactorKgPerMwh"]

        solarData := map[string]interface{}{
            "Max Panel Count":              maxPanelCount,
            "Max Area (m²)":                maxArea,
            "Max Sunshine Hours per Year": maxSunshineHours,
            "Carbon Offset Factor (kg/MWh)": carbonOffset,
        }
        return solarData, nil
    }

    return nil, fmt.Errorf("solar data not found in the response")
}
