package apiv1

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"sync"
)

//db, err := sql.Open("vvv", "./foo.db")

var OpenDbError = errors.New("can't open db")
var ConnectDbError = errors.New("can't connect db")

type SqliteTransactionLoggerSettings struct {
	DriveName      string
	DataSourceName string
}

type SqliteTransactionLogger struct {
	eventStream chan<- Event
	errorStream <-chan error
	db          *sql.DB
	wg          sync.WaitGroup
}

func (s *SqliteTransactionLogger) WriteDelete(key string) {
	log.Printf("SqliteTransactionLogger -> WriteDelete: key=%s", key)
	s.wg.Add(1)
	s.eventStream <- Event{EventType: EventDelete, Key: key}
}

func (s *SqliteTransactionLogger) WritePut(key, value string) {
	log.Printf("SqliteTransactionLogger -> WritePut: key=%s, value=%s", key, value)
	s.wg.Add(1)
	s.eventStream <- Event{EventType: EventPut, Key: key, Value: value}
}

func (s *SqliteTransactionLogger) Err() <-chan error {
	return s.errorStream
}

func (s *SqliteTransactionLogger) Close() {
	s.Close()
}

func (s *SqliteTransactionLogger) Wait() {
	s.wg.Wait()
}

func (s *SqliteTransactionLogger)Run()  {
	events := make(chan Event)
	errors := make(chan error)
	s.eventStream = events
	s.errorStream = errors

	go func() {
		query := `INSERT INTO transactions
			(event_type, key, value)
			VALUES ($1, $2, $3)`

		for e := range events {
			fmt.Printf("OK")
			if _, err := s.db.Exec(query, e.EventType, e.Key, e.Value); err != nil{
				errors <- err
			}
			s.wg.Done()
		}
	}()
}

func (s *SqliteTransactionLogger)createTable() error {
	query := "CREATE TABLE IF NOT EXISTS'transactions' ('id' INTEGER,'event_type' INTEGER,'key' TEXT,'value' TEXT, PRIMARY KEY('id' AUTOINCREMENT))"

	_, err := s.db.Exec(query)
	return err
}


func (s *SqliteTransactionLogger) ReadEvents() (<-chan Event, <-chan error) {
	outEventStream := make(chan Event)
	outErrorStream := make(chan error)

	fetchQuery := `SELECT id, event_type, "key", value from transactions t  ORDER by id`

	go func() {
		defer close(outEventStream)
		defer close(outErrorStream)

		rows, err := s.db.Query(fetchQuery)

		defer func(rows *sql.Rows) {
			err := rows.Close()
			if err != nil {
				outErrorStream <- err
			}
		}(rows)

		if err != nil{
			outErrorStream <- err
			return
		}

		var e Event
		for rows.Next() {
			err := rows.Scan(&e.Sequence, &e.EventType, &e.Key, &e.Value)
			if err != nil {
				outErrorStream <- err
				return
			}

			outEventStream <- e
		}

		if err := rows.Err(); err != nil{
			outErrorStream <- err
		}
	}()

	return outEventStream, outErrorStream
}


func NewSqliteTransactionLogger(settings SqliteTransactionLoggerSettings) (TransactionLogger, error) {

	db, err := sql.Open(settings.DriveName, settings.DataSourceName)
	if err != nil {
		return nil, OpenDbError
	}

	if db.Ping() != nil {
		return nil, ConnectDbError
	}

	sqliteLogger := &SqliteTransactionLogger{db: db}

	if err := sqliteLogger.createTable(); err != nil {
		return nil, err
	}

	//sqliteLogger.Run()

	return sqliteLogger, nil
}
