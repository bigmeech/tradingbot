package adapters

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bigmeech/tradingbot/clients"
	"github.com/bigmeech/tradingbot/pkg/types"
)

// RequestFormatter formats requests for the REST API using orderType and side.
type RequestFormatter func(orderType types.OrderType, side types.OrderSide, tradingPair string, amount, price float64) (string, string, interface{}, error)

type RestExecutor struct {
	restClient    *clients.RestClient
	formatRequest RequestFormatter
}

// NewRestExecutor initializes a RestExecutor with a REST client and a request formatter.
func NewRestExecutor(restClient *clients.RestClient, formatRequest RequestFormatter) *RestExecutor {
	return &RestExecutor{
		restClient:    restClient,
		formatRequest: formatRequest,
	}
}

// ExecuteOrder prepares and sends a request to the exchange's REST API to place an order.
func (re *RestExecutor) ExecuteOrder(orderType types.OrderType, side types.OrderSide, tradingPair string, amount, price float64) error {
	// Format the request using orderType and side
	endpoint, method, body, err := re.formatRequest(orderType, side, tradingPair, amount, price)
	if err != nil {
		return fmt.Errorf("failed to format request: %w", err)
	}

	// Marshal body to JSON if it's a POST request
	var jsonBody []byte
	if method == "POST" {
		jsonBody, err = json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
	}

	// Make the request using the REST client
	resp, err := re.restClient.DoRequest(method, endpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to execute order: %w", err)
	}

	// Handle response here if needed, e.g., check status code or parse response body
	fmt.Printf("Order executed successfully: %v\n", resp)
	return nil
}
