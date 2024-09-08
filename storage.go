package main

import (
	"fmt"
	"math/rand"
	"path/filepath"
	"time"

	bolt "go.etcd.io/bbolt"
)

const DBName = "picshuffle.db"

// GetImageName  returns random image file name and update database for next usage
func GetImageName(execPath, folder string, jpgs []string) (string, bool) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	i := r.Intn(len(jpgs))
	rs := jpgs[i]
	fmt.Printf("Initially selected %s out of %d images\n", rs, len(jpgs))
	baseName := filepath.Join(execPath, DBName)
	db, err := bolt.Open(baseName, 0600, nil)
	if err != nil {
		fmt.Printf("Failed to open storage %s\n", baseName)
		return rs, true
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		_, errBckt := tx.CreateBucketIfNotExists([]byte(folder))
		return errBckt
	})

	var candidates []string
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(folder))
		for _, k := range jpgs {
			v := b.Get([]byte(k))
			if v == nil { // file was not displayed yet, it is candidate
				candidates = append(candidates, string(k))
				fmt.Printf("Candidate %s\n", k)
			}
		}
		return nil
	})

	if len(candidates) == 0 { // all files used. clear bucket for next usage
		err = db.Update(func(tx *bolt.Tx) error {
			tx.DeleteBucket([]byte(folder))
			_, err := tx.CreateBucket([]byte(folder))
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			fmt.Printf("%s\n", err)
			return rs, true
		}
	} else {
		i = r.Intn(len(candidates))
		rs = candidates[i]
		fmt.Printf("Selected %s out of %d\n", rs, len(candidates))
	}
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(folder))
		return b.Put([]byte(rs), []byte(""))
	})
	if err != nil {
		fmt.Printf("%s\n", err)
	}
	// Store current wall paper location/name into the log bucket
	t := time.Now()
	st := t.Format(time.RFC3339)
	tago := t.AddDate(-1, 0, 0)
	stago := tago.Format(time.RFC3339)
	db.Update(func(tx *bolt.Tx) error {
		b, errBckt := tx.CreateBucketIfNotExists([]byte("::LOG"))
		if errBckt != nil {
			return errBckt
		}
		logValue := "Folder:" + folder + ", Image:" + rs
		errBckt = b.Put([]byte(st), []byte(logValue))
		if errBckt != nil {
			return errBckt
		}
		c := b.Cursor()
		k, _ := c.First()      // can not be nill as far as Put was OK
		if string(k) < stago { // record is more than year old. Discard
			b.Delete(k)
		}
		return nil
	})
	return rs, true
}

// PrintLog prints name of files set as Wall Paper.
func PrintLog(execPath string) {
	baseName := filepath.Join(execPath, DBName)
	db, err := bolt.Open(baseName, 0600, nil)
	if err != nil {
		fmt.Printf("Failed to open storage %s\n", baseName)
		return
	}
	defer db.Close()
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("::LOG"))
		if b == nil {
			fmt.Printf("Nothing to print for storage: %s, backet %s\n", baseName, "::LOG")
			return nil
		}
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			fmt.Printf("%s: %s\n", k, v)
		}
		return nil
	})
}
