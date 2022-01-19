package lndutil

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/mixer/clock"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
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

// TestSubscribeForwardingEventsColdStartBellowLimit tests that the forwarding
// history is fetched at an interval and returns the correct ForwardingHistory.
func TestSubscribeForwardingEventsColdStartBellowLimit(t *testing.T) {

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer mockDB.Close()
	db := sqlx.NewDb(mockDB, "sqlmock")

	mockClient := mockLightningClientForwardingHistory{
		ForwardingEvents: []*lnrpc.ForwardingEvent{},
		LastOffsetIndex:  0,
	}

	c := clock.NewMockClock()
	ctx := context.Background()
	errs, ctx := errgroup.WithContext(ctx)
	ctx, stopSub := context.WithCancel(ctx)

	mockTickerInterval := 30 * time.Second
	me := 1000

	opt := FwhOptions{
		MaxEvents: &me,
		Tick:      c.Tick(mockTickerInterval),
	}

	// Start subscribing in a goroutine to allow the test to continue simulating time through the
	// mocked time object.
	errs.Go(func() error {
		err = SubscribeForwardingEvents(ctx, mockClient, db, &opt)
		if err != nil {
			t.Fatal(err)
		}
		return nil
	})

	// Simulate passing of 3 intervals
	for i := 0; i < 3; i++ {
		mock.ExpectQuery("SELECT time_ns FROM forward ORDER BY time_ns DESC LIMIT 1;").
			WillReturnRows(mock.NewRows([]string{"time_ns"}).AddRow(1))
		c.AddTime(mockTickerInterval)
	}

	// Give the goroutine time to act on the mocked time interval
	time.Sleep(1 * time.Second)

	// Stop subscribing by canceling the context and ticking to the next iteration.
	stopSub()
	c.AddTime(mockTickerInterval)

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Fatal(err)
	}

	err = errs.Wait()
	if err != nil {
		t.Fatal(err)
	}
}
