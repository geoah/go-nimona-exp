package journal

import (
	"encoding/json"
	"fmt"
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
	journal     *SerialJournal
}

func TestJournalTestSuite(t *testing.T) {
	suite.Run(t, new(JournalTestSuite))
}

func (s *JournalTestSuite) SetupTest() {
	s.userID = "user_id"
	s.persistence = store.NewInMemoryStore()
	s.journal = NewJournal(s.persistence)
}

func (s *JournalTestSuite) TestIndexConvertion_Valid_Succeeds() {
	si0 := SerialIndex(0)
	bi0 := []byte{48}

	assert.Equal(s.T(), si0, indexFromBytes(bi0))
	assert.Equal(s.T(), bi0, indexToBytes(si0))

	fmt.Println("si0", si0)
	fmt.Println("bi0", bi0)
}

func (s *JournalTestSuite) TestPersistedRestore_Valid_Succeeds() {
	entry1payload := &testPayload{String: "entry1"}
	entry1payloadJSON, _ := json.Marshal(entry1payload)
	entry1index := SerialIndex(1)
	entry1 := NewSerialEntry(entry1index, entry1payloadJSON)

	entry2payload := &testPayload{String: "entry2"}
	entry2payloadJSON, _ := json.Marshal(entry2payload)
	entry2index := SerialIndex(2)
	entry2 := NewSerialEntry(entry2index, entry2payloadJSON)

	entry3payload := &testPayload{String: "entry3"}
	entry3payloadJSON, _ := json.Marshal(entry3payload)
	entry3index := SerialIndex(3)
	entry3 := NewSerialEntry(entry3index, entry3payloadJSON)

	index1, errEntry1 := s.journal.Restore(entry1)
	assert.Equal(s.T(), Index(indexToBytes(entry1index)), index1)
	assert.Nil(s.T(), errEntry1)

	index2, errEntry2 := s.journal.Restore(entry2)
	assert.Equal(s.T(), Index(indexToBytes(entry2index)), index2)
	assert.Nil(s.T(), errEntry2)

	index3, errEntry3 := s.journal.Restore(entry3)
	assert.Equal(s.T(), Index(indexToBytes(entry3index)), index3)
	assert.Nil(s.T(), errEntry3)
}

func (s *JournalTestSuite) TestPersistedAppend_Valid_Succeeds() {
	entry1payload := &testPayload{String: "entry1"}
	entry1payloadJSON, _ := json.Marshal(entry1payload)
	entry1index := Index(indexToBytes(SerialIndex(1)))

	entry2payload := &testPayload{String: "entry2"}
	entry2payloadJSON, _ := json.Marshal(entry2payload)
	entry2index := Index(indexToBytes(SerialIndex(2)))

	entry3payload := &testPayload{String: "entry3"}
	entry3payloadJSON, _ := json.Marshal(entry3payload)
	entry3index := Index(indexToBytes(SerialIndex(3)))

	index1, errEntry1 := s.journal.Append(entry1payloadJSON)
	assert.Equal(s.T(), entry1index, index1)
	assert.Nil(s.T(), errEntry1)

	index2, errEntry2 := s.journal.Append(entry2payloadJSON)
	assert.Equal(s.T(), entry2index, index2)
	assert.Nil(s.T(), errEntry2)

	index3, errEntry3 := s.journal.Append(entry3payloadJSON)
	assert.Equal(s.T(), entry3index, index3)
	assert.Nil(s.T(), errEntry3)
}

func (s *JournalTestSuite) TestAppend_InvalidParent_Failes() {
	entry2payload := &testPayload{String: "entry2"}
	entry2payloadJSON, _ := json.Marshal(entry2payload)
	entry2 := NewSerialEntry(2, entry2payloadJSON)

	lastIndex, errEntry2 := s.journal.Restore(entry2)
	assert.Equal(s.T(), Index(indexToBytes(SerialIndex(0))), lastIndex)
	assert.Equal(s.T(), ErrMissingParentIndex, errEntry2)
}
