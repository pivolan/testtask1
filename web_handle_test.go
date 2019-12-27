package testtask1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/labstack/gommon/random"
	"github.com/shopspring/decimal"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"testing"
)

func TestTestTask_HandleTransactionAction_InvalidHeader(t *testing.T) {
	transactionRequest := TransactionRequest{
		State:         STATE_LOST,
		Amount:        decimal.NewFromFloat(10.340009),
		TransactionId: random.String(20, random.Alphanumeric),
	}
	//invalid content-type
	request, _ := GenerateRequest("invalid", SOURCE_TYPE_GAME, "1c67d879-ba3b-48aa-b4fa-53b53b32d153", transactionRequest)
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		t.Errorf("bad request, err: %s", err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if !strings.Contains(string(body), `invalid header content-type, must be "application/json"`) {
		t.Errorf("response invalid, resp: %s", string(body))
	}
	//invalid source-type
	request, _ = GenerateRequest("application/json", "invalid_source", "1c67d879-ba3b-48aa-b4fa-53b53b32d153", transactionRequest)
	resp, err = client.Do(request)
	if err != nil {
		t.Errorf("bad request, err: %s", err)
	}
	body, err = ioutil.ReadAll(resp.Body)
	if !strings.Contains(string(body), `invalid header source-type`) {
		t.Errorf("response invalid, resp: %s", string(body))
	}
	//invalid auth uuid
	request, _ = GenerateRequest("application/json", SOURCE_TYPE_GAME, "1c67d879-ba3b-48aa-b4fa-53b53b32d153_", transactionRequest)
	resp, err = client.Do(request)
	if err != nil {
		t.Errorf("bad request, err: %s", err)
	}
	body, err = ioutil.ReadAll(resp.Body)
	if !strings.Contains(string(body), `invalid header auth, token is not uuid, token: 1c67d879-ba3b-48aa-b4fa-53b53b32d153_`) {
		t.Errorf("response invalid, resp: %s", string(body))
	}
	//invalid auth user not found
	request, _ = GenerateRequest("application/json", SOURCE_TYPE_GAME, "1c67d879-ba3b-48aa-b4fa-53b53b32d153", transactionRequest)
	resp, err = client.Do(request)
	if err != nil {
		t.Errorf("bad request, err: %s", err)
	}
	body, err = ioutil.ReadAll(resp.Body)
	if !strings.Contains(string(body), `invalid header auth, cannot find user with token: 1c67d879-ba3b-48aa-b4fa-53b53b32d153, err: record not found`) {
		t.Errorf("response invalid, resp: %s", string(body))
	}
}
func TestTestTask_HandleTransactionAction_NegativeBalance(t *testing.T) {
	transactionRequest := TransactionRequest{
		State:         STATE_LOST,
		Amount:        decimal.NewFromFloat(10.340009),
		TransactionId: random.String(20, random.Alphanumeric),
	}

	request, _ := GenerateRequest("application/json", SOURCE_TYPE_GAME, "1c67d879-ba3b-48aa-b4fa-53b53b32d15d", transactionRequest)
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		t.Errorf("bad request, err: %s", err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	response := TransactionResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		t.Errorf("error on unmarshall, err: %s, resp: %s", err, string(body))
	}
	if !strings.Contains(response.Error, `balance cannot be less than zero after transaction`) {
		t.Errorf("response must be with error, resp: %s", string(body))
	}
	fmt.Println(response)
}
func TestTestTask_HandleTransactionAction(t *testing.T) {
	userId := "efe90aa5-f8ac-42d5-a372-1876351afa86"

	var resp *http.Response
	var err error
	var sum decimal.Decimal
	//clear user data for test
	request, err := http.NewRequest(http.MethodPost, "http://localhost:8098/clear_user", nil)
	if err != nil {
		return
	}
	request.Header.Add("auth", userId)
	client := http.Client{}
	resp, err = client.Do(request)
	if err != nil {
		t.Errorf("error generate request, err: %s", err)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if !strings.Contains(string(body), "Success clear user data") {
		t.Errorf("cannot clear data for user: %s, err: %s", userId, string(body))
		return
	}
	for i := 0; i < 100; i++ {
		transactionRequest := TransactionRequest{
			State:         STATE_WIN,
			Amount:        decimal.NewFromFloat(10.340009),
			TransactionId: random.String(20, random.Alphanumeric),
		}

		request, _ := GenerateRequest("application/json", SOURCE_TYPE_GAME, userId, transactionRequest)
		client := http.Client{}
		resp, err = client.Do(request)
		if err != nil {
			t.Errorf("error generate request, err: %s", err)
			continue
		}
		body, err := ioutil.ReadAll(resp.Body)
		response := TransactionResponse{}
		err = json.Unmarshal(body, &response)
		if err != nil {
			t.Errorf("error on unmarshall, err: %s, resp: %s", err, string(body))
			continue
		}
		if response.Status != "success" {
			t.Errorf("response must be success, resp: %s", string(body))
			continue
		}
		sum = sum.Add(transactionRequest.Amount)
		if !response.Balance.Equal(sum) {
			t.Errorf("user_balance invalid, transaction: %s, response: %s", transactionRequest.Amount, response.Balance)
		}
	}
	fmt.Println("test parallel lost")
	wg := sync.WaitGroup{}
	for i := 0; i < 500; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			transactionRequest := TransactionRequest{
				State:         STATE_LOST,
				Amount:        decimal.NewFromFloat(10.340009),
				TransactionId: random.String(20, random.Alphanumeric),
			}

			request, _ := GenerateRequest("application/json", SOURCE_TYPE_GAME, userId, transactionRequest)
			client := http.Client{}
			resp, err = client.Do(request)
			if err != nil {
				t.Errorf("error generate request, err: %s", err)
				return
			}
			body, err := ioutil.ReadAll(resp.Body)
			response := TransactionResponse{}
			err = json.Unmarshal(body, &response)
			if err != nil {
				t.Errorf("error on unmarshall, err: %s, resp: %s", err, string(body))
				return
			}
			if response.Status != "success" {
				t.Errorf("response must be success, resp: %s", string(body))
				return
			}
			sum = sum.Add(transactionRequest.Amount)
			if !response.Balance.Equal(sum) {
				t.Errorf("user_balance invalid, transaction: %s, response: %s", transactionRequest.Amount, response.Balance)
			}
		}()
	}
	wg.Wait()
}

func GenerateRequest(contentType string, sourceType string, auth string, transactionRequest TransactionRequest) (request *http.Request, err error) {
	requestBody, _ := json.Marshal(transactionRequest)
	request, err = http.NewRequest(http.MethodPost, "http://localhost:8098/my_url", bytes.NewReader(requestBody))
	if err != nil {
		return
	}
	request.Header.Add("content-type", contentType)
	request.Header.Add("source-type", sourceType)
	request.Header.Add("auth", auth)
	return
}
