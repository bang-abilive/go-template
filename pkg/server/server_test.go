package server

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"ndinhbang/go-skeleton/pkg/config"

	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newTestServer builds a Server with default middlewares applied — ready for httptest.
func newTestServer(t *testing.T) *Server {
	t.Helper()
	srv := New(&config.ServerConfig{Port: 8080})
	srv.SetupMiddlewares()
	return srv
}

// ---------------------------------------------------------------------------
// Construction
// ---------------------------------------------------------------------------

func TestNew(t *testing.T) {
	t.Parallel()

	cfg := &config.ServerConfig{Port: 8080}
	srv := New(cfg)

	require.NotNil(t, srv, "New() must not return nil")
	assert.Same(t, cfg, srv.cfg, "New() must retain the exact config pointer")
	assert.NotNil(t, srv.echo, "New() must initialise the echo instance")
}

// ---------------------------------------------------------------------------
// Default middlewares — body limit, security headers, panic recovery
// ---------------------------------------------------------------------------

func TestSetupMiddlewares_Default(t *testing.T) {
	t.Parallel()

	srv := newTestServer(t)
	srv.echo.POST("/echo", func(c *echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})
	srv.echo.GET("/panic", func(_ *echo.Context) error {
		panic("intentional test panic")
	})

	tests := []struct {
		name       string
		method     string
		path       string
		body       string
		wantStatus int
		checkFn    func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		{
			name:       "small body is accepted",
			method:     http.MethodPost,
			path:       "/echo",
			body:       "hello",
			wantStatus: http.StatusOK,
		},
		{
			name:       "body over 1 MB is rejected with 413",
			method:     http.MethodPost,
			path:       "/echo",
			body:       strings.Repeat("a", int(MaxBodyLimitBytes)+1),
			wantStatus: http.StatusRequestEntityTooLarge,
		},
		{
			name:       "X-Content-Type-Options is nosniff",
			method:     http.MethodPost,
			path:       "/echo",
			body:       "safe",
			wantStatus: http.StatusOK,
			checkFn: func(t *testing.T, rec *httptest.ResponseRecorder) {
				t.Helper()
				assert.Equal(t, "nosniff", rec.Header().Get("X-Content-Type-Options"))
			},
		},
		{
			name:       "X-Frame-Options is SAMEORIGIN",
			method:     http.MethodPost,
			path:       "/echo",
			body:       "safe",
			wantStatus: http.StatusOK,
			checkFn: func(t *testing.T, rec *httptest.ResponseRecorder) {
				t.Helper()
				assert.Equal(t, "SAMEORIGIN", rec.Header().Get("X-Frame-Options"))
			},
		},
		{
			name:       "X-XSS-Protection header is present",
			method:     http.MethodPost,
			path:       "/echo",
			body:       "safe",
			wantStatus: http.StatusOK,
			checkFn: func(t *testing.T, rec *httptest.ResponseRecorder) {
				t.Helper()
				assert.NotEmpty(t, rec.Header().Get("X-XSS-Protection"))
			},
		},
		{
			name:       "panic inside handler is recovered and returns 500",
			method:     http.MethodGet,
			path:       "/panic",
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var bodyReader *strings.Reader
			if tt.body != "" {
				bodyReader = strings.NewReader(tt.body)
			} else {
				bodyReader = strings.NewReader("")
			}

			req := httptest.NewRequest(tt.method, tt.path, bodyReader)
			req.Header.Set(echo.HeaderContentType, "text/plain")
			rec := httptest.NewRecorder()

			srv.echo.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantStatus, rec.Code)
			if tt.checkFn != nil {
				tt.checkFn(t, rec)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// CORS middleware
// ---------------------------------------------------------------------------

func TestSetupMiddlewares_CORS(t *testing.T) {
	t.Parallel()

	srv := newTestServer(t)
	srv.echo.GET("/data", func(c *echo.Context) error {
		return c.String(http.StatusOK, "data")
	})

	tests := []struct {
		name                string
		method              string
		origin              string
		requestMethod       string
		wantStatus          int
		wantAllowOriginStar bool
	}{
		{
			name:                "simple GET with Origin receives wildcard allow-origin",
			method:              http.MethodGet,
			origin:              "https://example.com",
			wantStatus:          http.StatusOK,
			wantAllowOriginStar: true,
		},
		{
			name:          "preflight OPTIONS returns 204 with CORS headers",
			method:        http.MethodOptions,
			origin:        "https://example.com",
			requestMethod: http.MethodPost,
			wantStatus:    http.StatusNoContent,
		},
		{
			name:       "request without Origin header passes through unchanged",
			method:     http.MethodGet,
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(tt.method, "/data", nil)
			if tt.origin != "" {
				req.Header.Set(echo.HeaderOrigin, tt.origin)
			}
			if tt.requestMethod != "" {
				req.Header.Set(echo.HeaderAccessControlRequestMethod, tt.requestMethod)
			}
			rec := httptest.NewRecorder()

			srv.echo.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantStatus, rec.Code)
			if tt.wantAllowOriginStar {
				assert.Equal(t, "*", rec.Header().Get(echo.HeaderAccessControlAllowOrigin))
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Custom middlewares
// ---------------------------------------------------------------------------

func TestSetupMiddlewares_Custom(t *testing.T) {
	t.Parallel()

	srv := New(&config.ServerConfig{Port: 8080})

	var called atomic.Int32
	tracer := func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			called.Add(1)
			return next(c)
		}
	}

	srv.SetupMiddlewares(tracer)
	srv.echo.GET("/custom", func(c *echo.Context) error {
		return c.String(http.StatusOK, "custom")
	})

	req := httptest.NewRequest(http.MethodGet, "/custom", nil)
	rec := httptest.NewRecorder()
	srv.echo.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	assert.EqualValues(t, 1, called.Load(), "custom middleware must be called exactly once per request")
}

func TestSetupMiddlewares_CustomSkipsDefaults(t *testing.T) {
	t.Parallel()

	srv := New(&config.ServerConfig{Port: 8080})
	noop := func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error { return next(c) }
	}
	srv.SetupMiddlewares(noop)
	srv.echo.POST("/large", func(c *echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	oversizedBody := strings.Repeat("a", int(MaxBodyLimitBytes)+1)
	req := httptest.NewRequest(http.MethodPost, "/large", strings.NewReader(oversizedBody))
	req.Header.Set(echo.HeaderContentType, "text/plain")
	rec := httptest.NewRecorder()
	srv.echo.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code,
		"body-limit must not be enforced when custom middlewares override defaults")
	assert.Empty(t, rec.Header().Get("X-Content-Type-Options"),
		"secure middleware must not be present when using custom middlewares")
}

// ---------------------------------------------------------------------------
// Routes
// ---------------------------------------------------------------------------

func TestSetupRoutes(t *testing.T) {
	t.Parallel()

	srv := newTestServer(t)
	srv.SetupRoutes()

	tests := []struct {
		name       string
		method     string
		path       string
		wantStatus int
		wantBody   string
	}{
		{
			name:       "root endpoint returns welcome message",
			method:     http.MethodGet,
			path:       "/",
			wantStatus: http.StatusOK,
			wantBody:   "Welcome to the API!",
		},
		{
			name:       "health endpoint returns OK",
			method:     http.MethodGet,
			path:       "/health",
			wantStatus: http.StatusOK,
			wantBody:   "OK",
		},
		{
			name:       "hello endpoint returns greeting",
			method:     http.MethodGet,
			path:       "/api/v1/hello",
			wantStatus: http.StatusOK,
			wantBody:   "Hello, World!",
		},
		{
			name:       "unknown path returns 404",
			method:     http.MethodGet,
			path:       "/does-not-exist",
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(tt.method, tt.path, nil)
			rec := httptest.NewRecorder()
			srv.echo.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantStatus, rec.Code)
			if tt.wantBody != "" {
				assert.Contains(t, rec.Body.String(), tt.wantBody)
			}
		})
	}
}

func TestSetupRoutes_MethodNotAllowed(t *testing.T) {
	t.Parallel()

	srv := newTestServer(t)
	srv.SetupRoutes()

	tests := []struct {
		name   string
		method string
		path   string
	}{
		{"POST to root", http.MethodPost, "/"},
		{"DELETE to health", http.MethodDelete, "/health"},
		{"PUT to hello", http.MethodPut, "/api/v1/hello"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(tt.method, tt.path, nil)
			rec := httptest.NewRecorder()
			srv.echo.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
		})
	}
}

// ---------------------------------------------------------------------------
// Concurrent request safety
// ---------------------------------------------------------------------------

func TestConcurrentRequests(t *testing.T) {
	t.Parallel()

	srv := newTestServer(t)
	srv.SetupRoutes()

	const workers = 50
	var (
		wg      sync.WaitGroup
		success atomic.Int64
	)

	wg.Add(workers)
	for range workers {
		go func() {
			defer wg.Done()
			req := httptest.NewRequest(http.MethodGet, "/health", nil)
			rec := httptest.NewRecorder()
			srv.echo.ServeHTTP(rec, req)
			if rec.Code == http.StatusOK {
				success.Add(1)
			}
		}()
	}
	wg.Wait()

	assert.EqualValues(t, workers, success.Load(),
		"all concurrent requests must receive HTTP 200")
}

// ---------------------------------------------------------------------------
// Config helpers
// ---------------------------------------------------------------------------

func TestServerConfig_Address(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		port uint16
		want string
	}{
		{"standard HTTP port", 80, ":80"},
		{"standard HTTPS port", 443, ":443"},
		{"default app port", 8080, ":8080"},
		{"custom port", 9090, ":9090"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cfg := config.ServerConfig{Port: tt.port}
			assert.Equal(t, tt.want, cfg.ServerAddress())
		})
	}
}

