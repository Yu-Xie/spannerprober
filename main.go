package main

import (
	"log"
	"os"
	"sync"

	"github.com/Yu-Xie/spannerprober/consts"
	"github.com/Yu-Xie/spannerprober/spannerclient"
	"github.com/Yu-Xie/spannerprober/spanneroperations"
	"github.com/Yu-Xie/spannerprober/telemetry"
)

func main() {
	if os.Getenv("CREATE") == "true" {
		spanneroperations.CreateTable()
		return
	}

	testName, err := consts.TestName()
	if err != nil {
		log.Fatalln(err.Error())
	}
	log.Println(testName)

	log.Println("Initializing tracing...")
	closer := telemetry.InitializeTracing(consts.ProjectName, 1.0)

	dataClient, err := spannerclient.NewDataClient(consts.SpannerFullName)
	if err != nil {
		log.Fatalln(err.Error())
	}

	spanneroperations.Flow1(dataClient)

	closer()
	var wg sync.WaitGroup
	wg.Wait()
}
