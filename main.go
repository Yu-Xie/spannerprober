package main

import (
	"code.uber.internal/marketplace/spannerprober/spannerclient"
	"google.golang.org/genproto/googleapis/spanner/admin/database/v1"
	"context"
	"log"
	"time"
	"cloud.google.com/go/spanner"
	"sync"
	"github.com/google/uuid"
	"os"
)

func main() {
	log.Println(os.Environ())
	for {
		write()
		time.Sleep(time.Second * 5)
	}

	var wg sync.WaitGroup
	wg.Wait()

}

func write() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 10)
	defer cancel()

	client, err := spannerclient.NewDataClient()
	if err != nil {
		log.Printf("failed to create spanner client: %s", err.Error())
	}
	transactionFunc := func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		txn.BufferWrite([]*spanner.Mutation{
			spanner.InsertOrUpdate("TestTable",
				[]string{"UUID", "IndexUUID", "PlaceHolder"},
				[]interface{}{uuid.New().String(), "val2", "val3"}),
		})
		return nil
	}

	t, e := client.ReadWriteTransaction(ctx, transactionFunc)
	if e != nil {
		log.Printf("Commit Failed: %s", e.Error())
	}
	log.Printf("Commit Timestamp: %s", t.String())
}

func createTables(project, db string) error {
	ctx, cancel := context.WithTimeout(context.Background(), spannerclient.CreateSchemaTimeout)
	defer cancel()



	admin := spannerclient.NewAdminClient()
	op, err := admin.CreateDatabase(ctx, &database.CreateDatabaseRequest{
		Parent: project,
		CreateStatement: "CREATE DATABASE `" + db + "`",
		ExtraStatements: []string {
			spannerclient.CreateTable,
		},
	})
	if err != nil {
		return err
	}
	if _, err := op.Wait(ctx); err != nil {
		return err
	}
	log.Printf("Created database [%s]\n", db)
	return nil
}