# sne

sne: snowflake like IDs with millisecond time resolution, 11 bit sequence, 56 bit key, and 20 bit secondary key.

The timestamp is from a custom epoch: 1/1/2016 00:00:00.

The sequence number is 11 bits and starts at a random number and increases with each snowflake and rolls over.  It should not be relied on for anything.

A 20-bit secondary key, SID, element is provided.

The key, or ID, is 7 bytes and no assumptions are made about its contents or layout.  A snowflake generator is used for a single ID and SID combination.

## 128 bit
128 bit snowflakes have the following layout:

```
bits    
 0-40     Timestamp, in milliseconds
41-51     Sequence number
52-71     SID: secondary ID
72-127    ID
```
