package testtask1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/labstack/gommon/random"
	"github.com/shopspring/decimal"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestTestTask_HandleTransactionAction(t *testing.T) {
	requestStruct := TransactionRequest{
		State:         STATE_LOST,
		Amount:        decimal.NewFromFloat(10.340009),
		TransactionId: random.String(20, random.Alphanumeric),
	}
	requestBody, _ := json.Marshal(requestStruct)
	request, _ := http.NewRequest(http.MethodPost, "http://localhost:8098/my_url", bytes.NewReader(requestBody))
	request.Header.Add("content-type", "application/json")
	request.Header.Add("source-type", SOURCE_TYPE_GAME)
	request.Header.Add("auth", "1c67d879-ba3b-48aa-b4fa-53b53b32d15d")
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		t.Errorf("bad request, err: %s", err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body), resp.Header)
}
