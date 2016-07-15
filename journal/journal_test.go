package journal

import (
	"encoding/json"
	"testing"

	"github.com/nimona/go-nimona/store"
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
	s.journal = NewJournal(s.persistence)
}

func (s *JournalTestSuite) TestPersistedRestore_Valid_Succeeds() {
	entry1payload := &testPayload{String: "entry1"}
	entry1payloadJSON, _ := json.Marshal(entry1payload)
	entry1index := Index(1)
	entry1 := NewEntry(entry1index, entry1payloadJSON)

	entry2payload := &testPayload{String: "entry2"}
	entry2payloadJSON, _ := json.Marshal(entry2payload)
	entry2index := Index(2)
	entry2 := NewEntry(entry2index, entry2payloadJSON)

	entry3payload := &testPayload{String: "entry3"}
	entry3payloadJSON, _ := json.Marshal(entry3payload)
	entry3index := Index(3)
	entry3 := NewEntry(entry3index, entry3payloadJSON)

	index1, errEntry1 := s.journal.RestoreEntry(entry1)
	assert.Equal(s.T(), entry1index, index1)
	assert.Nil(s.T(), errEntry1)

	index2, errEntry2 := s.journal.RestoreEntry(entry2)
	assert.Equal(s.T(), entry2index, index2)
	assert.Nil(s.T(), errEntry2)

	index3, errEntry3 := s.journal.RestoreEntry(entry3)
	assert.Equal(s.T(), entry3index, index3)
	assert.Nil(s.T(), errEntry3)
}

func (s *JournalTestSuite) TestPersistedGet_Valid_Succeeds() {
	entry1payload := &testPayload{String: "entry1"}
	entry1payloadJSON, _ := json.Marshal(entry1payload)
	entry1index := Index(1)
	entry1 := NewEntry(entry1index, entry1payloadJSON)

	index1, errEntry1 := s.journal.RestoreEntry(entry1)
	assert.Equal(s.T(), entry1index, index1)
	assert.Nil(s.T(), errEntry1)

	getEntry1, err := s.journal.GetEntry(entry1index)
	assert.Equal(s.T(), entry1, getEntry1)
	assert.Nil(s.T(), err)
}

func (s *JournalTestSuite) TestPersistedAppend_Valid_Succeeds() {
	entry1payload := &testPayload{String: "entry1"}
	entry1payloadJSON, _ := json.Marshal(entry1payload)
	entry1index := Index(1)

	entry2payload := &testPayload{String: "entry2"}
	entry2payloadJSON, _ := json.Marshal(entry2payload)
	entry2index := Index(2)

	entry3payload := &testPayload{String: "entry3"}
	entry3payloadJSON, _ := json.Marshal(entry3payload)
	entry3index := Index(3)

	index1, errEntry1 := s.journal.AppendEntry(entry1payloadJSON)
	assert.Equal(s.T(), entry1index, index1)
	assert.Nil(s.T(), errEntry1)

	index2, errEntry2 := s.journal.AppendEntry(entry2payloadJSON)
	assert.Equal(s.T(), entry2index, index2)
	assert.Nil(s.T(), errEntry2)

	index3, errEntry3 := s.journal.AppendEntry(entry3payloadJSON)
	assert.Equal(s.T(), entry3index, index3)
	assert.Nil(s.T(), errEntry3)
}

func (s *JournalTestSuite) TestAppend_InvalidParent_Failes() {
	entry2payload := &testPayload{String: "entry2"}
	entry2payloadJSON, _ := json.Marshal(entry2payload)
	entry2 := NewEntry(2, entry2payloadJSON)

	lastIndex, errEntry2 := s.journal.RestoreEntry(entry2)
	assert.Equal(s.T(), Index(0), lastIndex)
	assert.Equal(s.T(), errEntry2, ErrMissingParentIndex)
}
