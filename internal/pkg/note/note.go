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
	keyLength  = 6
)

var (
	defaultExpiration = time.Hour
)

type Note struct {
	Data      string
	ExpiredAt time.Time
}

func NewNote(data string) *Note {
	return &Note{Data: data, ExpiredAt: time.Now().Add(defaultExpiration)}
}

func (n *Note) hash() []byte {
	hash := sha256.New()
	hash.Write([]byte(fmt.Sprintf("%s%d", n.Data, n.ExpiredAt.Unix())))
	return hash.Sum(nil)
}

//Save return note key witch used in url
func (n *Note) Save(db *bolt.DB) (string, error) {
	hash := n.hash()
	data, err := json.Marshal(n)
	if err != nil {
		return "", err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BucketName))
		if v := b.Get(hash[:keyLength]); v == nil {
			hash = hash[:keyLength]
		}
		if err := b.Put(hash, data); err != nil {
			return err
		}
		return nil
	})

	url := base64.URLEncoding.EncodeToString(hash)
	return url, err
}

func GetNote(db *bolt.DB, key string) (*Note, error) {
	rawKey, err := base64.URLEncoding.DecodeString(key)
	if err != nil {
		return nil, err
	}

	var note Note
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BucketName))
		data := b.Get(rawKey)

		if err := json.Unmarshal(data, &note); err != nil {
			return err
		}

		return b.Delete(rawKey)
	})
	return &note, nil
}
