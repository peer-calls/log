package log_test

import (
	"fmt"
	"os"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/peer-calls/log"
	"github.com/stretchr/testify/assert"
)

type testWriter struct {
	mockErr error
	b       strings.Builder
}

func newTestWriter() *testWriter {
	return &testWriter{}
}

func (t *testWriter) Write(b []byte) (int, error) {
	if t.mockErr != nil {
		return 0, t.mockErr
	}

	return t.b.Write(b)
}

func (t *testWriter) String() string {
	return t.b.String()
}

type testFormatter struct {
	*log.StringFormatter
	mockErr error
}

func newTestFormatter() *testFormatter {
	return &testFormatter{
		StringFormatter: log.NewStringFormatter(log.StringFormatterParams{
			DateLayout:               "-",
			DisableContextKeySorting: false,
		}),
	}
}

func (f *testFormatter) Format(message log.Message) ([]byte, error) {
	if f.mockErr != nil {
		return nil, f.mockErr
	}

	return f.StringFormatter.Format(message)
}

var errTest = fmt.Errorf("test err")

func TestLogger_Namespace(t *testing.T) {
	t.Parallel()

	log := log.New().WithNamespace("test").WithNamespaceAppended("test2")

	assert.Equal(t, "test:test2", log.Namespace())
}

