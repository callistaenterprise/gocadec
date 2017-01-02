package service

import (
        "github.com/boltdb/bolt"
        "fmt"
        "strconv"
        "github.com/callistaenterprise/gocadec/accountservice/model"
        "encoding/json"
        "log"
        "github.com/opentracing/opentracing-go"
        ct "github.com/eriklupander/cloudtoolkit"
)

var boltDB *bolt.DB

func QueryAccount(accountId string, span opentracing.Span) (model.Account, error) {
        childSpan := ct.Tracer.StartSpan("QueryAccount", opentracing.ChildOf(span.Context()))
        defer childSpan.Finish()
        account := model.Account{}
        err  := boltDB.View(func(tx *bolt.Tx) error {
                b := tx.Bucket([]byte("AccountBucket"))
                accountBytes := b.Get([]byte(accountId))
                if accountBytes == nil {
                    return fmt.Errorf("No account found for " + accountId)
                }
                json.Unmarshal(accountBytes, &account)
                return nil
        })
        if err != nil {
             return model.Account{}, err
        }
        return account, nil
}

func SeedAccounts() {
        boltDB.Update(func(tx *bolt.Tx) error {
                _, err := tx.CreateBucket([]byte("AccountBucket"))
                if err != nil {
                        return fmt.Errorf("create bucket: %s", err)
                }
                return nil
        })
        for i := 0; i < 100; i++ {
                key := strconv.Itoa(10000 + i)
                acc := model.Account{
                        Id: key,
                        Name: "Person_" + strconv.Itoa(i),
                }
                json, _ := json.Marshal(acc)
                boltDB.Update(func(tx *bolt.Tx) error {
                        b := tx.Bucket([]byte("AccountBucket"))
                        err := b.Put([]byte(key), json)
                        return err
                })

        }
        ct.Log.Println("Seeded 100 fake accounts...")
}

func OpenBoltDb() {
        var err error
        boltDB, err = bolt.Open("accounts.db", 0600, nil)
        if err != nil {
                log.Fatal(err)
        }
}
