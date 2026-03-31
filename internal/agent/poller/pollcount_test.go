package poller

import (
	"testing"

	models "github.com/mikhailpashkov/metrics/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPollCountPoller_GetMetrics(t *testing.T) {
	type fields struct {
		count int64
	}
	type want struct {
		mType string
		mName string
		delta *int64
		value *float64
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "count 0",
			fields: fields{
				count: 0,
			},
			want: want{
				mType: models.Counter,
				mName: "custom.PollCount",
				delta: func(val int64) *int64 { return &val }(0),
				value: nil,
			},
		},
		{
			name: "count 10",
			fields: fields{
				count: 10,
			},
			want: want{
				mType: models.Counter,
				mName: "custom.PollCount",
				delta: func(val int64) *int64 { return &val }(0),
				value: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPollCountPoller()
			got, err := p.GetMetrics()
			if err != nil {
				t.Errorf("GetMetrics() error = %v", err)
				return
			}
			assert.Len(t, got, 1, "PollCountPoller.GetMetrics() should return only 1 metric")
			assert.Equal(t, tt.want.mType, got[0].Type)
			assert.Equal(t, tt.want.mName, got[0].Name)

			if tt.want.delta != nil {
				require.NotNil(t, got[0].Delta)
				assert.Equal(t, *tt.want.delta, *got[0].Delta)
			}

			if tt.want.value != nil {
				require.NotNil(t, got[0].Value)
				assert.Equal(t, *tt.want.value, *got[0].Value)
			}
		})
	}
}

func TestPollCountPoller_IncrementCount(t *testing.T) {
	type fields struct {
		count    int64
		incTimes int64
	}
	type want struct {
		delta *int64
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "increment 1 time from 0",
			fields: fields{
				count:    0,
				incTimes: 1,
			},
			want: want{
				delta: func(val int64) *int64 { return &val }(1),
			},
		},
		{
			name: "increment 10 times from 0",
			fields: fields{
				count:    0,
				incTimes: 10,
			},
			want: want{
				delta: func(val int64) *int64 { return &val }(10),
			},
		},
		{
			name: "increment 1 time from 10",
			fields: fields{
				count:    10,
				incTimes: 1,
			},
			want: want{
				delta: func(val int64) *int64 { return &val }(11),
			},
		},
		{
			name: "increment 10 times from 10",
			fields: fields{
				count:    10,
				incTimes: 10,
			},
			want: want{
				delta: func(val int64) *int64 { return &val }(20),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPollCountPoller()
			p.count = tt.fields.count
			for range tt.fields.incTimes {
				p.IncrementCount()
			}
			got, err := p.GetMetrics()
			if err != nil {
				t.Errorf("GetMetrics() error = %v", err)
				return
			}
			assert.Len(t, got, 1, "PollCountPoller.GetMetrics() should return only 1 metric")
			require.NotNil(t, got[0].Delta)
			assert.Equal(t, *tt.want.delta, *got[0].Delta)
		})
	}
}

func TestPollCountPoller_ResetCount(t *testing.T) {
	type fields struct {
		count int64
	}
	type want struct {
		delta *int64
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "reset from 0",
			fields: fields{
				count: 0,
			},
			want: want{
				delta: func(val int64) *int64 { return &val }(0),
			},
		},
		{
			name: "reset from 10",
			fields: fields{
				count: 10,
			},
			want: want{
				delta: func(val int64) *int64 { return &val }(0),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPollCountPoller()
			p.count = tt.fields.count
			p.ResetCount()
			got, err := p.GetMetrics()
			if err != nil {
				t.Errorf("GetMetrics() error = %v", err)
				return
			}
			assert.Len(t, got, 1, "PollCountPoller.GetMetrics() should return only 1 metric")
			require.NotNil(t, got[0].Delta)
			assert.Equal(t, *tt.want.delta, *got[0].Delta)
		})
	}
}