func TestLogger(t *testing.T) {
	t.Parallel()

	type testEntry struct {
		namespace string
		level     log.Level
		message   string
		err       error
		ctx       log.Ctx
	}

	type testCase struct {
		config           string
		ctx              log.Ctx
		entries          []testEntry
		mockWriterErr    error
		mockFormatterErr error
		wantErr          error
		wantResult       string
	}

	testCases := []testCase{
		{
			config: "",
			entries: []testEntry{
				{"a", log.LevelInfo, "test", nil, nil},
			},
			wantResult: "",
		},
		{
			config: "a:b",
			entries: []testEntry{
				{"a:b", log.LevelInfo, "test", nil, nil},
			},
			wantResult: "-  info [                 a:b] test\n",
		},
		{
			config: "a",
			entries: []testEntry{
				{"a", log.LevelDebug, "test", nil, nil},
			},
		},
		{
			config: "a",
			entries: []testEntry{
				{"b", log.LevelInfo, "test", nil, nil},
			},
		},
		{
			config: "a:b:debug",
			entries: []testEntry{
				{"a:b", log.LevelDebug, "test", nil, nil},
			},
			wantResult: "- debug [                 a:b] test\n",
		},
		{
			config: "a:b:debug",
			entries: []testEntry{
				{"a:b", log.LevelDebug, "test", nil, nil},
			},
			mockWriterErr: errTest,
			wantErr:       fmt.Errorf("log write error: %w", errTest),
		},
		{
			config: "a:b:debug",
			entries: []testEntry{
				{"a:b", log.LevelDebug, "test", nil, nil},
			},
			mockFormatterErr: errTest,
			wantErr:          fmt.Errorf("log format error: %w", errTest),
		},
		{
			config: "a:b:debug",
			ctx: log.Ctx{
				"k1": "v1",
				"k2": "v2",
			},
			entries: []testEntry{
				{"a:b", log.LevelDebug, "test", nil, log.Ctx{"k2": "v3"}},
			},
			wantResult: "- debug [                 a:b] test k1=v1 k2=v3\n",
		},
		{
			config: "*:b:trace",
			entries: []testEntry{
				{"a:b", log.LevelTrace, "test1", nil, nil},
				{"a:c", log.LevelTrace, "test2", nil, nil},
			},
			wantResult: "- trace [                 a:b] test1\n",
		},
		{
			config: ":debug",
			entries: []testEntry{
				{"a:b", log.LevelDebug, "test1", nil, nil},
				{"a:c", log.LevelTrace, "test2", nil, nil},
			},
			wantResult: "- debug [                 a:b] test1\n",
		},
		{
			config: "a:b:trace,c:d",
			ctx: log.Ctx{
				"k1": "v1",
				"k2": "v2",
			},
			entries: []testEntry{
				{"a:b", log.LevelTrace, "test1", nil, log.Ctx{"k2": "v3"}},
				{"a:b", log.LevelDebug, "test2", nil, log.Ctx{"k3": "v3"}},
				{"a:b", log.LevelInfo, "test3", nil, log.Ctx{"k4": "v4"}},
				{"a:b", log.LevelWarn, "test4", nil, log.Ctx{"k5": "v5"}},
				{"a:b", log.LevelError, "", errTest, log.Ctx{"k6": "v6"}},
				{"a:b", log.LevelError, "err msg", errTest, log.Ctx{"k7": "v7"}},
				{"a:b", log.LevelError, "err msg", nil, log.Ctx{"k8": "v8"}},
				{"c:d", log.LevelTrace, "test1", nil, log.Ctx{"k2": "v3"}},
				{"c:d", log.LevelDebug, "test2", nil, log.Ctx{"k3": "v3"}},
				{"c:d", log.LevelInfo, "test3", nil, log.Ctx{"k4": "v4"}},
				{"c:d", log.LevelWarn, "test4", nil, log.Ctx{"k5": "v5"}},
				{"c:d", log.LevelError, "", errTest, log.Ctx{"k6": "v6"}},
				{"e:f", log.LevelTrace, "test1", nil, log.Ctx{"k2": "v3"}},
				{"e:f", log.LevelDebug, "test2", nil, log.Ctx{"k3": "v3"}},
				{"e:f", log.LevelInfo, "test3", nil, log.Ctx{"k4": "v4"}},
				{"e:f", log.LevelWarn, "test4", nil, log.Ctx{"k5": "v5"}},
				{"e:f", log.LevelError, "", errTest, log.Ctx{"k6": "v6"}},
			},
			wantResult: `- trace [                 a:b] test1 k1=v1 k2=v3
- debug [                 a:b] test2 k1=v1 k2=v2 k3=v3
-  info [                 a:b] test3 k1=v1 k2=v2 k4=v4
-  warn [                 a:b] test4 k1=v1 k2=v2 k5=v5
- error [                 a:b] test err k1=v1 k2=v2 k6=v6
- error [                 a:b] err msg: test err k1=v1 k2=v2 k7=v7
- error [                 a:b] err msg k1=v1 k2=v2 k8=v8
-  info [                 c:d] test3 k1=v1 k2=v2 k4=v4
-  warn [                 c:d] test4 k1=v1 k2=v2 k5=v5
- error [                 c:d] test err k1=v1 k2=v2 k6=v6
`,
		},
	}

	for i, tc := range testCases {
		descr := fmt.Sprintf("test case: %d", i)

		w := newTestWriter()

		formatter := newTestFormatter()

		root := log.New().WithConfig(log.NewConfigFromString(tc.config)).
			WithWriter(w).
			WithCtx(tc.ctx).
			WithFormatter(formatter)

		w.mockErr = tc.mockWriterErr
		formatter.mockErr = tc.mockFormatterErr

		var gotErr error

		for _, entry := range tc.entries {
			switch entry.level {
			case log.LevelError:
				_, gotErr = root.WithNamespace(entry.namespace).Error(entry.message, entry.err, entry.ctx)
			case log.LevelWarn:
				_, gotErr = root.WithNamespace(entry.namespace).Warn(entry.message, entry.ctx)
			case log.LevelInfo:
				_, gotErr = root.WithNamespace(entry.namespace).Info(entry.message, entry.ctx)
			case log.LevelDebug:
				_, gotErr = root.WithNamespace(entry.namespace).Debug(entry.message, entry.ctx)
			case log.LevelTrace:
				_, gotErr = root.WithNamespace(entry.namespace).Trace(entry.message, entry.ctx)
			case log.LevelDisabled:
				fallthrough
			case log.LevelUnknown:
				fallthrough
			default:
				panic(fmt.Sprintf("unexpected level: %s", entry.level))
			}
		}

		assert.Equal(t, tc.wantErr, gotErr, "%s: wantErr", descr)

		gotResult := w.String()

		assert.Equal(t, tc.wantResult, gotResult, "\n", "%s: wantResult", descr)
	}
}

