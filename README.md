# snöflinga
[![GoDoc](https://godoc.org/github.com/mohae/snoflinga?status.svg)](https://godoc.org/github.com/mohae/snoflinga)[![Build Status](https://travis-ci.org/mohae/snoflinga.png)](https://travis-ci.org/mohae/snoflinga)

snöflinga: snowflake like IDs

The timestamp is from Unix Epoch.

The sequence number starts at a random number and increases with each snowflake and rolls over.  It should not be relied on for anything.

The ID is 8 bytes and no assumptions are made about its contents or layout.  The snowflake generator is used for a single ID.

## 128 bit
128 bit snowflakes have the following layout:

```
bits    
0-51     Timestamp, in microseconds.
52-63    Sequence number.
64-127   ID
```
