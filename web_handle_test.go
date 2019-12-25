package testtask1

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"testing"
)

func TestTestTask_HandleTransactionAction(t *testing.T) {
	for i := 0; i < 10; i++ {
		v4, _ := uuid.NewV4()
		fmt.Printf(`"%s",`, v4)
	}
}
