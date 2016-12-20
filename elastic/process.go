// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"log"
	"os"
	"golang.org/x/net/context"
	elastic "gopkg.in/olivere/elastic.v5"
)

func Process_json_string(index string,json string) {
	errorlog := log.New(os.Stdout, "APP ", log.LstdFlags)

	// Obtain a client. You can also provide your own HTTP client here.
	client, err := elastic.NewClient(elastic.SetErrorLog(errorlog))
	if err != nil {
		// Handle error
		panic(err)
	}

	// Use the IndexExists service to check if a specified index exists.
	exists, err := client.IndexExists(index).Do(context.Background())
	if err != nil {
		// Handle error
		panic(err)
	}
	if !exists {
		// Create a new index.
		createIndex, err := client.CreateIndex(index).Do(context.Background())
		if err != nil {
			// Handle error
			panic(err)
		}
		if !createIndex.Acknowledged {
			// Not acknowledged
		}
	}

	_, err = client.Index().
		Index(index).
		Type("tweet").
		Id("9").
		BodyString(json).
		Do(context.Background())
	if err != nil {
		// Handle error
		panic(err)
	}

	// Flush to make sure the documents got written.
	_, err = client.Flush().Index(index).Do(context.Background())
	if err != nil {
		panic(err)
	}
}

func Process_json_bytes(index string,byteArray []byte) {
	errorlog := log.New(os.Stdout, "APP ", log.LstdFlags)

	// Obtain a client. You can also provide your own HTTP client here.
	client, err := elastic.NewClient(elastic.SetErrorLog(errorlog))
	if err != nil {
		// Handle error
		panic(err)
	}

	// Use the IndexExists service to check if a specified index exists.
	exists, err := client.IndexExists(index).Do(context.Background())
	if err != nil {
		// Handle error
		panic(err)
	}
	if !exists {
		// Create a new index.
		createIndex, err := client.CreateIndex(index).Do(context.Background())
		if err != nil {
			// Handle error
			panic(err)
		}
		if !createIndex.Acknowledged {
			// Not acknowledged
		}
	}

	n := len(byteArray)
	s := string(byteArray[:n])

	_, err = client.Index().
		Index(index).
		Type("tweet").
		Id("8").
		BodyString(s).
		Do(context.Background())
	if err != nil {
		// Handle error
		panic(err)
	}

	// Flush to make sure the documents got written.
	_, err = client.Flush().Index(index).Do(context.Background())
	if err != nil {
		panic(err)
	}
}
