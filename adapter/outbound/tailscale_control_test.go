//go:build with_gvisor && !no_tailscale

package outbound

import (
	"context"
	"path/filepath"
	"strings"
	"testing"

	C "github.com/metacubex/mihomo/constant"
)

func TestTailscaleStatusAfterClose(t *testing.T) {
	outbound, err := NewTailscale(TailscaleOption{Name: "test"})
	if err != nil {
		t.Fatal(err)
	}
	if err := outbound.Close(); err != nil {
		t.Fatal(err)
	}
	_, err = outbound.StartAndStatus(context.Background())
	if err == nil || !strings.Contains(err.Error(), "closed") {
		t.Fatalf("expected a closed error, got %v", err)
	}
}

func TestNewNodeClosesTheOneOnItsStateDir(t *testing.T) {
	previous := C.Path.HomeDir()
	home := t.TempDir()
	C.SetHomeDir(home)
	t.Cleanup(func() { C.SetHomeDir(previous) })
	stateDir := filepath.Join(home, "state")
	first, err := NewTailscale(TailscaleOption{Name: "first", StateDir: stateDir})
	if err != nil {
		t.Fatal(err)
	}
	second, err := NewTailscale(TailscaleOption{Name: "second", StateDir: stateDir})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = second.Close() })

	if _, err := first.StartAndStatus(context.Background()); err == nil {
		t.Fatal("expected the previous node to be closed")
	}
}
