package main

import (
	"github.com/boltdb/bolt"
	"log"
	"fmt"
	"encoding/json"
)
type sUser struct {
	Id  int
	Password string
}

var users = []byte("users")
var tokens = []byte("tokens")
//var usersPass = []byte("usersPass")

func getUser(login string) sUser {

	db, err := bolt.Open("database/bolt.db", 0644, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	key := []byte(login)
	var userJson sUser
	// retrieve the data
	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(users)
		if bucket == nil {
			return fmt.Errorf("Bucket %q not found!", users)
		}

		json.Unmarshal(bucket.Get(key), &userJson)

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	return userJson
}

func createUser(login string, password string) {
	db, err := bolt.Open("database/bolt.db", 0644, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	key := []byte(login)
	var user sUser
	user.Password=password

	// store some data
	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(users)
		if err != nil {
			return err
		}

		id, err := bucket.NextSequence()
		user.Id= int(id)

		encoded, err := json.Marshal(user)

		err = bucket.Put(key, encoded)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}

func setUser(login string, id int, password string) {
	db, err := bolt.Open("database/bolt.db", 0644, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	key := []byte(login)
	var user sUser
	user.Id=id
	user.Password=password

	// store some data
	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(users)
		if err != nil {
			return err
		}

		encoded, err := json.Marshal(user)

		err = bucket.Put(key, encoded)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}

func getToken(login string) string {
	db, err := bolt.Open("database/bolt.db", 0644, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	key := []byte(login)
	var token string
	// retrieve the data
	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(tokens)
		if bucket == nil {
			return fmt.Errorf("Bucket %q not found!", tokens)
		}

		token = string(bucket.Get(key))

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	return token
}

func createToken(login string) {
	db, err := bolt.Open("database/bolt.db", 0644, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	key := []byte(login)

	// store some data
	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(tokens)
		if err != nil {
			return err
		}

		token, err := GenerateRandomString(32)
		if err != nil {
			// Serve an appropriately vague error to the
			// user, but log the details internally.
		}

		err = bucket.Put(key, []byte(token))
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}