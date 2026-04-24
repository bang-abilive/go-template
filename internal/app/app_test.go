package app

import (
	"context"
	"net"
	"testing"
	"time"

	"ndinhbang/go-skeleton/pkg/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewApp(t *testing.T) {
	t.Parallel()
	cfg := &config.Config{Server: config.ServerConfig{Port: 0}}
	a := NewApp(cfg)
	require.NotNil(t, a)
	assert.Same(t, cfg, a.cfg)
}

func TestApp_Run_contextCancel(t *testing.T) {
	cfg := &config.Config{Server: config.ServerConfig{Port: 0}}
	a := NewApp(cfg)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	t.Cleanup(cancel)

	err := a.Run(ctx)
	assert.NoError(t, err)
}

func TestApp_Run_bindError(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	t.Cleanup(func() { _ = ln.Close() })
	port := uint16(ln.Addr().(*net.TCPAddr).Port)

	a := NewApp(&config.Config{Server: config.ServerConfig{Port: port}})
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	t.Cleanup(cancel)

	runErr := a.Run(ctx)
	require.Error(t, runErr, "port is already held by a listener, Start should fail")
}
