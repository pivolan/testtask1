package testtask1

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"net/http"
	"time"
)

func (b *TestTask) Init(dbConnect string, port string) (err error) {
	b.db, err = gorm.Open("postgres", dbConnect)
	if err != nil {
		err = fmt.Errorf("cannot connect to database postgres, err: %s", err)
		return
	}
	b.db.AutoMigrate(&TransactionBet{})
	b.db.AutoMigrate(&UserBalance{})
	b.StartCronTasks()
	http.HandleFunc("/my_url", b.HandleTransactionAction)

	log.Fatal(http.ListenAndServe("127.0.0.1:"+port, nil))

	return
}
func (b *TestTask) StartCronTasks() {
	go func() {
		for {
			<-time.After(10 * time.Minute)
			users := []UserBalance{}
			b.db.Find(&users)
			for _, user := range users {
				err := b.Cancel10LastOddUserTransactions(user.ID)
				if err != nil {
					log.Println(err)
				}
			}
		}
	}()
}
