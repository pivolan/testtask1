package testtask1

import (
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
	"time"
)

type TransactionBet struct {
	ID          string          `gorm:"type:string;primary_key;"`
	CreatedAt   time.Time       `sql:"index"`
	CancelledAt time.Time       `sql:"index"`
	Amount      decimal.Decimal `sql:"type:decimal(20,8);"`
	State       StateType
	UserID      uuid.UUID
	User        UserBalance
}
type UserBalance struct {
	Base
	Balance decimal.Decimal `sql:"type:decimal(20,8);"`
}

type Base struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

func (u *Base) BeforeCreate(scope *gorm.Scope) error {
	uuid4, err := uuid.NewV4()
	if err != nil {
		return err
	}
	return scope.SetColumn("ID", uuid4)
}