func TestLogger_WithNamespaceAppended(t *testing.T) {
	t.Parallel()

	w := newTestWriter()

	logger := log.New().WithConfig(log.LevelInfo).
		WithNamespaceAppended("a").
		WithNamespaceAppended("b").
		WithWriter(w).
		WithFormatter(newTestFormatter())

	_, err := logger.Info("test1", nil)
	assert.NoError(t, err)

	_, err = logger.Trace("test2", nil)
	assert.NoError(t, err)

	gotStr := w.String()

	assert.Equal(t, "-  info [                 a:b] test1\n", gotStr)
}

func TestLogger_Ctx(t *testing.T) {
	t.Parallel()

	logger := log.New()

	assert.Equal(t, log.Ctx(nil), logger.Ctx())
	assert.Equal(t, log.Ctx{"a": "b"}, logger.WithCtx(log.Ctx{"a": "b"}).Ctx())
}

func TestNewFromEnv(t *testing.T) {
	t.Parallel()

	envKey := "TEST_LOG"

	old := os.Getenv(envKey)

	defer os.Setenv(envKey, old)

	os.Setenv(envKey, "**:a:trace,:info")

	logger := log.NewFromEnv(envKey)

	assert.Equal(t, true, logger.IsLevelEnabled(log.LevelInfo))
	assert.Equal(t, false, logger.IsLevelEnabled(log.LevelDebug))

	assert.Equal(t, true, logger.WithNamespace("a").IsLevelEnabled(log.LevelInfo))
	assert.Equal(t, true, logger.WithNamespace("a").IsLevelEnabled(log.LevelDebug))

	assert.Equal(t, true, logger.WithNamespace("c:b:a").IsLevelEnabled(log.LevelInfo))
	assert.Equal(t, true, logger.WithNamespace("c:b:a").IsLevelEnabled(log.LevelDebug))

	assert.Equal(t, true, logger.WithNamespace("b").IsLevelEnabled(log.LevelInfo))
	assert.Equal(t, false, logger.WithNamespace("b").IsLevelEnabled(log.LevelDebug))
}

func BenchmarkLogger_Disabled(b *testing.B) {
	logger := log.New().WithNamespace("test").WithConfig(log.LevelDisabled)

	var thread int64

	b.RunParallel(func(pb *testing.PB) {
		curThread := atomic.AddInt64(&thread, 1)

		var n int64

		for pb.Next() {
			curN := atomic.AddInt64(&n, 1)
			_, _ = logger.Info("benchmark", log.Ctx{"thread": curThread, "n": curN})
		}
	})
}

func BenchmarkLogger_Enabled(b *testing.B) {
	logger := log.New().WithNamespace("test").WithConfig(log.LevelInfo)

	var thread int64

	b.RunParallel(func(pb *testing.PB) {
		curThread := atomic.AddInt64(&thread, 1)

		var n int64

		for pb.Next() {
			curN := atomic.AddInt64(&n, 1)
			_, _ = logger.Info("benchmark", log.Ctx{"thread": curThread, "n": curN})
		}
	})
}

func BenchmarkLogger_EnabledWithoutSorting(b *testing.B) {
	logger := log.New().WithNamespace("test").WithConfig(log.LevelInfo).
		WithFormatter(log.NewStringFormatter(log.StringFormatterParams{
			DateLayout:               "",
			DisableContextKeySorting: true,
		}))

	var thread int64

	b.RunParallel(func(pb *testing.PB) {
		curThread := atomic.AddInt64(&thread, 1)

		var n int64

		for pb.Next() {
			curN := atomic.AddInt64(&n, 1)
			_, _ = logger.Info("benchmark", log.Ctx{"thread": curThread, "n": curN})
		}
	})
}
