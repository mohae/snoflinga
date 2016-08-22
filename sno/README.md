# snö

snö: snowflake like IDs with millisecond time resolution, a 10 bit sequence, a 64 bit ID, and a 13 bit secondary ID.

The timestamp is from a custom epoch: 1/1/2016 00:00:00.

The sequence number starts at a random number and increases with each snowflake and rolls over.  It should not be relied on for anything.

A 13-bit secondary ID, SID, element is provided.

The ID is 8 bytes and no assumptions are made about its contents or layout.  A snowflake generator is used for a single ID and SID combination.

## 128 bit
128 bit snowflakes have the following layout:

```
bits    
 0-40     Timestamp, in milliseconds
41-50     Sequence number
51-63     SID: secondary ID
64-127    ID
```
