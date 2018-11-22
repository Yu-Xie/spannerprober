package spanneroperations

import (
	"context"
	"log"
	"regexp"
	"time"

	"github.com/Yu-Xie/spannerprober/consts"
	"github.com/Yu-Xie/spannerprober/spannerclient"
	"google.golang.org/genproto/googleapis/spanner/admin/database/v1"
)

var TestTable = testTable("TestTable")

type Table struct {
	TableName, CreateSQL, DropSQL string
	Cols                          []string
}

func testTable(tableName string) *Table {
	return &Table{
		TableName: tableName,
		CreateSQL: `
CREATE TABLE ` + tableName + ` (UUID STRING(36) NOT NULL,
    IndexUUID STRING(36),
	Placeholder STRING(MAX),
) PRIMARY KEY (UUID)`,
		DropSQL: `DROP TABLE ` + tableName,
		Cols:    []string{"UUID", "IndexUUID", "Placeholder"},
	}
}

func CreateTable() {
	adminClient, err := spannerclient.NewAdminClient()
	if err != nil {
		log.Fatalln(err.Error())
	}
	matches := regexp.MustCompile("^(.*)/databases/(.*)$").FindStringSubmatch(consts.SpannerFullName)
	if matches == nil || len(matches) != 3 {
		log.Fatalf("Invalid database id %s", consts.SpannerFullName)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*10)
	defer cancel()

	op, err := adminClient.CreateDatabase(ctx, &database.CreateDatabaseRequest{
		Parent:          matches[1],
		CreateStatement: "CREATE DATABASE `" + matches[2] + "`",
		ExtraStatements: []string{TestTable.CreateSQL},
	})
	if err != nil {
		log.Fatalln(err.Error())
	}
	if _, err := op.Wait(ctx); err != nil {
		log.Fatalln(err.Error())
	}
	log.Printf("Created database [%s]\n", consts.SpannerDatabaseName)
}
