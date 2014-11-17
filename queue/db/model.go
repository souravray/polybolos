/*
* @Author: souravray
* @Date:   2014-11-08 00:57:48
* @Last Modified by:   souravray
* @Last Modified time: 2014-11-18 00:11:05
 */

package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"sync"
)

type Model struct {
	TableName string
	DB        *sql.DB
	Tx        *sql.Tx
	RWLock    sync.RWMutex
}

type QueueIteam struct {
	key  sql.NullString
	task sql.NullString
}

func NewModel(connectionString, name string) (model *Model, err error) {
	var db *sql.DB
	db, err = sql.Open("sqlite3", connectionString)
	if err != nil {
		return
	}
	model = &Model{TableName: name, DB: db}
	query := `
		PRAGMA automatic_index = OFF;
		PRAGMA cache_size = 32768;
		PRAGMA cache_spill = OFF;
		PRAGMA foreign_keys = OFF;
		PRAGMA journal_size_limit = 67110000;
		PRAGMA locking_mode = NORMAL;
		PRAGMA page_size = 4096;
		PRAGMA recursive_triggers = OFF;
		PRAGMA secure_delete = OFF;
		PRAGMA synchronous = FULL;
		PRAGMA temp_store = MEMORY;
		PRAGMA journal_mode = WAL;
		PRAGMA wal_autocheckpoint = 16384;

		CREATE TABLE IF NOT EXISTS queue (
        key text not null primary key,
        task blob
    );
    `
	_, err = model.DB.Exec(query)
	if err != nil {
		log.Fatal(err)
		return
	}
	return
}

func (m *Model) begin() {
	m.Tx, _ = m.DB.Begin()
}

func (m *Model) BatchTransaction() {
	m.RWLock.Lock()
	if m.Tx != nil {
		oldTx := m.Tx
		defer oldTx.Commit()
	}
	m.begin()
	m.RWLock.Unlock()
}

func (m *Model) TransactionEnd() {
	if m.Tx != nil {
		m.RWLock.Lock()
		m.Tx.Commit()
		m.RWLock.Unlock()
	}
}

func (m *Model) Add(key string, task []byte) (err error) {
	m.RWLock.RLock()
	stmt, err := m.Tx.Prepare("INSERT INTO queue (key, task) VALUES (?, ?)")
	defer stmt.Close()
	if err != nil {
		log.Fatal(err)
		return
	}
	_, err = stmt.Exec(key, task)
	m.RWLock.RUnlock()
	if err != nil {
		log.Fatal(err)
		return
	}
	return
}

func (m *Model) Update(key string, task []byte) (err error) {
	m.RWLock.RLock()
	stmt, err := m.Tx.Prepare("UPDATE queue SET task = ? WHERE key = ?")
	defer stmt.Close()
	if err != nil {
		log.Fatal(err)
		return
	}
	_, err = stmt.Exec(task, key)
	if err != nil {
		log.Fatal(err)
		return
	}
	m.RWLock.RUnlock()
	return
}

func (m *Model) Delete(key string) (err error) {
	m.RWLock.RLock()
	stmt, err := m.Tx.Prepare("DELETE FROM queue WHERE key = ?")
	defer stmt.Close()
	if err != nil {
		log.Fatal(err)
		return
	}
	_, err = stmt.Exec(key)
	if err != nil {
		log.Fatal(err)
		return
	}
	m.RWLock.RUnlock()
	return
}

// func (m *Model) Read() chan *bytes.Buffer {
// }
