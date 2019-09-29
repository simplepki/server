package ledger

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"	
)


type DBConnectionInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Engine   string `json:"engine"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
}

func getDB() (*gorm.DB, error) {
	log.Println("getting new aurora db connection")
	config := &aws.Config{}
	session := session.Must(session.NewSession(config))

	sm := secretsmanager.New(session)

	queryString := "mysql"

	query := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(queryString),
		VersionStage: aws.String("AWSCURRENT"),
	}

	log.Println("getting db config from secrets manager")
	output, err := sm.GetSecretValue(query)
	if err != nil {
		log.Fatal(err)
	}
	

	var dbConnectionInfo DBConnectionInfo
	err = json.Unmarshal([]byte(*output.SecretString), &dbConnectionInfo)
	
	if err != nil {
		log.Fatalln(err.Error())
		return nil, err
	}

	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		dbConnectionInfo.Username,
		dbConnectionInfo.Password,
		dbConnectionInfo.Host,
		dbConnectionInfo.Port,
		"ledger")
	
	db, err := gorm.Open("mysql", connectionString)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	log.Println("connection open")

	if !db.HasTable(&AuroraLedgerRecord{}){
		log.Println("creating aurora ledger table")
		db.CreateTable(&AuroraLedgerRecord{})
	}

	log.Println("auto-migrating aurora ledger table")
	db.AutoMigrate(&AuroraLedgerRecord{})

	return db, nil
}

type AWSAuroraLedger struct{}

type AuroraLedgerRecord struct {
	gorm.Model
	Name        string `gorm:"primary_key;size:255"`
	Account     string `gorm:"primary_key;size:255"`
	Certificate []byte
}

func (ledger AWSAuroraLedger) Publish(lrecord LedgerRecord) error {
	log.Printf("publishing cert %s for account %s\n", lrecord.Name, lrecord.Account)
	db, err := getDB()
	if err != nil {
		db.Close()
		return err
	}
	defer db.Close()

	var record AuroraLedgerRecord
	db.Where(&AuroraLedgerRecord{Name: lrecord.Name, Account: lrecord.Account}).Find(&record)

	record.Name = lrecord.Name
	record.Account = lrecord.Account
	record.Certificate = lrecord.Certificate

	db.Save(&record)
	log.Println("published")
	return nil
}
