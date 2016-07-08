package journal

import (
	"encoding/json"
	"testing"

	"github.com/nimona/go-nimona/store"
	"github.com/nimona/go-nimona/stream"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type testPayload struct {
	String string `json:"str"`
}

type JournalTestSuite struct {
	suite.Suite
	userID      string
	persistence store.Store
	journal     *Journal
}

func TestJournalTestSuite(t *testing.T) {
	suite.Run(t, new(JournalTestSuite))
}

func (s *JournalTestSuite) SetupTest() {
	s.userID = "user_id"
	s.persistence = store.NewInMemoryStore()
	s.journal = NewJournal(s.userID, s.persistence)
}

func (s *JournalTestSuite) TestPersistedAppend_Valid_Succeeds() {
	entry1payload := &testPayload{String: "entry1"}
	entry1payloadJSON, _ := json.Marshal(entry1payload)
	entry1 := NewEntry(1, entry1payloadJSON)

	entry2payload := &testPayload{String: "entry2"}
	entry2payloadJSON, _ := json.Marshal(entry2payload)
	entry2 := NewEntry(2, entry2payloadJSON)

	entry3payload := &testPayload{String: "entry3"}
	entry3payloadJSON, _ := json.Marshal(entry3payload)
	entry3 := NewEntry(3, entry3payloadJSON)

	errEntry1 := s.journal.AppendEntry(entry1)
	assert.Nil(s.T(), errEntry1)

	errEntry2 := s.journal.AppendEntry(entry2)
	assert.Nil(s.T(), errEntry2)

	errEntry3 := s.journal.AppendEntry(entry3)
	assert.Nil(s.T(), errEntry3)

	getEntry2, err := s.journal.GetEntry(2)
	assert.Equal(s.T(), entry2, getEntry2)
	assert.Nil(s.T(), err)
}

func (s *JournalTestSuite) TestAppend_InvalidParent_Failes() {
	entry2payload := &testPayload{String: "entry2"}
	entry2payloadJSON, _ := json.Marshal(entry2payload)
	entry2 := NewEntry(2, entry2payloadJSON)

	errEntry2 := s.journal.AppendEntry(entry2)
	assert.Equal(s.T(), errEntry2, stream.ErrMissingParentIndex)
}
