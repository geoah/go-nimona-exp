# Journal

Nimona's Journal represents the simplest form of an transaction log.  
Its only purpose is to store binary messages that can later be retrieved
in the same order they were written in.

Other services can listen for new Entries by adding themselves as Notifiees,
using `Notify()`. When a new Entry has been added to the journal (after it 
has been persisted) the Notifiee's `ProcessJournalEntry()` method will be
called with the new entry.

## Replication

For most parts the Journal doesn't really care how services handle replication,
synchronization, etc. Nevertheless, the `Restore()` method can be used to
restore entries (index & payload) on a replica journal; if the previous index
is not present, it will return an error as well as the last available index.
This can be used to figure out which entries are missing on the replica.

## Terminology

#### Index

The Index is an incremental number (uint64) that holds the entry's position
in the log. The first Entry in the log has an Index of 0 (zero).

#### Payload

The payload of the Entry. This is simply a byte array. The Journal doesn't
care about what it contains, its encoding, or anything else. It will be stored
and retrieved as-is.

#### Entry

An Entry is just a pair of an Index and a Payload and represents a single
entity in our Journal.

#### Journal

A Journal is a collection of Entries.

* There cannot be more than one Entry with the same Index.
* When appending a new Entry in the Journal, all Indexes up to the one we are
  trying to append must exist. If the latest present Index is 10, the next
  Entry to be added MUST have an Index of 11.

## Literature

* [The Log: What every software engineer should know about real-time data's unifying abstraction](https://engineering.linkedin.com/distributed-systems/log-what-every-software-engineer-should-know-about-real-time-datas-unifying) by Jay Kreps
* [Using logs to build a solid data infrastructure](http://www.confluent.io/blog/using-logs-to-build-a-solid-data-infrastructure-or-why-dual-writes-are-a-bad-idea/) by Martin Kleppmann
* [Raft: In Search of an Understandable Consensus Algorithm](https://raft.github.io/raft.pdf) by Diego Ongaro and John Ousterhout