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

func (ledger AWSAuroraLedger) GetChainForRecord(certUrl url.URL) ([]LedgerRecord, error) {
	log.Printf("getting chain for cert %s for account %s\n", certUrl.Path, certUrl.Host)
	db, err := getDB()
	if err != nil {
		db.Close()
		return []LedgerRecord{}, err
	}
	defer db.Close()

	var paths []string
	if strings.HasPrefix(certUrl.Path, "/") {
		paths = strings.Split(certUrl.Path[1:len(certUrl.Path)], "/")
	} else {
		paths = strings.Split(certUrl.Path,"/")
	}

	searchPaths := make([]string, len(paths))
	for idx, path := range paths {
		insertIdx := (len(paths)-1) - idx
		switch idx {
		case 0:
			searchPaths[insertIdx] = path
		default:
			searchPaths[insertIdx] = fmt.Sprintf("%s/%s", searchPaths[insertIdx+1], path)
		}
	}

	log.Printf("looking up certs for account %s using paths %#v\n", certUrl.Host, searchPaths)
	for _, path := range searchPaths {
		searchPath := fmt.Sprintf("spiffe://%s", path)
		var foundRecord AuroraLedgerRecord
		db.Where(&AuroraLedgerRecord{Account: certUrl.Host, Name: searchPath}).First(&foundRecord)
		log.Printf("account: %v, path:%v, record: %#v\n", certUrl.Host, path, foundRecord)
	}


	return []LedgerRecord{}, nil
}