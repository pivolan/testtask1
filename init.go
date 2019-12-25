package testtask1

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"net/http"
	"time"
)

func (b *TestTask) Init(dsn string, port string) (err error) {
	err = b.ConnectDb(dsn)
	if err != nil {
		err = fmt.Errorf("cannot connect to database postgres, err: %s", err)
		return
	}
	b.db.DropTable(&TransactionBet{})
	log.Println("connect to db success")
	b.MigrateDatabase()
	log.Println("migrations done")
	b.StartCronTasks()
	log.Println("periodically tasks started")
	b.StartWebListen(port)
	return
}
func (b *TestTask) ConnectDb(dsn string) (err error) {
	b.db, err = gorm.Open("postgres", dsn)
	if err != nil {
		err = fmt.Errorf("cannot connect to database postgres, err: %s", err)
		return
	}
	return
}
func (b *TestTask) MigrateDatabase() {
	b.db.AutoMigrate(&TransactionBet{})
	b.db.AutoMigrate(&UserBalance{})
}
func (b *TestTask) StartWebListen(port string) {
	http.HandleFunc("/my_url", b.HandleTransactionAction)

	log.Fatal(http.ListenAndServe("127.0.0.1:"+port, nil))
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
