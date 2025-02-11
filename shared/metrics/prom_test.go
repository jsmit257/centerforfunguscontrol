package metrics

import (
	"context"
	"testing"

	"github.com/jsmit257/huautla/types"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func Test_NewHandler(*testing.T) {
	_ = NewHandler()
}

func Test_GetContextCID(t *testing.T) {
	t.Parallel()
	tcs := map[string]struct {
		ctx context.Context
		cid types.CID
	}{
		"happy_path": {
			ctx: context.WithValue(context.TODO(),
				Cid,
				types.CID("happy_path")),
			cid: "happy_path",
		},
		"null_attr": {
			ctx: context.TODO(),
			cid: "context has no cid attribute: context.todoCtx{emptyCtx:context.emptyCtx{}}",
		},
	}

	for name, tc := range tcs {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			cid := GetContextCID(tc.ctx)
			require.Equal(t, tc.cid, cid)
		})
	}
}

func Test_GetContextLog(t *testing.T) {
	t.Parallel()
	tcs := map[string]struct {
		ctx context.Context
	}{
		"happy_path": {
			ctx: context.WithValue(context.TODO(),
				Log,
				logrus.NewEntry(logrus.New())),
		},
		"null_attr": {
			ctx: context.TODO(),
		},
	}

	for name, tc := range tcs {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			l := GetContextLog(tc.ctx)
			require.NotNil(t, l.Error)
		})
	}
}

func Test_GetContextMetrics(t *testing.T) {
	t.Parallel()
	tcs := map[string]struct {
		ctx context.Context
	}{
		"happy_path": {
			ctx: context.WithValue(context.TODO(),
				Metrics,
				DataMetrics),
		},
		"null_attr": {
			ctx: context.TODO(),
		},
	}

	for name, tc := range tcs {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			m := GetContextMetrics(tc.ctx)
			require.NotNil(t, m.MustCurryWith)
		})
	}
}
