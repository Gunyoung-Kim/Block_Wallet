package db

import (
	"github.com/Gunyoung-Kim/blockchain/utils"
	"github.com/boltdb/bolt"
)

var db *bolt.DB // varaible for singleton pattern of *blot.DB

const (
	dbName = "blockchain.db" // DB Name

	dataBucket   = "data"   // Bucket name for checkPoint of blockChain
	blocksBucket = "blocks" // Bucket name for blocks

	checkPoint = "checkPoint" // Key for dataBucket, all data for dataBucket use this key
)

//DB return *bolt.DB which is designed by singleton pattern
// create bucker for dataBucket and blocksBucket if not exist
func DB() *bolt.DB {
	if db == nil {
		dbPointer, err := bolt.Open(dbName, 0600, nil)
		utils.HandleError(err)
		db = dbPointer
		db.Update(func(t *bolt.Tx) error {
			_, err = t.CreateBucketIfNotExists([]byte(dataBucket))
			utils.HandleError(err)
			_, err = t.CreateBucketIfNotExists([]byte(blocksBucket))
			return err
		})
		utils.HandleError(err)
	}

	return db
}

//Close database
func Close() {
	DB().Close()
}

// ------------------- functions for dataBucket --------------

//CheckPoint read checkPoint from DB(dataBucket) and return slice of byte
//use transaction for read-only
func CheckPoint() []byte {
	var data []byte
	DB().View(func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(dataBucket))
		data = bucket.Get([]byte(checkPoint))
		return nil
	})
	return data
}

//SaveCheckPoint save checkPoint of blockChain in dataBucket
func SaveCheckPoint(data []byte) {
	err := DB().Update(func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(dataBucket))
		err := bucket.Put([]byte(checkPoint), data)
		return err
	})
	utils.HandleError(err)
}

// ------------------- functions for blocksBucket --------------

//Block read a Block from DB(blocksBucket) and return slice of byte
//use transaction for read-only
func Block(hash string) []byte {
	var data []byte
	DB().View(func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(blocksBucket))
		data = bucket.Get([]byte(hash))
		return nil
	})
	return data
}

//SaveBlock save a block in blocksBucket
func SaveBlock(hash string, data []byte) {
	err := DB().Update(func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(blocksBucket))
		err := bucket.Put([]byte(hash), data)
		return err
	})
	utils.HandleError(err)
}
