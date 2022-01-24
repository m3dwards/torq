package lndutil

import (
	"context"
	"fmt"
	"github.com/cockroachdb/errors"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/lncapital/torq/migrations"
	"github.com/lncapital/torq/testutil"
	"github.com/mixer/clock"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"log"
	"os"
	"testing"
	"time"
)

var db *sqlx.DB

func TestMain(m *testing.M) {

	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "timescale/timescaledb",
		Tag:        "latest-pg14",
		Env: []string{
			"POSTGRES_PASSWORD=torq",
			"POSTGRES_USER=torq",
			"POSTGRES_DB=torq",
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	hostAndPort := resource.GetHostPort("5432/tcp")

	databaseUrl := fmt.Sprintf("postgres://torq:torq@%s/torq?sslmode=disable", hostAndPort)

	log.Println("Connecting to database on url: ", databaseUrl)
	resource.Expire(120) // Tell docker to hard kill the container in 120 seconds

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	pool.MaxWait = 120 * time.Second
	if err = pool.Retry(func() error {
		db, err = sqlx.Connect("postgres", databaseUrl)
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	err = migrations.MigrateUp(databaseUrl)
	if err != nil {
		log.Fatalf("Could not migrate DB: %s", err)
	}

	//Run tests
	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

// mockLightningClientForwardingHistory is used to moc responses from GetNodeInfo
type mockLightningClientForwardingHistory struct {
	CustomMaxEvents  int32
	ForwardingEvents []*lnrpc.ForwardingEvent
	LastOffsetIndex  uint32
	Error            error
}

// TODO: Use fuzzy tests:
//	 https://go.dev/doc/fuzz/
//   https://go.dev/blog/fuzz-beta

func (c mockLightningClientForwardingHistory) ForwardingHistory(ctx context.Context,
	in *lnrpc.ForwardingHistoryRequest,
	opts ...grpc.CallOption) (*lnrpc.ForwardingHistoryResponse, error) {

	if c.Error != nil {
		return nil, c.Error
	}

	r := lnrpc.ForwardingHistoryResponse{
		ForwardingEvents: c.ForwardingEvents,
		LastOffsetIndex:  c.LastOffsetIndex,
	}

	return &r, nil
}

func TestFetchForwardingHistoryError(t *testing.T) {

	mClient := mockLightningClientForwardingHistory{
		Error: errors.New("Some error"),
	}

	ctx := context.Background()
	_, err := fetchForwardingHistory(ctx, mClient, 0, 1000)

	testutil.Given(t, "While fetching forwarding history")

	testutil.WhenF(t, "If lnrcp.ForwardingHistory returns an error.")

	if err != nil {
		testutil.Successf(t, "fetchForwardingHistory returns the error")
	} else {
		testutil.Errorf(t, "fetchForwardingHistory returns the error")
	}

}

func TestSubscribeForwardingEvents(t *testing.T) {

	ctx := context.Background()
	errs, ctx := errgroup.WithContext(ctx)
	ctx, stopSubFwE := context.WithCancel(ctx)
	c := clock.NewMockClock(time.Unix(0, 0))

	mockTickerInterval := 30 * time.Second
	me := 1000
	opt := FwhOptions{
		MaxEvents: &me,
		Tick:      c.Tick(mockTickerInterval),
	}

	mclient := mockLightningClientForwardingHistory{
		ForwardingEvents: []*lnrpc.ForwardingEvent{
			{
				ChanIdIn:    1234,
				ChanIdOut:   2345,
				AmtIn:       11,
				AmtOut:      10,
				Fee:         1,
				FeeMsat:     1000,
				AmtInMsat:   11000,
				AmtOutMsat:  10000,
				TimestampNs: uint64(c.Now().UnixNano()),
			},
			{
				ChanIdIn:    1234,
				ChanIdOut:   2345,
				AmtIn:       11,
				AmtOut:      10,
				Fee:         1,
				FeeMsat:     1000,
				AmtInMsat:   11000,
				AmtOutMsat:  10000,
				TimestampNs: uint64(c.Now().UnixNano()) + 500000000,
			},
			{ // Duplicate record used for testing
				ChanIdIn:    1234,
				ChanIdOut:   2345,
				AmtIn:       11,
				AmtOut:      10,
				Fee:         1,
				FeeMsat:     1000,
				AmtInMsat:   11000,
				AmtOutMsat:  10000,
				TimestampNs: uint64(c.Now().UnixNano()) + 500000000,
			},
			{
				ChanIdIn:    1234,
				ChanIdOut:   2345,
				AmtIn:       11,
				AmtOut:      10,
				Fee:         1,
				FeeMsat:     1000,
				AmtInMsat:   11000,
				AmtOutMsat:  10000,
				TimestampNs: uint64(c.Now().UnixNano()) + 1000000000,
			},
		},
		LastOffsetIndex: 0,
	}

	// Start subscribing in a goroutine to allow the test to continue simulating time through the
	// mocked time object.
	errs.Go(func() error {
		err := SubscribeForwardingEvents(ctx, mclient, db, &opt)
		if err != nil {
			t.Fatal(errors.Wrapf(err, "SubscribeForwardingEvents(%v, %v, %v, %v)", ctx,
				mclient, db, &opt))
		}
		return nil
	})

	// Simulate passing intervals
	numbTicks := 3
	for i := 0; i < numbTicks; i++ {

		c.AddTime(mockTickerInterval)

	}

	// Give the goroutine time to act on the mocked time interval
	time.Sleep(400 * time.Millisecond)

	testutil.Given(t, "While running SubscribeForwardingEvents")

	testutil.WhenF(t, "We need to check that fetchLastForwardTime returns the expected nanosecond.")
	{
		var expected uint64 = 1000000000
		returned, err := fetchLastForwardTime(db)
		switch {
		case err != nil:
			testutil.Fatalf(t, "We get an error: %v", err)
		case returned != expected:
			testutil.Errorf(t, " "+
				"We expected %d got %d", expected, returned)
		case returned == expected:
			testutil.Successf(t, "We got the expected nanosecond ")
		}
	}

	testutil.WhenF(t, "We need to check that storeForwardingHistory only stores unique records.")
	{
		expectedTotal := 4
		actualTotal := len(mclient.ForwardingEvents)

		if expectedTotal != actualTotal {
			testutil.Errorf(t, "We expected to mock %d ForwardingEvents but there where %",
				expectedTotal, actualTotal)
		}

		var expectedUnique = 3
		var returned int
		err := db.QueryRow("select count(*) from forward;").Scan(&returned)

		switch {
		case err != nil:
			testutil.Fatalf(t, "We get an error: %v", err)
		case returned != expectedUnique:
			testutil.Errorf(t, "We expected to store %d records but stored %d", expectedUnique,
				returned)
		case returned == expectedUnique:
			testutil.Successf(t, "We stored the expected number of records")
		}
	}

	// Stop subscribing by canceling the context and ticking to the next iteration.
	stopSubFwE()
	c.AddTime(mockTickerInterval)

	// Check for potential errors from the goroutine (SubscribeForwardingEvents)
	err := errs.Wait()
	if err != nil {
		t.Fatal(err)
	}

}
