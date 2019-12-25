package testtask1

import (
	"github.com/labstack/gommon/random"
	uuid "github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
	"log"
	"testing"
	"time"
)

const DEFAULT_TEST_DSN = `host=localhost port=5432 user=postgres dbname=testtask1 sslmode=disable`

func TestGetUserBalance(t *testing.T) {
	b := TestTask{}
	err := b.ConnectDb(DEFAULT_TEST_DSN)
	if err != nil {
		log.Fatalln(err)
	}
	//create temp user
	balance, err := decimal.NewFromString("376.9006")
	user := UserBalance{
		Balance: balance,
	}
	b.db.Create(&user)
	newUser := UserBalance{}
	b.db.Find(&newUser, "id=?", user.ID)
	if user.Balance != balance {
		t.Fail()
		return
	}
	//get balance on empty
	balance, err = GetUserBalance(user.ID, b.db)
	if err != nil {
		t.Error("no balance, err:", err)
	}
	if !balance.Equal(decimal.Zero) {
		t.Errorf("not null balance for empty userBalance, balance: %s\n", balance)
	}
	//add transaction and update balance
	err = b.AddTransaction(user.ID, "tid1", STATE_WIN, decimal.NewFromFloat(10.15))
	if err != nil {
		t.Error(err)
	}
	//get balance
	balance, err = GetUserBalance(user.ID, b.db)
	if err != nil {
		t.Errorf("no balance, err: %s\n", err)
	}
	if !balance.Equal(decimal.NewFromFloat(10.15)) {
		t.Errorf("user balance not equal to 10.15, balance: %s\n", balance)
	}
	b.db.Delete(&TransactionBet{}, "user_id=?", user.ID)
	b.db.Delete(&user)
}
func TestCancelTransactions(t *testing.T) {
	b := TestTask{}
	err := b.ConnectDb(DEFAULT_TEST_DSN)
	if err != nil {
		log.Fatalln(err)
	}
	v4, _ := uuid.NewV4()
	b.AddTransaction(v4, "tid", STATE_WIN, decimal.NewFromInt(5))
	var transaction TransactionBet
	b.db.Find(&transaction, "id=?", "tid")
	CancelTransactions(v4, []TransactionBet{transaction}, b.db)
	b.db.Find(&transaction, "id=?", "tid")
	if transaction.CancelledAt == nil {
		t.Errorf("not changed tranaction cancelled_at")
	}
	now := time.Now()
	transaction.CancelledAt = &now
	b.db.Save(&transaction)
	b.db.Find(&transaction, "id=?", "tid")
	if transaction.CancelledAt == nil {
		t.Errorf("not changed tranaction cancelled_at by save method")
	}
	b.db.Delete(&transaction)
}
func TestGetLast10OddTransactionUser(t *testing.T) {
	b := TestTask{}
	err := b.ConnectDb(DEFAULT_TEST_DSN)
	if err != nil {
		log.Fatalln(err)
	}
	//create temp user
	balance, err := decimal.NewFromString("376.9006")
	user := UserBalance{
		Balance: balance,
	}
	b.db.Create(&user)
	//add 25 transactions and check ordering
	sum := int64(0)
	for i := int64(0); i < 20; i++ {
		sum += i
		err := b.AddTransaction(user.ID, random.String(20, random.Alphanumeric), STATE_WIN, decimal.NewFromInt(i))
		if err != nil {
			t.Errorf("error on add transaction on step %d, err: %s\n", i, err)
		}
	}
	for i := int64(5); i < 10; i++ {
		sum -= i
		err := b.AddTransaction(user.ID, random.String(20, random.Alphanumeric), STATE_LOST, decimal.NewFromInt(i))
		if err != nil {
			t.Errorf("error on add transaction on step %d, err: %s\n", i, err)
		}
	}
	//check all transactions
	var transactions []TransactionBet
	err = b.db.Find(&transactions, "user_id=? AND cancelled_at IS NULL", user.ID).Order("order_uuid").Error
	if err != nil {
		t.Error(err)
		return
	}
	if len(transactions) != 25 {
		t.Errorf("transactions count is not 20, count: %d", len(transactions))
	}
	//check balance
	balance, err = GetUserBalance(user.ID, b.db)
	if err != nil {
		t.Error(err)
		return
	}
	if !balance.Equal(decimal.NewFromInt(sum)) {
		t.Errorf("sum balance not equal to %d, balance: %s", sum, balance)
	}
	lastTransactions, err := GetLast10OddTransactionUser(user.ID, b.db)
	if err != nil {
		t.Error(err)
		return
	}
	if len(lastTransactions) != 10 {
		t.Error("not 10 last odd transactions, ", len(lastTransactions))
	}
	amounts := []int64{-9, -7, -5, 18, 16, 14, 12, 10, 8, 6}
	expectedSum := sum
	for i, transaction := range lastTransactions {
		expectedAmount := decimal.NewFromInt(amounts[i])
		expectedSum -= expectedAmount.IntPart()
		if !expectedAmount.Equal(transaction.Amount) {
			t.Errorf("amount transaction not equal to expected %s amount, transaction.amount = %s\n", expectedAmount, transaction.Amount)
		}
	}

	err = b.Cancel10LastOddUserTransactions(user.ID)
	if err != nil {
		t.Fatalf("error on cancel10LastOdd %s", err)
		return
	}
	//check all transactions after cancel
	err = b.db.Find(&transactions, "user_id=? AND cancelled_at IS NULL", user.ID).Order("order_uuid").Error
	if err != nil {
		t.Error(err)
		return
	}
	if len(transactions) != 15 {
		t.Errorf("transactions count is not 15, count: %d", len(transactions))
	}
	//check balance after cancel
	balance, err = GetUserBalance(user.ID, b.db)
	if err != nil {
		t.Error(err)
		return
	}
	if !balance.Equal(decimal.NewFromInt(expectedSum)) {
		t.Errorf("sum balance not equal to %d, balance: %s", expectedSum, balance)
	}

	b.db.Delete(&TransactionBet{}, "user_id=?", user.ID)
	b.db.Delete(&UserBalance{}, "id=?", user.ID)
}
