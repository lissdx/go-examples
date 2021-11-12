package main

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	apiv1 "github.com/lissdx/go-examples/key-value-store/api/v1"
	"io/ioutil"
	"log"
	"net/http"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func notAllowedHandler(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, "Not Allowed", http.StatusMethodNotAllowed)
}

func keyValuePutHandler(kvsStorekeeper apiv1.KeyValueStorekeeper, kvsLogger apiv1.TransactionLogger) func(w http.ResponseWriter, r *http.Request) {
	var storekeeper, transactionLogger = kvsStorekeeper, kvsLogger
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		key := vars["key"]

		value, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = storekeeper.Put(key, string(value))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		transactionLogger.WritePut(key, string(value))

		w.WriteHeader(http.StatusCreated)

		log.Printf("PUT key=%s value=%s\n", key, string(value))
	}

}

func keyValueGetHandler(kvsStorekeeper apiv1.KeyValueStorekeeper) func(w http.ResponseWriter, r *http.Request) {
	var storekeeper = kvsStorekeeper
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		key := vars["key"]

		value, err := storekeeper.Get(key)
		if errors.Is(err, apiv1.ErrorNoSuchKey) {
			http.Error(w, fmt.Errorf("%w, key=%s", err, key).Error(), http.StatusNoContent)
			log.Printf("%v,GET key=%s", err, key)
			return
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write([]byte(value))

		log.Printf("GET key=%s value=%s\n", key, value)
	}
}

func keyValueDeleteHandler(kvsStorekeeper apiv1.KeyValueStorekeeper, kvsLogger apiv1.TransactionLogger) func(w http.ResponseWriter, r *http.Request) {
	var storekeeper, transactionLogger = kvsStorekeeper, kvsLogger
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		key := vars["key"]

		err := storekeeper.Delete(key)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		transactionLogger.WriteDelete(key)

		log.Printf("DELETE key=%s\n", key)
	}

}

func initTransactionLogger(transactionLogger apiv1.TransactionLogger, kvsStore apiv1.KeyValueStorekeeper)  error {
	eventStream, errorStream := transactionLogger.ReadEvents()
	var err error
	count, ok, e := 0, true, apiv1.Event{}

	for ok && err == nil {
		select {
		case err, ok = <- errorStream:
		case e, ok = <- eventStream:
			switch e.EventType {
			case apiv1.EventDelete:
				_ = kvsStore.Delete(e.Key)
				log.Printf("key DELETED: %v", e.Key)
				count++
			case apiv1.EventPut:
				_ = kvsStore.Put(e.Key, e.Value)
				log.Printf("key PUTED: key:%v, value:%v", e.Key, e.Value)
				count++
			}
		}
	}

	transactionLogger.Run()
	return err
}

func main() { // Create a new mux router
	kvStorage := apiv1.NewKeyValueStore()
	kvsTransactionLogger, err := apiv1.NewSqliteTransactionLogger(apiv1.SqliteTransactionLoggerSettings{DataSourceName: "./put-delete.db", DriveName: "sqlite3"})

	if err != nil{
		log.Fatal(err)
	}

	if err := initTransactionLogger(kvsTransactionLogger, kvStorage); err != nil{
		log.Fatalf("can't run initTransactionLogger err: %v", err)
	}

	r := mux.NewRouter()

	r.Use(loggingMiddleware)

	r.HandleFunc("/v1/{key}", keyValueGetHandler(kvStorage)).Methods("GET")
	r.HandleFunc("/v1/{key}", keyValuePutHandler(kvStorage, kvsTransactionLogger)).Methods("PUT")
	r.HandleFunc("/v1/{key}", keyValueDeleteHandler(kvStorage, kvsTransactionLogger)).Methods("DELETE")

	r.HandleFunc("/v1", notAllowedHandler)
	r.HandleFunc("/v1/{key}", notAllowedHandler)

	log.Fatal(http.ListenAndServe(":8080", r))
}
