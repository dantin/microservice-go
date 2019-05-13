package dbclient

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/boltdb/bolt"

	"github.com/dantin/microservice-go/account/model"
)

// IBoltClient is an interface that interacts with BoltDB.
type IBoltClient interface {
	OpenBoltDB()
	QueryAccount(accountID string) (model.Account, error)
	Seed()
}

// BoltClient is the implementation of `IBoltClient`.
type BoltClient struct {
	boltDB *bolt.DB
}

// OpenBoltDB opens a connection to BoltDB.
func (bc *BoltClient) OpenBoltDB() {
	var err error
	bc.boltDB, err = bolt.Open("accounts.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
}

// Seed stars seeding accounts.
func (bc *BoltClient) Seed() {
	bc.initializeBucket()
	bc.seedAccounts()
}

// initializeBucket creates an `AccountBucket` in our BoltDB.  It will
// overwrite any existing bucket of the same name.
func (bc *BoltClient) initializeBucket() {
	bc.boltDB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte("AccountBucket"))
		if err != nil {
			return fmt.Errorf("create bucket failed: %s", err)
		}
		return nil
	})
}

// seedAccounts make-believe account objects into the `AccountBucket` bucket.
func (bc *BoltClient) seedAccounts() {
	total := 100
	for i := 0; i < total; i++ {
		// Generate a key 10000 or larger.
		key := strconv.Itoa(10000 + i)

		// Crate an instance of our Account struct.
		acc := model.Account{
			ID:   key,
			Name: "Person_" + strconv.Itoa(i),
		}

		// Serialize the struct to JSON.
		jsonBytes, _ := json.Marshal(acc)

		// Write the data into the `AccountBucket`.
		bc.boltDB.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("AccountBucket"))
			err := b.Put([]byte(key), jsonBytes)
			return err
		})
	}
	fmt.Printf("Seeded %v fake accounts...\n", total)
}

// QueryAccount find `Account` by `id`.
func (bc *BoltClient) QueryAccount(accountID string) (model.Account, error) {
	// Allocate an empty Account instance.
	account := model.Account{}

	// Read an object from the bucket using boltDB.view.
	err := bc.boltDB.View(func(tx *bolt.Tx) error {
		// Read the bucket from DB.
		b := tx.Bucket([]byte("AccountBucket"))

		// Read the value identified by accountID.
		accountBytes := b.Get([]byte(accountID))
		if accountBytes == nil {
			return fmt.Errorf("No account found for %s", accountID)
		}

		// Unmarshal the returned bytes into the account struct.
		json.Unmarshal(accountBytes, &account)

		return nil
	})

	if err != nil {
		return model.Account{}, nil
	}

	return account, nil
}
