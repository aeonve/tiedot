package main

import (
	"encoding/json"
	"fmt"
	"github.com/HouzuoGuo/tiedot/db"
	"os"
)

// You are encouraged to use (nearly) all tiedot public functions concurrently.
// There are few exceptions - see individual package/functions for details.

func embeddedExample() {
	dir := "/tmp/MyDatabase"
	os.RemoveAll(dir)
	defer os.RemoveAll(dir)

	// Open database
	myDB, err := db.OpenDB(dir)
	if err != nil {
		panic(err)
	}

	// Create collection
	if err := myDB.Create("A"); err != nil {
		panic(err)
	}
	if err := myDB.Create("B"); err != nil {
		panic(err)
	}

	// Rename collection
	if err := myDB.Rename("B", "C"); err != nil {
		panic(err)
	}

	// Which collections do I have?
	for name := range myDB.StrCol {
		fmt.Printf("I have a collection called %s\n", name)
	}

	// Drop collection
	if err := myDB.Drop("C"); err != nil {
		panic(err)
	}

	// Start using collection
	A := myDB.Use("A")

	// Collection insert/update/delete operations require the document to be a map[string]interface{}, otherwise index may not work

	// Insert a docuemnt
	docID, err := A.Insert(map[string]interface{}{"Url": "http://google.com", "Owner": "Google Inc."})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Inserted document at %d (document ID)\n", docID)

	// Update document
	var doc map[string]interface{}
	json.Unmarshal([]byte(`{"Url": "http://www.google.com.au", "Owner": "Google Inc."}`), &doc)
	newID, err := A.Update(docID, doc)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Updated document %d to %v, new ID is %d\n", docID, doc, newID)
	// Updated document may or may not retain the original ID.
	// Optionally, each document may have a unique/persistent/never changing ID assigned to it, this is called "UID"
	// Check out "InsertWithUID", "ReadByUID", "UpdateByUID", "DeleteByUID"

	// Read document
	var readback map[string]interface{}
	if err := A.Read(newID, &readback); err != nil {
		panic(err)
	}
	fmt.Printf("Read document ID %d: %v\n", newID, readback)

	// Delete document
	A.Delete(123) // An ID which does not exist does no harm

	// Create index
	if err := A.Index([]string{"a", "b", "c"}); err != nil {
		panic(err)
	}

	// Which indexes do I have on collection A?
	for path := range A.StrHT {
		fmt.Printf("I have an index on path %s\n", path)
	}

	// Remove index
	if err := A.Unindex([]string{"a", "b", "c"}); err != nil {
		panic(err)
	}

	// Execute query
	result := make(map[uint64]struct{})
	var query interface{}
	json.Unmarshal([]byte(`"all"`), &query)
	if err := db.EvalQueryV2(query, A, &result); err != nil {
		panic(err)
	}
	for id := range result {
		// Map keys are query results - result document IDs
		fmt.Printf("Query returned document ID %d\n", id)
		// Optionally, use A.read(ID) to get actual document content from query result
	}

	// Gracefully close database
	myDB.Close()
}
