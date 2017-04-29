

#### redishacker

knows nothing about elastic, moves data from hackernews into redis

#### elasticredis

knows nothing about hackernews, moves data from from redis into elastic

#### elastichacker

combo of **both** above

moves data from hackernews into elastic via redis

**cmd.bulkredis** talks to the interface **redisc.Hscan** through 2 parameters:

1. a [Redis Hashmap](https://redis.io/commands#hash) key where the field is the elastic document ID and the value is the JSON document.
2. a channel which contains the struct redisc.Datum with the elastic document ID and JSON document.



## Fetch

**Fetch** pulls the data out of Hacker News and writes it to Redis.

The data is categorized in Redis two different ways.

Both structures are hashes

##### By Index Name

key = indexname in this case hackernews  
field = hackernews ID   
value = GOB data structure {}   

```
type P struct {
	Itype string
	Id    int
	Json  []byte
}
```

##### By Hackernews Type {story, comment}

key = hackernews type  
field = hackernews ID  
value = JSON byte array  

## Elastic

Read data out of Redis and put it on a channel   
Each Hackernews Type {comment, story} will have its own channel.

Once the channel hits a threshold of n items the elastic bulk
processor inserts the JSON data into Elastic.

Initially, I thought I could go directly from fetch to elastic
bypassing Redis, but then realized that fetching the data off
the internet is INFINITELY slower than reading data from Redis.

And it also decouples these processes which can actually run
simultaneously anyway.  Having Redis sitting there as a cache
is what it does really well anyway and Redis always acts as
a really nice interface anyway.
