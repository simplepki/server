package ledger

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
	"net/url"
	"strings"

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
	db.DB().SetConnMaxLifetime(300 * time.Second)

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
	Name        string `gorm:"size:255"`
	Account     string `gorm:"size:255"`
	Certificate string `gorm:"size:5464"`
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

func (ledger AWSAuroraLedger) GetChainForRecord(account string, certUrl url.URL) ([]LedgerRecord, error) {
	log.Printf("getting chain for cert %#v for account %s\n", certUrl, account)
	db, err := getDB()
	if err != nil {
		db.Close()
		return []LedgerRecord{}, err
	}
	defer db.Close()

	var path string
	if strings.HasSuffix(certUrl.Path,"/") {
		path = fmt.Sprintf("%s%s",certUrl.Host, certUrl.Path[0:len(certUrl.Path)-1])
	} else {
		path = fmt.Sprintf("%s%s",certUrl.Host, certUrl.Path)
	}

	paths := strings.Split(path, "/")
	searchPaths := make([]string,len(paths))
	for idx, p := range paths {
		switch idx {
		case 0:
			searchPaths[idx] = p
		default:
			searchPaths[idx] = fmt.Sprintf("%s/%s", searchPaths[idx-1], p)
		}
	}

	log.Printf("looking up certs for account %s using paths %#v\n", account, searchPaths)
	for _, path := range searchPaths {
		searchPath := fmt.Sprintf("spiffe://%s", path)
		var foundRecord AuroraLedgerRecord
		db.Where(&AuroraLedgerRecord{Account: account, Name: searchPath}).First(&foundRecord)
		log.Printf("account: %v, path:%v, record: %#v\n", certUrl.Host, path, foundRecord)
	}


	return []LedgerRecord{}, nil
}