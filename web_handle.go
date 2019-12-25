package testtask1

import (
	"encoding/json"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
	"io/ioutil"
	"log"
	"net/http"
)

func (b *TestTask) HandleTransactionAction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(404)
		_, _ = w.Write([]byte(`Not found`))
		return
	}
	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(400)
		w.Write([]byte(`invalid header content-type, must be "application/json"`))
		return
	}
	if !InArray(r.Header.Get("Source-Type"), []string{string(SOURCE_TYPE_GAME), string(SOURCE_TYPE_PAYMENT), string(SOURCE_TYPE_SERVER)}) {
		w.WriteHeader(400)
		w.Write([]byte(`invalid header source-type`))
		return
	}
	//UserId, skip all authentication steps, we suppose all is ok
	userId, err := getUserIdFromRequest(r.Header.Get("auth"))
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(`invalid header auth, no user found with token`))
		return

	}
	var transactionRequest TransactionRequest
	response := TransactionResponse{Status: STATUS_FAIL}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responseJsonError(err, w)
	}
	err = json.Unmarshal(body, &transactionRequest)
	if err != nil {
		responseJsonError(err, w)
		return
	}
	//success json unmarshal
	if transactionRequest.Amount.LessThan(decimal.Zero) {
		err = fmt.Errorf("field 'amount' cannot be less than 0, amount: %s", transactionRequest.Amount)
		responseJsonError(err, w)
		return

	}
	//Add transaction
	err = b.AddTransaction(userId, transactionRequest.TransactionId, transactionRequest.State, transactionRequest.Amount)
	if err != nil {
		response.Error = err.Error()
		responseJsonError(err, w)
		return
	}
	balance, err := GetUserBalance(userId, b.db)
	if err != nil {
		responseJsonError(err, w)
		return
	}
	response.Error = ""
	response.Status = STATUS_SUCCESS
	response.Balance = balance
	responseBody, err := json.Marshal(&response)
	if err != nil {
		log.Println("cannot format response body, err: ", err)
		w.WriteHeader(500)
		_, err = w.Write([]byte(`Cannot format response body`))
		if err != nil {
			log.Println("cannot write response")
		}
		return
	}
	_, err = w.Write(responseBody)
	if err != nil {
		log.Println("cannot write response after marshal response")
	}
	return
}

//fake method
func getUserIdFromRequest(authenticationToken string) (userId uuid.UUID, err error) {
	userId, err = uuid.FromString(authenticationToken)
	if err != nil {
		return
	}
	return
}

func responseJsonError(err error, w http.ResponseWriter) {
	response := TransactionResponse{Status: STATUS_FAIL}
	response.Error = err.Error()
	responseBody, err := json.Marshal(&response)
	if err != nil {
		log.Println("cannot format response body, err: ", err)
		w.WriteHeader(500)
		_, err = w.Write([]byte(`Cannot format response body`))
		if err != nil {
			log.Println("cannot write response")
		}
		return
	}
	_, err = w.Write(responseBody)
	if err != nil {
		log.Println("cannot write response after marshal response")
	}
	return
}
