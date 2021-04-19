package session

import (
	"encoding/json"
	"errors"
	"github.com/adigunhammedolalekan/cashtroops/types"
	"github.com/dgraph-io/badger/v2"
	"time"
)

const (
	sevenDays = time.Hour * 24 * 7
)

var (
	ErrTokenExpired    = errors.New("token has expired. please re-authenticate")
	ErrUnAuthenticated = errors.New("unauthenticated user. token not found")
)

type Store interface {
	Create(token *types.Token) error
	Get(key string) (*types.User, error)
}

type sessionStore struct {
	db *badger.DB
}

func New(path string) (Store, error) {
	db, err := badger.Open(badger.DefaultOptions(path))
	if err != nil {
		return nil, err
	}
	return &sessionStore{db: db}, nil
}

func (store *sessionStore) Create(token *types.Token) error {
	data, err := json.Marshal(token)
	if err != nil {
		return err
	}
	return store.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(token.Key), data)
	})
}

func (store *sessionStore) Get(key string) (*types.User, error) {
	token := &types.Token{}
	err := store.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		var value []byte
		err = item.Value(func(val []byte) error {
			value = append(value, val...)
			return nil
		})
		if err != nil {
			return err
		}
		if err := json.Unmarshal(value, token); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, ErrUnAuthenticated
	}
	diff := time.Now().Sub(token.Created)
	if diff > sevenDays {
		return nil, ErrTokenExpired
	}
	return token.User, nil
}
