package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"crypto-wallet-service/config"

	"github.com/redis/go-redis/v9"
)


type PriceResponse map[string]map[string]float64


type CoinGeckoService struct {
	apiURL      string
	redisClient *redis.Client
}


func NewCoinGeckoService(redisClient *redis.Client) *CoinGeckoService {
	return &CoinGeckoService{
		apiURL:      config.AppConfig.CoinGecko.APIURL,
		redisClient: redisClient,
	}
}


func (s *CoinGeckoService) GetCryptoPrices(currencies []string) (PriceResponse, error) {
	ctx := context.Background()
	cacheKey := "crypto_prices"


	cached, err := s.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var data PriceResponse
		if err := json.Unmarshal([]byte(cached), &data); err == nil {
			return data, nil
		}
	}

	
	data, err := s.fetchFromAPI(currencies)
	if err != nil {
		return nil, err
	}


	jsonData, _ := json.Marshal(data)
	s.redisClient.Set(ctx, cacheKey, jsonData, config.GetCacheDuration())

	return data, nil
}


func (s *CoinGeckoService) GetPrice(currency string) (float64, error) {
	currency = strings.ToUpper(currency)

	
	if currency == "IDR" {
		return 1.0, nil
	}

	currencies := []string{"bitcoin", "ethereum", "tether"}
	prices, err := s.GetCryptoPrices(currencies)
	if err != nil {
		return 0, err
	}

	
	coinMap := map[string]string{
		"BTC":  "bitcoin",
		"ETH":  "ethereum",
		"USDT": "tether",
	}

	coinID, exists := coinMap[currency]
	if !exists {
		return 0, fmt.Errorf("unsupported currency: %s", currency)
	}

	if priceData, ok := prices[coinID]; ok {
		if idrPrice, ok := priceData["idr"]; ok {
			return idrPrice, nil
		}
	}

	return 0, fmt.Errorf("price not found for currency: %s", currency)
}


func (s *CoinGeckoService) fetchFromAPI(currencies []string) (PriceResponse, error) {
	
	ids := strings.Join(currencies, ",")
	url := fmt.Sprintf("%s/simple/price?ids=%s&vs_currencies=idr", s.apiURL, ids)

	
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch prices from CoinGecko: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("CoinGecko API returned status code: %d", resp.StatusCode)
	}

	var data PriceResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return data, nil
}


func (s *CoinGeckoService) ClearCache() error {
	ctx := context.Background()
	return s.redisClient.Del(ctx, "crypto_prices").Err()
}
