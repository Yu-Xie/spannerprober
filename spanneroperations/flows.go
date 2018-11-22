package spanneroperations

import (
	"context"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/spanner"
	"github.com/Yu-Xie/spannerprober/consts"
	"github.com/google/uuid"
	"go.opencensus.io/trace"
)

const (
	roTransactionTimeout = time.Second * 10
	rwTransactionTimeout = time.Second * 10
)

var Flow1 = func(client *spanner.Client) {
	name := "flow1"
	ctx := context.Background()

	// add a parent span for the flow
	flowCtx, flowSpan := trace.StartSpan(ctx, flowSpanName(name))
	defer flowSpan.End()

	// warm up sessions by a dummy read
	log.Println("Start warming up sessions")
	runReadRowWithSpan("warm-up", "dummy", flowCtx, client)

	// blind write to insert rows
	log.Println("Start pre-populate writes")
	keys := prepopulate(flowCtx, client, 2)

	// read only
	log.Println("Start read only operations")
	for _, key := range keys {
		// ignore results and errors
		_, err := runReadRowWithSpan("read-only", key, flowCtx, client)
		if err != nil {
			log.Println(err.Error())
		}
	}

	// one read-write with 1 write after 2 reads
	log.Println("Start read write transactions")
	rwTxn := func(context context.Context, transaction *spanner.ReadWriteTransaction) error {
		for _, key := range keys {
			_, err := transaction.ReadRow(ctx, TestTable.TableName, spanner.Key{key}, TestTable.Cols)
			if err != nil {
				return err
			}
		}
		writeRowToBuffer(transaction, []interface{}{uuid.New().String(), uuid.New().String(), "payload"})
		return nil
	}
	runRWTransactionWithSpan("read-write", flowCtx, rwTxn, client)
}

func blindWrite(name string, ctx context.Context, client *spanner.Client) (string, error) {
	key := uuid.New().String()
	indexKey := uuid.New().String()
	transactionFunc := func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		return writeRowToBuffer(txn, []interface{}{key, indexKey, "payload"})
	}
	return key, runRWTransactionWithSpan(name, ctx, transactionFunc, client)
}

func runRWTransactionWithSpan(
	spanName string,
	parentCtx context.Context,
	f func(context.Context, *spanner.ReadWriteTransaction) error,
	client *spanner.Client) error {
	// span for the write
	ctx, span := trace.StartSpan(parentCtx, spanName)
	defer span.End()

	// start the write
	ctx, cancel := context.WithTimeout(ctx, rwTransactionTimeout)
	defer cancel()
	timestamp, err := client.ReadWriteTransaction(ctx, f)
	if err != nil {
		log.Println(err.Error())
		span.AddAttributes(trace.StringAttribute("msg", err.Error()))
	} else {
		span.AddAttributes(trace.StringAttribute("commit_time", timestamp.String()))
	}
	return err
}

func writeRowToBuffer(txn *spanner.ReadWriteTransaction, vals []interface{}) error {
	return txn.BufferWrite([]*spanner.Mutation{
		spanner.InsertOrUpdate(TestTable.TableName,
			TestTable.Cols,
			vals),
	})
}

func prepopulate(ctx context.Context, client *spanner.Client, numRows int) []string {
	uuids := make([]string, 0, numRows)
	for i := 0; i < numRows; {
		if key, err := blindWrite("blind-write", ctx, client); err == nil {
			i++
			uuids = append(uuids, key)
		}
	}
	return uuids
}

func runReadRowWithSpan(spanName string, key string, parentCtx context.Context, client *spanner.Client) (*spanner.Row, error) {
	ctx, span := trace.StartSpan(parentCtx, spanName)
	defer span.End()

	ctx, cancel := context.WithTimeout(ctx, roTransactionTimeout)
	defer cancel()

	row, err := client.Single().ReadRow(ctx, TestTable.TableName, spanner.Key{key}, TestTable.Cols)
	if err != nil {
		span.AddAttributes(trace.StringAttribute("msg", err.Error()))
	}
	return row, err
}

func flowSpanName(flowName string) string {
	testName, _ := consts.TestName()
	return fmt.Sprintf(`flow_%s_%s`, flowName, testName)
}
