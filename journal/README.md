# Journal

Nimona's Journal represents the simplest form of an event log.  
Its only purpose is to store binary encoded events that can later be retrieved
in the same order they were written in.

## Replication

For most parts the Journal doesn't really care how services handle replication,
synchronization, etc. Nevertheless, the `Restore()` method can be used to
restore entries (index & payload) on a replica journal; if the previous index
is not present, it will return an error as well as the last available index.
This can be used to figure out which entries are missing on the replica.

## Literature

* [The Log: What every software engineer should know about real-time data's unifying abstraction](https://engineering.linkedin.com/distributed-systems/log-what-every-software-engineer-should-know-about-real-time-datas-unifying) by Jay Kreps
* [Using logs to build a solid data infrastructure](http://www.confluent.io/blog/using-logs-to-build-a-solid-data-infrastructure-or-why-dual-writes-are-a-bad-idea/) by Martin Kleppmann
* [Raft: In Search of an Understandable Consensus Algorithm](https://raft.github.io/raft.pdf) by Diego Ongaro and John Ousterhout