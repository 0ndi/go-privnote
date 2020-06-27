package note

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"time"
)

const (
	BucketName = "notes"
)

type Note struct {
	Data      string
	ExpiredAt time.Time
}

func NewNote(data string) *Note {
	return &Note{Data: data}
}

func (n *Note) hash() []byte {
	hash := sha256.New()
	hash.Write([]byte(fmt.Sprintf("%s%d", n.Data, n.ExpiredAt.Unix())))
	return hash.Sum(nil)
}

//Save return note key witch used in url
func (n *Note) Save(db *bolt.DB) (string, error) {
	hash := n.hash()

	tx, err := db.Begin(true)
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	data, err := json.Marshal(n)
	if err != nil {
		return "", err
	}

	b := tx.Bucket([]byte(BucketName))
	if err := b.Put(hash, data); err != nil {
		return "", err
	}

	if err := tx.Commit(); err != nil {
		return "", err
	}

	url := base64.URLEncoding.EncodeToString(hash)
	return url, nil
}

func GetNote(db *bolt.DB, key string) (*Note, error) {
	rawKey, err := base64.URLEncoding.DecodeString(key)
	if err != nil {
		return nil, err
	}

	tx, err := db.Begin(false)
	if err != nil {
		return nil, err
	}

	b := tx.Bucket([]byte(BucketName))
	data := b.Get(rawKey)

	var note Note
	if err := json.Unmarshal(data, &note); err != nil {
		return nil, err
	}
	return &note, nil
}
