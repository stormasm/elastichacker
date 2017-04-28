package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/olivere/elastic"
	"github.com/stormasm/elastichacker/redisc"
	"golang.org/x/sync/errgroup"
)

func main() {
	var (
		index    = flag.String("index", "", "Elasticsearch index name")
		typ      = flag.String("type", "", "Elasticsearch type name")
		n        = flag.Int("n", 0, "Number of documents to bulk insert")
		bulkSize = flag.Int("bulk-size", 0, "Number of documents to collect before committing")
	)
	flag.Parse()
	log.SetFlags(0)
	rand.Seed(time.Now().UnixNano())

	if *index == "" {
		log.Fatal("missing index parameter")
	}
	if *typ == "" {
		log.Fatal("missing type parameter")
	}
	if *n <= 0 {
		log.Fatal("n must be a positive number")
	}
	if *bulkSize <= 0 {
		log.Fatal("bulk-size must be a positive number")
	}

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

	// var newStory chan redisc.Datum = make(chan redisc.Datum, 100)

	docsc := make(chan redisc.Datum)

	// Goroutine to create documents
	g.Go(func() error {
		defer close(docsc)

		for i := 0; i < *n; i++ {

			// Construct the json string
			jsonstr := `{"user":"olivere","city":"santafe","age":56}`

			d := redisc.Datum{
				Hnid: strconv.Itoa(i),
				Json: jsonstr,
			}

			// Send over to 2nd goroutine, or cancel
			select {
			case docsc <- d:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
		return nil
	})

	// Second goroutine will consume the documents sent from the first and bulk insert into ES
	g.Go(func() error {
		bulk := client.Bulk().Index(*index).Type(*typ)

		//		count := 0

		for d := range docsc {
			// Enqueue the document
			countstr := d.Hnid
			mydoc := d.Json
			bulk.Add(elastic.NewBulkStringRequest().Id(countstr).SetSource(mydoc))

			//			count = count + 1

			if bulk.NumberOfActions() >= *bulkSize {
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
