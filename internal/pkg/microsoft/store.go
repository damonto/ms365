package microsoft

import (
	"encoding/json"
	"os"
	"time"

	"github.com/damonto/msonline/internal/pkg/logger"
	"github.com/syndtr/goleveldb/leveldb"
)

// AccessToken is the microsoft graph api access token sturct
type AccessToken struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpireDate   time.Time `json:"expire_date"`
}

// Store is the microsoft graph api access token store.
type Store struct {
	db *leveldb.DB
}

// NewStore create access token store instance
func NewStore() *Store {
	dir, err := os.Getwd()
	if err != nil {
		logger.Sugar.Fatalf("Can get workdir %v", err)
		os.Exit(1)
	}

	db, err := leveldb.OpenFile(dir+"/database", nil)
	if err != nil {
		logger.Sugar.Fatalf("Unable to create database directory %v", err)
		os.Exit(1)
	}

	return &Store{
		db: db,
	}
}

// Put a new access token into level db
func (s *Store) Put(key string, accessToken AccessToken) error {
	defer s.db.Close()
	b, err := json.Marshal(accessToken)
	if err != nil {
		return err
	}

	err = s.db.Put([]byte(key), b, nil)
	if err != nil {
		return err
	}

	return nil
}

// Get retrieve a access token
func (s *Store) Get(key string) (accessToken AccessToken, err error) {
	defer s.db.Close()
	data, err := s.db.Get([]byte(key), nil)
	if err != nil {
		return accessToken, nil
	}

	err = json.Unmarshal(data, &accessToken)
	return
}

// Delete an item from level db
func (s *Store) Delete(key string) error {
	defer s.db.Close()
	return s.db.Delete([]byte(key), nil)
}

// All retrieve all access token from level db
func (s *Store) All() (accessToken []AccessToken, err error) {
	defer s.db.Close()
	iter := s.db.NewIterator(nil, nil)
Loop:
	for iter.Next() {
		var token AccessToken
		err := json.Unmarshal(iter.Value(), &token)
		if err != nil {
			break Loop
		}

		accessToken = append(accessToken, token)
	}
	iter.Release()
	if err = iter.Error(); err != nil {
		return accessToken, err
	}

	return accessToken, nil
}
