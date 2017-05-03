package main

import (
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"time"
)

type sUser struct {
	Id       int
	Password string
}

type toDo struct {
	Id     int
	Text   string
	Status int
	Date   string
	User   string
}

func getUserNameByToken(token string) string {
	db, err := bolt.Open("database/bolt.db", 0644, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	user := ""

	db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		bucket := tx.Bucket([]byte("tokens"))

		cursor := bucket.Cursor()

		for key, value := cursor.First(); key != nil; key, value = cursor.Next() {
			if string(value) == token {
				user = string(key)
			}
		}

		return nil
	})

	return user
}

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
		bucket := tx.Bucket([]byte("users"))
		if bucket == nil {
			return fmt.Errorf("Bucket %q not found!", []byte("users"))
		}

		json.Unmarshal(bucket.Get(key), &userJson)

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	return userJson
}

func createUser(login string, password string) int {
	db, err := bolt.Open("database/bolt.db", 0644, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	key := []byte(login)

	var userExist int = 0
	// retrieve the data
	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("users"))
		if bucket == nil {
			return fmt.Errorf("Bucket %q not found!", []byte("users"))
		}

		if bucket.Get(key) != nil {
			userExist = 1
		}

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	if userExist == 1 {
		return 0
	}

	var user sUser
	user.Password = password

	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("users"))
		if err != nil {
			return err
		}

		id, err := bucket.NextSequence()
		user.Id = int(id)

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

	return 1
}

func setUser(login string, id int, password string) {
	db, err := bolt.Open("database/bolt.db", 0644, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	key := []byte(login)
	var user sUser
	user.Id = id
	user.Password = password

	// store some data
	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("Users"))
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
		bucket := tx.Bucket([]byte("tokens"))
		if bucket == nil {
			return fmt.Errorf("Bucket %q not found!", []byte("tokens"))
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
		bucket, err := tx.CreateBucketIfNotExists([]byte("tokens"))
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

func addToDo(login string, text string) {
	db, err := bolt.Open("database/bolt.db", 0644, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var toDo toDo

	// store some data
	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("todos"))
		if err != nil {
			return err
		}

		id, err := bucket.NextSequence()

		current_time := time.Now().Local()

		toDo.Id = int(id)
		toDo.Status = 0
		toDo.Text = text
		toDo.Date = current_time.Format("02.01.2006")
		toDo.User = login

		encoded, err := json.Marshal(toDo)

		err = bucket.Put([]byte(string(id)), encoded)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}

func getToDoS(login string) string {
	db, err := bolt.Open("database/bolt.db", 0644, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ToDoS := make([]toDo, 0)

	// retrieve the data
	db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		bucket := tx.Bucket([]byte("todos"))

		cursor := bucket.Cursor()

		var valueJson toDo

		for key, value := cursor.First(); key != nil; key, value = cursor.Next() {
			json.Unmarshal(value, &valueJson)
			tmpToDoS := make([]toDo, 0)
			if valueJson.User == login {
				tmpToDoS = append(tmpToDoS, valueJson)
				ToDoS = append(tmpToDoS, ToDoS...)
			}
		}

		return nil
	})

	result, err := json.Marshal(ToDoS)

	if err != nil {
		log.Fatal(err)
	}

	return string(result)
}

func editToDo(id int, text string, status int) {
	db, err := bolt.Open("database/bolt.db", 0644, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var toDo toDo

	var changeKey []byte

	db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		bucket := tx.Bucket([]byte("todos"))

		cursor := bucket.Cursor()

		for key, value := cursor.First(); key != nil; key, value = cursor.Next() {
			if string(key) == string(id) {
				json.Unmarshal(value, &toDo)
				changeKey = key
			}
		}

		return nil
	})

	// store some data
	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("todos"))
		if err != nil {
			return err
		}

		id, err := bucket.NextSequence()

		toDo.Status = status
		if text != "" {
			toDo.Text = text
		}

		encoded, err := json.Marshal(toDo)

		err = bucket.Put(changeKey, encoded)

		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}

func getToDoById(id int) string {
	db, err := bolt.Open("database/bolt.db", 0644, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var toDo string

	db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		bucket := tx.Bucket([]byte("todos"))

		cursor := bucket.Cursor()

		for key, value := cursor.First(); key != nil; key, value = cursor.Next() {
			if string(key) == string(id) {
				toDo = string(value)
			}
		}

		return nil
	})

	return toDo
}

func deleteToDo(id int) {
	db, err := bolt.Open("database/bolt.db", 0644, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var toDo toDo

	var changeKey []byte

	db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		bucket := tx.Bucket([]byte("todos"))

		cursor := bucket.Cursor()

		for key, value := cursor.First(); key != nil; key, value = cursor.Next() {
			if string(key) == string(id) {
				json.Unmarshal(value, &toDo)
				changeKey = key
			}
		}

		return nil
	})

	// store some data
	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("todos"))
		if err != nil {
			return err
		}

		err = bucket.Delete(changeKey)

		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}
