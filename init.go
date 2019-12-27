package testtask1

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	uuid "github.com/satori/go.uuid"
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
	log.Println("connect to db success")
	b.MigrateDatabase()
	log.Println("migrations done")
	b.Fixtures()
	log.Println("fixtures done")
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
	b.db.DB().SetConnMaxLifetime(10 * time.Second)
	b.db.DB().SetMaxOpenConns(1000)
	b.db.DB().SetMaxIdleConns(100)
	return
}
func (b *TestTask) MigrateDatabase() {
	b.db.AutoMigrate(&TransactionBet{})
	b.db.AutoMigrate(&UserBalance{})
}
func (b *TestTask) Fixtures() {
	userUuidList := []string{"1c67d879-ba3b-48aa-b4fa-53b53b32d15d", "efe90aa5-f8ac-42d5-a372-1876351afa86", "d4d3d3d4-2fc0-4fc6-a796-bf857ef6ae0e", "f1ae2bc5-f5a8-44ab-96f9-64d4982b9c97", "9be78564-75c5-48ea-864f-c8434121007f", "1695f173-d51f-4ff9-8304-5e93dfa69ba0", "08b41b1c-0123-4b35-84ca-7b5ebbc4a57c", "db46eca6-0295-4514-9822-b1ff480e2546", "6d09abfc-a953-4545-aad7-efeb3ec78eaa", "4854a4b8-a991-472e-9f41-10cc68609c4d"}
	for _, userId := range userUuidList {
		uuidV4, _ := uuid.FromString(userId)
		userBalance := UserBalance{}
		userBalance.ID = uuidV4
		b.db.Create(&userBalance)
	}
}
func (b *TestTask) StartWebListen(port string) {
	http.HandleFunc("/my_url", b.HandleTransactionAction)
	http.HandleFunc("/clear_user", b.HandleUserClearAction)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
func (b *TestTask) StartCronTasks() {
	go func() {
		for {
			time.Sleep(10 * time.Minute)
			users := []UserBalance{}
			b.db.Find(&users)
			log.Printf("start task, found %d users to update transactions", len(users))
			for _, user := range users {
				log.Printf("start delete last 10 odd transactions for user: %s\n", user.ID)
				err := b.Cancel10LastOddUserTransactions(user.ID)
				if err != nil {
					log.Println(err)
				}
			}
		}
	}()
}