// ---------------------------------------------------------------------------
// Start — brief lifecycle (line coverage)
// ---------------------------------------------------------------------------

// TestStart_BriefLifecycle verifies that Start() binds and returns nil on
// graceful shutdown. It does NOT test in-flight request behaviour — see
// TestServe_GracefulShutdown_CompletesInFlightRequests for that.
func TestStart_BriefLifecycle(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // pre-cancel: bind → immediate graceful shutdown → nil

	srv := New(&config.ServerConfig{Port: 0})
	require.NoError(t, srv.Start(ctx))
}

// ---------------------------------------------------------------------------
// Serve — graceful shutdown behaviour
// ---------------------------------------------------------------------------

func TestServe_GracefulShutdown_CompletesInFlightRequests(t *testing.T) {
	t.Parallel()

	// Test flow:
	//  1. Test runner starts srv.Serve in a goroutine.
	//  2. Server binds port 0 → sends actual address via addrCh.
	//  3. Test runner creates client and fires GET /slow.
	//  4. Client sends request, blocks waiting for response.
	//  5. Handler starts executing → closes handlerStarted ("I'm running").
	//  6. Test runner receives signal → immediately calls cancel().
	//  7. Server receives shutdown signal but sees /slow in-flight → waits.
	//  8. Handler finishes time.Sleep → returns 200 OK.
	//  9. Server sees no more in-flight requests → Serve() returns nil.
	// 10. Test runner verifies response status and clean shutdown.

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	addrCh := make(chan net.Addr, 1)
	handlerStarted := make(chan struct{}) // closed when handler begins processing

	srv := New(&config.ServerConfig{Port: 0})
	srv.echo.GET("/slow", func(c *echo.Context) error {
		close(handlerStarted) // step 5 — deterministic, no time.Sleep race
		time.Sleep(300 * time.Millisecond)
		return c.String(http.StatusOK, "done")
	})

	sc := echo.StartConfig{
		Address:          ":0",
		HideBanner:       true,
		HidePort:         true,
		ListenerAddrFunc: func(addr net.Addr) { addrCh <- addr }, // step 2
	}
	done := make(chan error, 1)
	go func() { done <- srv.Serve(ctx, sc) }() // step 1

	addr := <-addrCh // step 2 received

	// Step 3 — build request with its own timeout so the test never hangs
	// if the server fails to respond at all.
	reqCtx, reqCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer reqCancel()
	req, err := http.NewRequestWithContext(reqCtx, http.MethodGet,
		"http://"+addr.String()+"/slow", nil)
	require.NoError(t, err)

	type result struct {
		resp *http.Response
		err  error
	}
	resultCh := make(chan result, 1)
	go func() { // step 4
		resp, err := http.DefaultClient.Do(req)
		resultCh <- result{resp, err}
	}()

	// Step 6 — wait until handler is mid-execution before cancelling.
	// select + timeout prevents the test hanging if the handler is never reached
	// (e.g. routing misconfiguration, middleware short-circuit, server crash).
	select {
	case <-handlerStarted: // step 5 received
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for /slow handler to start")
	}
	cancel() // step 6 — graceful shutdown while /slow is still running

	// Step 10 — verify results
	res := <-resultCh
	// Precondition: if the request errored, the assertion below is meaningless.
	require.NoError(t, res.err, "in-flight request must not be dropped during graceful shutdown")
	require.NotNil(t, res.resp)
	t.Cleanup(func() { res.resp.Body.Close() }) // avoid fd leak

	assert.Equal(t, http.StatusOK, res.resp.StatusCode)

	// Server must exit cleanly after all in-flight requests complete.
	require.NoError(t, <-done)
}

