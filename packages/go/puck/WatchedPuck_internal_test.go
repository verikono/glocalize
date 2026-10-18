package puck

import (
	"errors"
	"strings"
	"testing"
	"time"

	NTFY "github.com/fsnotify/fsnotify"
)

func mustPanic(t *testing.T, fn func()) {
	t.Helper()
	defer func() {
		if recover() == nil {
			t.Fatalf("expected panic")
		}
	}()
	fn()
}

func TestPuck_NonFileProtocol(t *testing.T) {
	p := &Puck{url: "http://example.com/schema.prisma", protocol: ProtocolHttp}

	if got := p.Path(); got != "" {
		t.Fatalf("Path = %q, want empty", got)
	}
	if got := p.Extension(); got != "" {
		t.Fatalf("Extension = %q, want empty", got)
	}

	mustPanic(t, func() { _ = p.Exists() })
	mustPanic(t, func() { _ = p.Size() })
	mustPanic(t, func() { _ = p.resolveResource() })

	if err := p.assertValidPath(); err != nil {
		t.Fatalf("assertValidPath: %v", err)
	}
	if err := p.assertValidUrl(); err != nil {
		t.Fatalf("assertValidUrl: %v", err)
	}
	if err := p.assertValidFormat(); err != nil {
		t.Fatalf("assertValidFormat: %v", err)
	}
	if err := p.assertValidPuck(); err != nil {
		t.Fatalf("assertValidPuck: %v", err)
	}

	if err := p.OnChange(WithWatcher(func(FSChangeEvent) error { return nil })); err == nil || !strings.Contains(err.Error(), "cannot watch") {
		t.Fatalf("OnChange: %v", err)
	}
	if err := p.StopWatching(); err == nil || !strings.Contains(err.Error(), "cannot stop watching") {
		t.Fatalf("StopWatching: %v", err)
	}
}

func TestPuck_InvalidProtocol(t *testing.T) {
	p := &Puck{url: "gopher://example.com/x.prisma", protocol: "gopher"}

	mustPanic(t, func() { _ = p.Protocol() })
	if err := p.assertValidProtocol(); err == nil || !strings.Contains(err.Error(), "unsupported protocol") {
		t.Fatalf("assertValidProtocol: %v", err)
	}
	if err := p.assertValidPuck(); err == nil || !strings.Contains(err.Error(), "unsupported protocol") {
		t.Fatalf("assertValidPuck: %v", err)
	}
}

func TestPuck_NonAbsolutePath(t *testing.T) {
	p := &Puck{url: "file://relative.prisma", protocol: ProtocolFile}

	if err := p.assertValidPath(); err == nil || !strings.Contains(err.Error(), "must be an absolute") {
		t.Fatalf("assertValidPath: %v", err)
	}
	if err := p.assertValidPuck(); err == nil || !strings.Contains(err.Error(), "must be an absolute") {
		t.Fatalf("assertValidPuck: %v", err)
	}
	if _, err := p.Read(); err == nil || !strings.Contains(err.Error(), "must be an absolute") {
		t.Fatalf("Read: %v", err)
	}
	if err := p.Sync(); err == nil || !strings.Contains(err.Error(), "must be an absolute") {
		t.Fatalf("Sync: %v", err)
	}
}

func TestPuck_UrlMismatch(t *testing.T) {
	p := &Puck{url: "file:///one.prisma", protocol: ProtocolFile, glob: "/two.prisma"}

	if err := p.assertValidUrl(); err == nil {
		t.Fatal("assertValidUrl should error")
	}
	if err := p.assertValidPuck(); err == nil {
		t.Fatal("assertValidPuck should error")
	}
}

func TestPuck_FormatEmptyUrlPanics(t *testing.T) {
	p := &Puck{}
	if got := p.protocolFromUrl(); got != "" {
		t.Fatalf("protocolFromUrl = %q, want empty", got)
	}
	mustPanic(t, func() { _ = p.Format() })
}

func TestWatchTarget_Run(t *testing.T) {
	t.Run("delivers events", func(t *testing.T) {
		events := make(chan NTFY.Event)
		errs := make(chan error)
		target := &watchTarget{watcher: &NTFY.Watcher{Events: events, Errors: errs}}
		got := make(chan FSChangeEvent, 1)
		target.fn = func(event FSChangeEvent) error {
			got <- event
			return nil
		}

		done := make(chan struct{})
		go func() { target.run(); close(done) }()

		events <- NTFY.Event{Name: "/tmp/x.prisma", Op: NTFY.Write}
		select {
		case event := <-got:
			if event.Path != "/tmp/x.prisma" || event.Op != NTFY.Write {
				t.Fatalf("unexpected event: %+v", event)
			}
		case <-time.After(time.Second):
			t.Fatal("timed out waiting for event")
		}

		close(events)
		<-done
	})

	t.Run("drains errors", func(t *testing.T) {
		events := make(chan NTFY.Event)
		errs := make(chan error)
		target := &watchTarget{
			watcher: &NTFY.Watcher{Events: events, Errors: errs},
			fn:      func(FSChangeEvent) error { return nil },
		}

		done := make(chan struct{})
		go func() { target.run(); close(done) }()

		errs <- errors.New("boom")
		close(events)
		<-done
	})

	t.Run("stops on closed events channel", func(t *testing.T) {
		events := make(chan NTFY.Event)
		errs := make(chan error)
		target := &watchTarget{
			watcher: &NTFY.Watcher{Events: events, Errors: errs},
			fn:      func(FSChangeEvent) error { return nil },
		}

		done := make(chan struct{})
		go func() { target.run(); close(done) }()

		close(events)
		select {
		case <-done:
		case <-time.After(time.Second):
			t.Fatal("run did not stop after events channel closed")
		}
	})

	t.Run("stops on closed errors channel", func(t *testing.T) {
		events := make(chan NTFY.Event)
		errs := make(chan error)
		target := &watchTarget{
			watcher: &NTFY.Watcher{Events: events, Errors: errs},
			fn:      func(FSChangeEvent) error { return nil },
		}

		done := make(chan struct{})
		go func() { target.run(); close(done) }()

		close(errs)
		select {
		case <-done:
		case <-time.After(time.Second):
			t.Fatal("run did not stop after errors channel closed")
		}
	})
}
