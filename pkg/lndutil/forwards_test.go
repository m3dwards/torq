package lndutil

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/cockroachdb/errors"
	"github.com/jmoiron/sqlx"
	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/lncapital/torq/testutil"
	"github.com/mixer/clock"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"regexp"
	"testing"
	"time"
)

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

// TestSubscribeForwardingIntervals tests that the forwarding
// history is fetched at intervals.
func TestSubscribeForwardingIntervals(t *testing.T) {

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer mockDB.Close()
	db := sqlx.NewDb(mockDB, "sqlmock")

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
				TimestampNs: uint64(c.Now().UnixNano()) + 1000000000,
			},
		},
		LastOffsetIndex: 0,
	}

	// Start subscribing in a goroutine to allow the test to continue simulating time through the
	// mocked time object.
	errs.Go(func() error {
		err = SubscribeForwardingEvents(ctx, mclient, db, &opt)
		if err != nil {
			t.Fatal(errors.Wrapf(err, "SubscribeForwardingEvents(%v, %v, %v, %v)", ctx,
				mclient, db, &opt))
		}
		return nil
	})

	// Simulate passing intervals

	numbTicks := 3
	for i := 0; i < numbTicks; i++ {
		mock.ExpectQuery("SELECT time_ns FROM forward ORDER BY time_ns DESC LIMIT 1;").
			WillReturnRows(mock.NewRows([]string{"time_ns"}).AddRow(0))

		mock.ExpectBegin()
		for _, event := range mclient.ForwardingEvents {
			mock.ExpectExec(regexp.QuoteMeta(querySfwh)).WithArgs(
				convMicro(event.TimestampNs), event.TimestampNs, event.FeeMsat,
				event.ChanIdIn, event.ChanIdOut, event.AmtInMsat,
				event.AmtOutMsat).WillReturnResult(sqlmock.NewResult(1, 1))
		}
		mock.ExpectCommit()

		c.AddTime(mockTickerInterval)

	}

	// Give the goroutine time to act on the mocked time interval
	time.Sleep(1000 * time.Millisecond)

	testutil.Given(t, "Given the need to test fwd events subscriptions.")

	testutil.WhenF(t, "When checking that the loop repeats.")

	err = mock.ExpectationsWereMet()
	if err != nil {
		testutil.Errorf(t, "We should see the database be queried %d times : %v", numbTicks,
			err)
	} else {
		testutil.Successf(t, "We should see the database be queried %d times",
			numbTicks)
	}

	time.Sleep(1000 * time.Millisecond)
	// Stop subscribing by canceling the context and ticking to the next iteration.
	stopSubFwE()
	c.AddTime(mockTickerInterval)

	// Check for potential errors from the goroutine (SubscribeForwardingEvents)
	err = errs.Wait()
	if err != nil {
		t.Fatal(err)
	}

}