func TestServe_GracefulShutdown_TimeoutCallsOnShutdownError(t *testing.T) {
	t.Parallel()

	// Test flow:
	//  1. Server starts with GracefulTimeout = 100ms.
	//  2. /very-slow handler takes 2s — longer than the timeout.
	//  3. cancel() triggers shutdown; server.Shutdown(100ms-ctx) fires.
	//  4. After 100ms, context expires → Shutdown returns DeadlineExceeded.
	//  5. OnShutdownError is called; gracefulShutdown goroutine finishes.
	//  6. wg.Wait() unblocks → Serve() returns nil promptly.
	//  NOTE: Echo v5 does NOT call server.Close() — handler goroutine
	//        continues running for up to 2s after Serve() returns.
	//        reqCancel() closes the client connection to avoid fd/goroutine leak.

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	addrCh := make(chan net.Addr, 1)
	handlerStarted := make(chan struct{})

	srv := New(&config.ServerConfig{Port: 0})
	srv.echo.GET("/very-slow", func(c *echo.Context) error {
		close(handlerStarted)
		time.Sleep(2 * time.Second)
		return c.String(http.StatusOK, "done")
	})

	var shutdownErr atomic.Pointer[error]
	sc := echo.StartConfig{
		Address:          ":0",
		HideBanner:       true,
		HidePort:         true,
		GracefulTimeout:  100 * time.Millisecond,
		ListenerAddrFunc: func(addr net.Addr) { addrCh <- addr },
		OnShutdownError:  func(err error) { shutdownErr.Store(&err) },
	}
	done := make(chan error, 1)
	go func() { done <- srv.Serve(ctx, sc) }()

	addr := <-addrCh

	// Dedicated client so CloseIdleConnections() in t.Cleanup only affects this test.
	// http.DefaultClient is shared across tests — mutating it would cause flakiness.
	testClient := &http.Client{Transport: &http.Transport{}}
	t.Cleanup(func() {
		testClient.CloseIdleConnections() // release any fd held by this client
	})

	// reqCancel is deferred: fires when this test function returns (~100ms after
	// cancel()). It closes the client-side connection, letting the
	// testClient.Do goroutine exit without leaking.
	//
	// Expected side-effect: the server-side handler goroutine is still sleeping
	// (2s total). When it wakes ~1.9s later and calls c.String(), the connection
	// is already gone → Echo logs "write: broken pipe". This is intentional and
	// harmless — the test has already passed by then. It is a direct consequence
	// of Echo v5 not calling server.Close() after GracefulTimeout (active
	// connections are not force-terminated). If you see "broken pipe" in the
	// test output after a PASS, this comment explains why.
	reqCtx, reqCancel := context.WithCancel(context.Background())
	defer reqCancel()
	req, err := http.NewRequestWithContext(reqCtx, http.MethodGet,
		"http://"+addr.String()+"/very-slow", nil)
	require.NoError(t, err)
	go testClient.Do(req) //nolint:errcheck

	select {
	case <-handlerStarted:
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for /very-slow handler to start")
	}

	shutdownAt := time.Now()
	cancel()

	// Verify Serve() returns in a window that proves the timeout fired:
	//   lower bound: at least GracefulTimeout has elapsed (shutdown actually waited)
	//   upper bound: well under handler duration (didn't wait for the 2s handler)
	const (
		gracefulTimeout = 100 * time.Millisecond
		handlerDuration = 2 * time.Second
		upperBound      = handlerDuration / 2 // 1s — comfortably below handler
	)
	select {
	case serveErr := <-done:
		require.NoError(t, serveErr)
		elapsed := time.Since(shutdownAt)
		t.Logf("Serve() returned after %v (GracefulTimeout=%v, handler=%v)",
			elapsed, gracefulTimeout, handlerDuration)
		assert.GreaterOrEqual(t, elapsed, gracefulTimeout,
			"Serve() returned too fast — GracefulTimeout may not have been respected")
		assert.Less(t, elapsed, upperBound,
			"Serve() waited too long — may be blocked on handler instead of timing out")
	case <-time.After(upperBound + time.Second):
		t.Fatal("Serve() did not return after GracefulTimeout expired")
	}

	// wg.Wait() in start() guarantees gracefulShutdown goroutine is done before
	// done fires → OnShutdownError has been called — no race.
	storedErr := shutdownErr.Load()
	require.NotNil(t, storedErr, "OnShutdownError must be called when GracefulTimeout expires")
	assert.ErrorIs(t, *storedErr, context.DeadlineExceeded)
}

