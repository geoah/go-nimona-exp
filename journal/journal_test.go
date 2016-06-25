package journal

import (
	"testing"

	"github.com/nimona/go-nimona/store"
	"github.com/nimona/go-nimona/stream"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

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
	entry1 := NewEntry(1, &Payload{String: "entry1"})
	entry2 := NewEntry(2, &Payload{String: "entry2"})
	entry3 := NewEntry(3, &Payload{String: "entry3"})

	errEntry1 := s.journal.AppendEntry(entry1)
	assert.Nil(s.T(), errEntry1)

	errEntry2 := s.journal.AppendEntry(entry2)
	assert.Nil(s.T(), errEntry2)

	errEntry3 := s.journal.AppendEntry(entry3)
	assert.Nil(s.T(), errEntry3)

	getEntry2, err := s.journal.GetEntry(entry2.GetIndex())
	assert.Equal(s.T(), entry2, getEntry2)
	assert.Nil(s.T(), err)
}

func (s *JournalTestSuite) TestAppend_InvalidParent_Failes() {
	entry2 := NewEntry(2, &Payload{String: "entry2"})

	errEntry2 := s.journal.AppendEntry(entry2)
	assert.Equal(s.T(), errEntry2, stream.ErrMissingParentIndex)
}
