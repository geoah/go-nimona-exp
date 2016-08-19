package journal

import (
	"encoding/json"
	"os"
	"testing"

	mj "github.com/jbenet/go-multicodec/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type testPayload struct {
	String string `json:"str"`
}

type JournalTestSuite struct {
	suite.Suite
	userID   string
	journal  *SequentialJournal
	file     *os.File
	filePath string
}

func TestJournalTestSuite(t *testing.T) {
	suite.Run(t, new(JournalTestSuite))
}

func (s *JournalTestSuite) SetupTest() {
	s.userID = "user_id"

	s.filePath = "/tmp/nimona-journal-test.mjson"
	f, err := os.OpenFile(s.filePath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0777)
	if err != nil {
		panic(err)
	}

	s.file = f
	mc := mj.Codec(true)
	s.journal = NewJournal(mc, s.file, s.file)
}

func (s *JournalTestSuite) TeardownTest() {
	s.file.Close()
	os.Remove(s.filePath)
}

func (s *JournalTestSuite) TestFilePersistence_Valid_Succeeds() {
	entry1payload := &testPayload{String: "entry1"}
	entry1payloadJSON, _ := json.Marshal(entry1payload)
	entry1index := Index(1)

	entry2payload := &testPayload{String: "entry2"}
	entry2payloadJSON, _ := json.Marshal(entry2payload)
	entry2index := Index(2)

	entry3payload := &testPayload{String: "entry3"}
	entry3payloadJSON, _ := json.Marshal(entry3payload)
	entry3index := Index(3)

	index1, errEntry1 := s.journal.Append(entry1payloadJSON)
	assert.Equal(s.T(), entry1index, index1)
	assert.Nil(s.T(), errEntry1)

	index2, errEntry2 := s.journal.Append(entry2payloadJSON)
	assert.Equal(s.T(), entry2index, index2)
	assert.Nil(s.T(), errEntry2)

	index3, errEntry3 := s.journal.Append(entry3payloadJSON)
	assert.Equal(s.T(), entry3index, index3)
	assert.Nil(s.T(), errEntry3)

	s.file.Close()

	f, err := os.OpenFile(s.filePath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0777)
	if err != nil {
		panic(err)
	}

	s.file = f
	mc := mj.Codec(true)
	journal := NewJournal(mc, s.file, s.file)

	nm := &MockNotifiee{}
	journal.Notify(nm)

	nm.On("ProcessJournalEntry", NewSequentialEntry(entry1index, entry1payloadJSON)).Return()
	nm.On("ProcessJournalEntry", NewSequentialEntry(entry2index, entry2payloadJSON)).Return()
	nm.On("ProcessJournalEntry", NewSequentialEntry(entry3index, entry3payloadJSON)).Return()

	err = journal.Replay()
	assert.Nil(s.T(), err)
}

func (s *JournalTestSuite) TestPersistedRestore_Valid_Succeeds() {
	entry1payload := &testPayload{String: "entry1"}
	entry1payloadJSON, _ := json.Marshal(entry1payload)
	entry1index := Index(1)
	entry1 := NewSequentialEntry(entry1index, entry1payloadJSON)

	entry2payload := &testPayload{String: "entry2"}
	entry2payloadJSON, _ := json.Marshal(entry2payload)
	entry2index := Index(2)
	entry2 := NewSequentialEntry(entry2index, entry2payloadJSON)

	entry3payload := &testPayload{String: "entry3"}
	entry3payloadJSON, _ := json.Marshal(entry3payload)
	entry3index := Index(3)
	entry3 := NewSequentialEntry(entry3index, entry3payloadJSON)

	index1, errEntry1 := s.journal.Restore(entry1)
	assert.Equal(s.T(), entry1index, index1)
	assert.Nil(s.T(), errEntry1)

	index2, errEntry2 := s.journal.Restore(entry2)
	assert.Equal(s.T(), entry2index, index2)
	assert.Nil(s.T(), errEntry2)

	index3, errEntry3 := s.journal.Restore(entry3)
	assert.Equal(s.T(), entry3index, index3)
	assert.Nil(s.T(), errEntry3)
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
	entry2 := NewSequentialEntry(2, entry2payloadJSON)

	lastIndex, errEntry2 := s.journal.Restore(entry2)
	assert.Equal(s.T(), Index(0), lastIndex)
	assert.Equal(s.T(), ErrMissingParentIndex, errEntry2)
}
