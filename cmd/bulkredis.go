package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/olivere/elastic"
	"github.com/stormasm/elastichacker/redisc"
	"golang.org/x/sync/errgroup"
)

// index    Elasticsearch index name
// typ      Elasticsearch type name
// bulkSize Number of documents to collect before committing

func insert(index string, typ string, bulkSize int) {

	// Do a trace log
	tracelog := log.New(os.Stdout, "", 0)
	client, err := elastic.NewClient(elastic.SetTraceLog(tracelog))
	// Or with nothing...
	// client, err := elastic.NewClient()

	if err != nil {
		// Handle error
		log.Fatal(err)
	}

	// Setup a group of goroutines from the excellent errgroup package
	g, ctx := errgroup.WithContext(context.TODO())

	// The first goroutine will emit documents and send it to the second goroutine
	// via the docsc channel.
	// The second Goroutine will simply bulk insert the documents.
	docsc := make(chan redisc.Datum)

	// Goroutine to create documents
	g.Go(func() error {
		defer close(docsc)

		// Eventually one could pass in an array of strings which
		// would be the keys one can pull from redis...
		go redisc.Hscan("story", docsc)
		go redisc.Hscan("comment", docsc)

		// This is a hack, need to fix it...
		var input string
		fmt.Scanln(&input)

		return nil
	})

	// Second goroutine will consume the documents sent from the first and bulk insert into ES
	g.Go(func() error {
		bulk := client.Bulk().Index(index).Type(typ)

		for d := range docsc {
			// Enqueue the document
			myid := d.Id
			mydoc := d.Json
			bulk.Add(elastic.NewBulkStringRequest().Id(myid).SetSource(mydoc))

			if bulk.NumberOfActions() >= bulkSize {
				// Commit
				res, err := bulk.Do(ctx)
				if err != nil {
					return err
				}
				if res.Errors {
					// Look up the failed documents with res.Failed(), and e.g. recommit
					return errors.New("bulk commit failed")
				}
				// "bulk" is reset after Do, so you can reuse it
			}

			select {
			default:
			case <-ctx.Done():
				return ctx.Err()
			}
		}

		// Commit the final batch before exiting
		if bulk.NumberOfActions() > 0 {
			_, err = bulk.Do(ctx)
			if err != nil {
				return err
			}
		}
		return nil
	})

	// Wait until all goroutines are finished
	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	insert("warehouse", "product", 3)
}
