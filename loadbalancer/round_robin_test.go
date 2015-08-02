package loadbalancer_test

import (
	"reflect"
	"testing"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/loadbalancer"
	"github.com/go-kit/kit/loadbalancer/static"
	"golang.org/x/net/context"
)

func TestRoundRobinDistribution(t *testing.T) {
	var (
		ctx       = context.Background()
		counts    = []int{0, 0, 0}
		endpoints = []endpoint.Endpoint{
			func(context.Context, interface{}) (interface{}, error) { counts[0]++; return struct{}{}, nil },
			func(context.Context, interface{}) (interface{}, error) { counts[1]++; return struct{}{}, nil },
			func(context.Context, interface{}) (interface{}, error) { counts[2]++; return struct{}{}, nil },
		}
	)

	lb := loadbalancer.NewRoundRobin(static.Publisher(endpoints))

	for i, want := range [][]int{
		{1, 0, 0},
		{1, 1, 0},
		{1, 1, 1},
		{2, 1, 1},
		{2, 2, 1},
		{2, 2, 2},
		{3, 2, 2},
	} {
		e, err := lb.Endpoint()
		if err != nil {
			t.Fatal(err)
		}
		e(ctx, struct{}{})
		if have := counts; !reflect.DeepEqual(want, have) {
			t.Fatalf("%d: want %v, have %v", i, want, have)
		}

	}
}

func TestRoundRobinBadPublisher(t *testing.T) {
	t.Skip("TODO")
}
