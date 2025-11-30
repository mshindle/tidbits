# Cache

When you need a simple cache with Set, Get and TTL, but you do  not need the
overhead you get from Redis.

## Simple

Here's a straightforward approach that will work if key size is small.
 * Store key-value pairs in memory
 * Expire keys after a certain time
 * Handle concurrent reads & writes
 * Small key size (< 10,000)

## LRU
