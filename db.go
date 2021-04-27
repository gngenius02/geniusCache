package main

import (
	"encoding/hex"
	"fmt"
	"log"

	badger "github.com/dgraph-io/badger/v3"
)

// func getDataFromFile(path string, db *badger.DB) {
// 	file, _ := os.Open(path)
// 	defer wg.Done()
// 	defer file.Close()

// 	wb := db.NewWriteBatch()
// 	defer wb.Cancel()

// 	scanner := bufio.NewScanner(file)
// 	for scanner.Scan() {
// 		line := scanner.Text()
// 		if !strings.Contains(line, ",") {
// 			continue
// 		}

// 		str := strings.Split(line, ",")
// 		seedString, hashString := str[0], str[1]
// 		hashBytes, _ := hex.DecodeString(hashString)

// 		wb.Set(hashBytes, []byte(seedString))
// 	}
// 	fmt.Println(path, "finished.")
// 	wb.Flush()

// }

// func getFileList(path string, db *badger.DB) {
// 	dir, err := os.ReadDir(path)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	for _, file := range dir {
// 		if file.IsDir() {
// 			continue
// 		}
// 		fmt.Println("db is currently adding file:", file.Name())
// 		wg.Add(1)
// 		go getDataFromFile(path+"/"+file.Name(), db)
// 	}
// 	wg.Wait()
// }

// func decodeHex(s string) []byte {
// 	bytes, _ := hex.DecodeString(s)
// 	return bytes
// }

// func checkDB(b string) func(txn *badger.Txn) error {
// 	return func(txn *badger.Txn) error {
// 		item, err := txn.Get(hex2Binary(b))
// 		handle(err)

// 		val, err := item.ValueCopy(nil)
// 		handle(err)

// 		fmt.Printf("%s The value of %s, is %s\n", time.Now().Format("Stamp"), b, string(val))
// 		return nil
// 	}
// 	//	return func(txn *badger.Txn) error {
// 	//		defer wg.Done()
// 	//		item, err := txn.Get(decodeHex(b))
// 	//		if err != nil { return err }
// 	//
// 	//		val, err := item.ValueCopy(nil)
// 	//		if err != nil { return err }
// 	//
// 	//		fmt.Printf("%s The value of %s, is %s\n", time.Now(), b, string(val))
// 	//		return nil
// 	//	}
// }
type transaction struct {
	seed string
	hash string
}

var trxs []transaction

func hex2Binary(s string) []byte {
	bytes, _ := hex.DecodeString(s)
	return bytes
}

func doWritesToDB(db *badger.DB) {
	wb := db.NewWriteBatch()
	defer wb.Cancel()
	for _, v := range trxs {
		sstr, hstr := v.seed, v.hash
		fmt.Println(sstr, hstr)
		hashBytes, _ := hex.DecodeString(hstr)
		wb.Set(hashBytes, []byte(sstr))
	}
	wb.Flush()
	trxs = nil
}

func saveToDb(s, h string, db *badger.DB) {
	trxs = append(trxs, transaction{s, h})
	if len(trxs) >= 10 {
		fmt.Println("saving data to database.", len(trxs), trxs)
		doWritesToDB(db)
	}
}

func exists(b, original string) func(txn *badger.Txn) error {
	return func(txn *badger.Txn) error {
		defer wg.Done()
		item, err := txn.Get(hex2Binary(b))
		if err != nil {
			return err
		}
		val, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}
		if string(val) == original {
			return nil
		}
		fmt.Printf("%s The value of %s, is %s\n", timeStamp(), original, string(val))
		return nil
	}
}

func connectToDB() *badger.DB {
	opts := badger.DefaultOptions("/var/lib/badger")
	// opts.Logger = nil
	db, err := badger.Open(opts)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\t - \tDB opened\n", timeStamp())
	return db
}