func TestServe_GracefulShutdown_ClosesIdleConnections(t *testing.T) {
	t.Parallel()

	// Test flow:
	//  1. Server starts.
	//  2. Client sends a quick request → gets response → connection becomes idle (keep-alive).
	//  3. cancel() triggers graceful shutdown.
	//  4. server.Shutdown() must close the idle connection immediately.
	//  5. Serve() returns promptly — not waiting for keep-alive expiry (minutes).

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	addrCh := make(chan net.Addr, 1)
	srv := New(&config.ServerConfig{Port: 0})
	srv.echo.GET("/quick", func(c *echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	sc := echo.StartConfig{
		Address:          ":0",
		HideBanner:       true,
		HidePort:         true,
		ListenerAddrFunc: func(addr net.Addr) { addrCh <- addr },
	}
	done := make(chan error, 1)
	go func() { done <- srv.Serve(ctx, sc) }()

	addr := <-addrCh

	// Dedicated transport — two deliberate settings:
	//   DisableKeepAlives=false : ensure connections ARE reused (idle pool active)
	//   IdleConnTimeout=60s     : far longer than the test duration, so any
	//                             connection closure comes from server.Shutdown(),
	//                             not from a client-side idle timeout firing first.
	transport := &http.Transport{
		IdleConnTimeout: 60 * time.Second,
	}
	client := &http.Client{Transport: transport}

	reqCtx, reqCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer reqCancel()
	req, err := http.NewRequestWithContext(reqCtx, http.MethodGet,
		"http://"+addr.String()+"/quick", nil)
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)
	// Drain body fully BEFORE closing — HTTP/1.1 transport only returns a
	// connection to the idle pool when the body has been consumed entirely.
	// Closing without reading marks the connection "dirty"; transport closes
	// it instead of recycling, so the test would not actually exercise the
	// idle-connection path.
	_, err = io.Copy(io.Discard, resp.Body)
	require.NoError(t, err)
	resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)
	// Connection is now idle in the keep-alive pool.

	start := time.Now()
	cancel()

	// Idle connections are closed immediately by server.Shutdown() —
	// Serve() must return in well under 1s, not waiting for keep-alive timeout.
	select {
	case serveErr := <-done:
		require.NoError(t, serveErr)
		assert.Less(t, time.Since(start), time.Second,
			"idle connection should be closed immediately, not after keep-alive timeout")
	case <-time.After(5 * time.Second):
		t.Fatal("Serve() timed out — idle connection was not closed promptly")
	}
}

// ---------------------------------------------------------------------------
// Example as executable documentation
// ---------------------------------------------------------------------------

func ExampleServer_SetupRoutes() {
	srv := New(&config.ServerConfig{Port: 8080})
	srv.SetupMiddlewares()
	srv.SetupRoutes()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	srv.echo.ServeHTTP(rec, req)

	fmt.Println(rec.Code, strings.TrimSpace(rec.Body.String()))
	// Output: 200 OK
}
