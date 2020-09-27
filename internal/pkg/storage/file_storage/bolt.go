package file_storage

import (
	"encoding/base64"
	"encoding/json"
	"github.com/0ndi/go-privnote/internal/pkg/note"
	"github.com/boltdb/bolt"
	"github.com/pkg/errors"
)

const (
	BucketName = "notes"
)

type Storage struct {
	db *bolt.DB
}

func GetDB(dbPath string) (*bolt.DB, error) {
	db, err := bolt.Open(dbPath, 0600, nil)
	if err != nil {
		err := errors.Wrap(err, "bolt.Open")
		return nil, err
	}

	if err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(BucketName))
		if err != nil {
			err := errors.Wrap(err, "CreateBucketIfNotExists")
			return err
		}
		return nil
	}); err != nil {
		err := errors.Wrap(err, "Update")
		return nil, err
	}
	return db, nil
}

func NewStorage(db *bolt.DB) (*Storage, error) {
	return &Storage{db: db}, nil
}

//Save return note key wich used in url
func (s Storage) Save(n *note.Note) (string, error) {
	hash := n.Hash()
	data, err := json.Marshal(n)
	if err != nil {
		err := errors.Wrap(err, "Marhasl")
		return "", err
	}

	err = s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BucketName))
		if v := b.Get(hash[:note.KeyLength]); v == nil {
			hash = hash[:note.KeyLength]
		}
		if err := b.Put(hash, data); err != nil {
			err := errors.Wrap(err, "db.Put")
			return err
		}
		return nil
	})

	url := base64.URLEncoding.EncodeToString(hash)
	return url, err
}

func (s *Storage) GetNote(key string) (*note.Note, error) {
	rawKey, err := base64.URLEncoding.DecodeString(key)
	if err != nil {
		err := errors.Wrap(err, "DecodeString")
		return nil, err
	}

	var note note.Note
	err = s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BucketName))
		data := b.Get(rawKey)
		if data == nil {
			return errors.New("not found")
		}

		if err := json.Unmarshal(data, &note); err != nil {
			err := errors.Wrap(err, "Unmarshal")
			return err
		}

		return b.Delete(rawKey)
	})
	return &note, err
}
