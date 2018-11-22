package consts

import (
	"errors"
	"fmt"
	"log"
	"os"
)

const (
	ProjectName         = "horizon-spanner-benchmark"
	SpannerDatabaseName = "testdb"

	EnvSpannerName  = "SPANNER_NAME"
	EnvSourceRegion = "SOURCE_REGION"
)

var (
	SourceRegion = os.Getenv(EnvSourceRegion)
	SpannerName  = os.Getenv(EnvSpannerName)

	SpannerInstanceName = SpannerName + "test"
	SpannerLeaderRegion = getLeaderRegion(SpannerName)
	SpannerFullName     = fmt.Sprintf(
		"projects/%s/instances/%s/databases/%s",
		ProjectName,
		SpannerInstanceName,
		SpannerDatabaseName)
)

func TestName() (string, error) {
	if SpannerName == "" {
		return "", errors.New(fmt.Sprintf("env var %s required", EnvSpannerName))
	} else if SourceRegion == "" {
		return "", errors.New(fmt.Sprintf("env var %s required", EnvSourceRegion))
	}
	return fmt.Sprintf(`from_%s_to_%s_leader_%s`, SourceRegion, SpannerName, SpannerLeaderRegion), nil
}

func getLeaderRegion(name string) string {
	switch name {
	case "nam3":
		return "us-east4"
	case "nam6":
		return "us-central1"
	case "eur3":
		return "europe-west1"
	case "nam-eur-asia1":
		return "us-central1"
	default:
		log.Fatalf("invalid spanner region: %s\n", name)
	}
	return ""
}
