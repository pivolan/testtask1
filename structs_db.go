package testtask1

import (
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
	"time"
)

type TransactionBet struct {
	ID          string          `gorm:"primary_key;"`
	OrderId     uint64          `gorm:"type:bigserial;AUTO_INCREMENT"`
	CreatedAt   time.Time       `gorm:"index"`
	CancelledAt *time.Time      `gorm:"index"`
	Amount      decimal.Decimal `gorm:"type:decimal(20,8);"`
	State       StateType
	UserID      uuid.UUID
	User        UserBalance
}
type UserBalance struct {
	Base
	Balance decimal.Decimal `gorm:"type:decimal(20,8);"`
}

type Base struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

func (u *Base) BeforeCreate(scope *gorm.Scope) error {
	if uuid.Equal(u.ID, uuid.UUID{}) {
		uuid4, err := uuid.NewV4()
		if err != nil {
			return err
		}
		return scope.SetColumn("ID", uuid4)
	}
	return nil
}
func (UserBalance) TableName() string {
	return "user_balance"
}
func (TransactionBet) TableName() string {
	return "transaction_bet"
}
