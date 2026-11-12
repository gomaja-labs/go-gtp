package gtpv1

import (
	"errors"
	"github.com/gomaja-labs/go-gtp/gtpv1/message"
	"golang.org/x/sys/unix"
	"testing"
	"time"

	"github.com/vishvananda/netlink"
)

func TestKernelLinkRollback(t *testing.T) {
	failure := errors.New("synthetic setup failure")
	for _, stage := range []string{"add", "up", "mtu", "success"} {
		t.Run(stage, func(t *testing.T) {
			exists := false
			deleted := 0
			ops := kernelLinkOps{
				add: func(netlink.Link) error {
					if stage == "add" {
						return failure
					}
					exists = true
					return nil
				},
				up: func(netlink.Link) error {
					if stage == "up" {
						return failure
					}
					return nil
				},
				mtu: func(netlink.Link, int) error {
					if stage == "mtu" {
						return failure
					}
					return nil
				},
				del: func(netlink.Link) error { exists = false; deleted++; return nil },
			}
			err := configureKernelLink(&netlink.GTP{}, ops)
			if stage == "success" {
				if err != nil || !exists || deleted != 0 {
					t.Fatal("successful link was removed")
				}
				return
			}
			if !errors.Is(err, failure) || exists {
				t.Fatalf("failed link leaked: exists=%v error=%v", exists, err)
			}
			if (stage == "add" && deleted != 0) || (stage != "add" && deleted != 1) {
				t.Fatalf("wrong rollback count: %d", deleted)
			}
		})
	}
}

// This exercises the descriptor acquisition used by EnableKernelGTP without
// requiring a kernel GTP device or netlink privileges.
func TestKernelSocketDescriptorRemainsNonblocking(t *testing.T) {
	sock := listenUDP(t, "127.0.0.1", 0)
	u, err := NewUPlaneConnFromUDPConn(sock)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = u.Close() })
	var flags int
	if err := withKernelSocket(u.packetConn(), func(fd uintptr) error {
		var err error
		flags, err = unix.FcntlInt(fd, unix.F_GETFL, 0)
		// Restore the flag even on a regression so cleanup cannot hang.
		if err == nil && flags&unix.O_NONBLOCK == 0 {
			return unix.SetNonblock(int(fd), true)
		}
		return err
	}); err != nil {
		t.Fatal(err)
	}
	if flags&unix.O_NONBLOCK == 0 {
		t.Fatal("kernel descriptor acquisition made the shared socket blocking")
	}
	u = serveUDP(t, sock, nil)
	peer := listenUDP(t, "127.0.0.1", 0)
	sendUDP(t, peer, u.LocalAddr(), message.NewEchoRequest(47))
	_, _ = readUDP(t, peer)
	done := make(chan error, 1)
	go func() { done <- u.Close() }()
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(time.Second):
		t.Fatal("Close blocked after kernel descriptor acquisition")
	}
}

func TestFileFDChangesSharedSocketFlags(t *testing.T) {
	sock := listenUDP(t, "127.0.0.1", 0)
	f, err := sock.File()
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = f.Close() }()
	fd := f.Fd()
	flags, err := unix.FcntlInt(fd, unix.F_GETFL, 0)
	if err != nil {
		t.Fatal(err)
	}
	if err := unix.SetNonblock(int(fd), true); err != nil {
		t.Fatal(err)
	}
	if flags&unix.O_NONBLOCK != 0 {
		t.Fatal("File/Fd no longer reproduces the shared blocking flag")
	}
}
