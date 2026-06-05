package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fickleDude/gophemart/internal/config"
	"github.com/fickleDude/gophemart/internal/model"
)

func GetOrderAccrual(number string, client http.Client) (*model.Order, error) {
	cfg := config.GetConfig()
	baseURL := fmt.Sprintf("%s/api/orders", cfg.AccrualSystenAddress())
	url := fmt.Sprintf("%s/%s", baseURL, number)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	var order model.Order
	if err := json.NewDecoder(response.Body).Decode(&order); err != nil {
		return nil, err
	}
	defer response.Body.Close()
	return &order, nil

}
