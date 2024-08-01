package persistence

import (
	"encoding/json"
	"log/slog"
	"reflect"
	"strings"

	bolt "go.etcd.io/bbolt"
)

type DAO[T DataObject] struct {
	db     *bolt.DB
	bucket []byte
}

func NewDAO[T DataObject](bucket string) *DAO[T] {
	slog.Debug("New DAO with", "bucket", bucket, "type", reflect.TypeOf(new(T)))
	return &DAO[T]{database(), []byte(bucket)}
}

func (dao *DAO[T]) Find(query string) []T {
	tx, err := dao.db.Begin(false)
	if err != nil {
		slog.Error(err.Error())
		return nil
	}
	defer tx.Rollback()

	result := []T{}
	err = tx.Bucket(dao.bucket).ForEach((func(k, v []byte) error {
		if query != "" {
			if !strings.Contains(strings.ToLower(string(k)), strings.ToLower(query)) {
				return nil
			}
		}

		var item T
		err := json.Unmarshal(v, &item)
		if err == nil {
			result = append(result, item)
		}
		return err
	}))
	if err != nil {
		slog.Error(err.Error())
	}
	return result
}

func (dao *DAO[T]) Save(item T) error {
	tx, err := dao.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	b, err := tx.CreateBucketIfNotExists(dao.bucket)
	if err != nil {
		return err
	}

	item.UpdateTimestamp()
	data, err := json.Marshal(item)
	if err != nil {
		return err
	}

	if err = b.Put(item.Identifier(), data); err == nil {
		err = tx.Commit()
	}
	return err
}

type DataObject interface {
	Identifier() []byte
	UpdateTimestamp()
}
