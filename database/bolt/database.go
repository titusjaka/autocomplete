package bolt

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"time"

	"github.com/boltdb/bolt"
	"github.com/pkg/errors"

	"github.com/titusjaka/autocomplete"
	"github.com/titusjaka/autocomplete/database"
)

const defaultBucketName = "autocomplete"

// DB is a bolt implementation of Database interface
// It uses boltDB and saves data on disk
type DB struct {
	conn   *bolt.DB
	logger *log.Logger
}

// New returns new bolt DB structure
func New(path string) (*DB, error) {
	db, err := bolt.Open(path, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, errors.Wrapf(err, "couldn't open database '%s'", path)
	}
	if err = db.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucket([]byte(defaultBucketName))
		if err != nil && err != bolt.ErrBucketExists {
			return err
		}
		return nil
	}); err != nil {
		return nil, errors.Wrap(err, "failed to create bucket")
	}

	return &DB{
		conn:   db,
		logger: log.New(ioutil.Discard, "", 0),
	}, nil
}

// SetLogger sets logger for DB structure
func (d *DB) SetLogger(logger *log.Logger) error {
	if logger == nil {
		return errors.New("logger is nil")
	}
	d.logger = logger
	d.logger.Printf("[DEBUG] logger set successfully")
	return nil
}

// Get fetches data from bolt DB. If key is no presented in DB, database.ErrNotFound returned
func (d *DB) Get(request string) ([]*autocomplete.WidgetData, error) {
	d.logger.Printf("[DEBUG] Getting data for request '%s' from bolt DB", request)
	var rawData []byte
	_ = d.conn.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(defaultBucketName))
		rawData = b.Get([]byte(request))
		return nil
	})
	if rawData == nil {
		d.logger.Printf("[DEBUG] No data found for request '%s'", request)
		return nil, database.ErrNotFound
	}
	var wd []*autocomplete.WidgetData
	if err := json.Unmarshal(rawData, &wd); err != nil {
		d.logger.Printf("[ERROR] Failed to unmarshal data fetched from bolt DB: %v", err)
		return nil, errors.Wrap(err, "failed to unmarshal widget data")
	}
	return wd, nil
}

// Write writes data to bolt DB
func (d *DB) Write(request string, data []*autocomplete.WidgetData) error {
	d.logger.Printf("[DEBUG] Writing data for request '%s' to bolt DB", request)
	rawData, err := json.Marshal(data)
	if err != nil {
		d.logger.Printf("[ERROR] Failed to marshal data to JSON: %v", err)
		return errors.Wrap(err, "failed to marshal widget data")
	}
	return d.conn.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(defaultBucketName))
		err = b.Put([]byte(request), rawData)
		if err != nil {
			d.logger.Printf("[ERROR] Failed to write data to bolt DB, bucket '%s': %v", defaultBucketName, err)
		}
		return err
	})
}

// Close closes connection to bolt DB
func (d *DB) Close() error {
	return d.conn.Close()
}
