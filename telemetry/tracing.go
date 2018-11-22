package telemetry

import (
	"log"
	"time"

	"contrib.go.opencensus.io/exporter/stackdriver"
	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
)

// InitializeTracing initialize tracing for the process.
// sampleRate can be set to 1.0 for local debugging, but should
// be set to a very small value when running load testing benchmarks.
func InitializeTracing(projectName string, sampleRate float64) (closer func()) {
	se, err := stackdriver.NewExporter(stackdriver.Options{
		ProjectID:            projectName,
		BundleDelayThreshold: time.Second / 10,
		BundleCountThreshold: 10,
	})
	if err != nil {
		log.Fatalf("StatsExporter err: %v", err)
	}
	// Let's ensure that data is uploaded before the program exits

	// Enable tracing on the exporter
	trace.RegisterExporter(se)

	// Enable metrics collection
	view.RegisterExporter(se)

	// Views that we are interested in
	views := []*view.View{
		// latency
		ocgrpc.ClientRoundtripLatencyView,
		ocgrpc.ClientServerLatencyView,
		ocgrpc.ServerLatencyView,

		// others
		ocgrpc.ClientCompletedRPCsView,
		ocgrpc.ClientSentBytesPerRPCView,
		ocgrpc.ClientReceivedBytesPerRPCView,
	}

	// Register Views
	if err := view.Register(views...); err != nil {
		log.Fatalf("failed to register views: %s", err.Error())
	}

	// Enable the trace sampler.
	// We are always sampling for demo purposes only: it is very high
	// depending on the QPS, but sufficient for the purpose of this quick demo.
	// More realistically perhaps tracing 1 in 10,000 might be more useful
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.ProbabilitySampler(
		sampleRate,
	)})

	// Returns a closer which ensures that data is uploaded before the program exits
	return se.Flush
}
