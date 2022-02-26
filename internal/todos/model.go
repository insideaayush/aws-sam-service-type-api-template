package todos

import (
	"bytes"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
)

type Todo struct {
	PartitionKey string    `json:"pk" dynamodbav:"pk"`
	SortKey      string    `json:"sk" dynamodbav:"sk"`
	ID           uuid.UUID `json:"id" dynamodbav:"id"`
	CreatedAt    string    `json:"created_at" dynamodbav:"created_at"`
	UpdatedAt    string    `json:"updated_at" dynamodbav:"updated_at"`

	Title  string `json:"title" dynamodbav:"title"`
	Body   string `json:"body" dynamodbav:"body"`
	Status Status `json:"status" dynamodbav:"status"`
}

func (td *Todo) IsValid() bool {
	return !(td.Title == "")
}

type Status int

const (
	CREATED Status = iota + 1
	INPROGRESS
	BLOCKED
	ONHOLD
	COMPLETED
	BACKLOG
)

func (s Status) String() string {
	m := map[Status]string{
		CREATED:    "CREATED",
		INPROGRESS: "INPROGRESS",
		BLOCKED:    "BLOCKED",
		ONHOLD:     "ONHOLD",
		COMPLETED:  "COMPLETED",
		BACKLOG:    "BACKLOG",
	}
	return m[s]
}

func ToStatus(s string) (Status, error) {
	m := map[string]Status{
		"CREATED":    CREATED,
		"INPROGRESS": INPROGRESS,
		"BLOCKED":    BLOCKED,
		"ONHOLD":     ONHOLD,
		"COMPLETED":  COMPLETED,
		"BACKLOG":    BACKLOG,
	}
	if _, ok := m[s]; !ok {
		return -1, errors.New("invalid status")
	}
	return m[s], nil
}

func (s Status) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(s.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (s *Status) UnmarshalJSON(b []byte) error {
	var j string
	if err := json.Unmarshal(b, &j); err != nil {
		return err
	}
	*s, _ = ToStatus(j)
	return nil
}
