package note

import (
	"crypto/sha256"
	"fmt"
	"time"
)

const (
	KeyLength = 6
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

func (n *Note) Hash() []byte {
	hash := sha256.New()
	hash.Write([]byte(fmt.Sprintf("%s%d", n.Data, n.ExpiredAt.UnixNano())))
	return hash.Sum(nil)
}
