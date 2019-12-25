package testtask1

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
	"time"
)

func (b *TestTask) AddTransaction(userId uuid.UUID, transactionId string, state StateType, amount decimal.Decimal) (err error) {
	tx := b.db.Model(&UserBalance{}).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return err
	}
	balance, err := GetUserBalance(userId, tx)
	if err != nil {
		err = fmt.Errorf("error on get balance, err: %s", err)
		return
	}
	if state == STATE_LOST {
		amount = amount.Neg()
	}
	if balance.Add(amount).LessThan(decimal.Zero) {
		err = fmt.Errorf("balance cannot be less than zero after transaction, transaction amount: %s, current balance: %s", amount, balance)
		return
	}
	if err = tx.Create(&TransactionBet{
		ID:        transactionId,
		CreatedAt: time.Now(),
		Amount:    amount,
		State:     state,
		UserID:    userId,
	}).Error; err != nil {
		err = fmt.Errorf("error on create transaction, err: %s", err)
		tx.Rollback()
		return
	}
	if err = tx.Model(&UserBalance{}).Where("id = ?", userId).Update("balance", amount.Add(balance)).Error; err != nil {
		err = fmt.Errorf("error on update user balance, transaction id: %s, err: %s", transactionId, err)
		tx.Rollback()
		return
	}
	tx.Commit()
	return
}
func (b *TestTask) Cancel10LastOddUserTransactions(userId uuid.UUID) (err error) {
	tx := b.db.Model(&UserBalance{}).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return err
	}
	balance, err := GetUserBalance(userId, tx)
	if err != nil {
		err = fmt.Errorf("error on get balance, err: %s", err)
		return
	}
	transactions, err := GetLast10OddTransactionUser(userId, tx)
	if err != nil {
		err = fmt.Errorf("error on get last 10 odd transactions, err: %s", err)
		return
	}
	balanceResult := decimal.Zero.Add(balance)
	for _, transaction := range transactions {
		balanceResult = balanceResult.Add(transaction.Amount)
	}
	if balance.LessThan(decimal.Zero) {
		err = fmt.Errorf("error, after cancel 10 transactions user balance will less than zero, current user balance: %s, after user balance: %s", balance, balanceResult)
		return
	}
	err = CancelTransactions(userId, transactions, tx)
	if err != nil {
		tx.Rollback()
		return
	}
	afterUserBalance, err := GetUserBalance(userId, tx)
	if err != nil {
		tx.Rollback()
		err = fmt.Errorf("cannot check balance after cancel transaction, rollback this commit, userId: %s, err: %s", userId, err)
		return
	}
	if afterUserBalance.LessThan(decimal.Zero) {
		err = fmt.Errorf("after cancel transaction balance are less than zero, rollback all, user balance before: %s, after: %s", balance, afterUserBalance)
		tx.Rollback()
		return
	}
	tx.Model(&UserBalance{}).Where("id=?", userId).Update("balance", afterUserBalance)
	tx.Commit()
	return
}
func GetUserBalance(userId uuid.UUID, tx *gorm.DB) (balance decimal.Decimal, err error) {
	row := tx.Model(&TransactionBet{}).Raw("SELECT COALESCE(SUM(amount),0) AS total FROM transaction_bet WHERE user_id=? AND cancelled_at IS NULL", userId).Row()
	err = row.Scan(&balance)
	if err != nil {
		err = fmt.Errorf("error on raw query, err: %s", err)
		return
	}
	return
}
func GetLast10OddTransactionUser(userId uuid.UUID, tx *gorm.DB) (transactions []TransactionBet, err error) {
	tempTransactions := []TransactionBet{}
	if err = tx.Order("created_at DESC").Limit(20).Find(&tempTransactions, "user_id = ?", userId).Error; err != nil {
		err = fmt.Errorf("error on find 20 last transaction for user: %s, err: %s", userId, err)
		return
	}
	for i, transaction := range tempTransactions {
		if i%2 == 0 {
			transactions = append(transactions, transaction)
		}
	}
	return
}
func CancelTransactions(userId uuid.UUID, transactions []TransactionBet, tx *gorm.DB) (err error) {
	for i, transaction := range transactions {
		if err = tx.Model(&TransactionBet{}).Where("id=?", transaction.ID).Update("cancelled_at", time.Now()).Error; err != nil {
			err = fmt.Errorf("error on cancel transactions %d for user: %s, err: %s", i, userId, err)
			return
		}
	}
	return
}
