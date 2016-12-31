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
