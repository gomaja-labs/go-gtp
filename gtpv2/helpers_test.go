package gtpv2_test

import (
	"fmt"
	"math"
	"net"
	"strconv"
	"testing"

	"github.com/gomaja-labs/go-gtp/gtpv2"
)

var dummyAddr net.Addr = &net.UDPAddr{IP: net.IP{0x00, 0x00, 0x00, 0x00}, Port: 2123}

const testIMSIPrefix = "00101123456789"

// newSessionFixture returns a Conn holding four active sessions. Session i
// (1-based) has S11 TEID i and an IMSI ending in the digit i. Every test builds
// its own fixture so that tests are independent and repeatable.
func newSessionFixture(t *testing.T) (*gtpv2.Conn, []*gtpv2.Session) {
	t.Helper()
	conn := gtpv2.NewConn(dummyAddr, gtpv2.IFTypeS11MMEGTPC, 0)
	sessions := make([]*gtpv2.Session, 4)
	for i := range sessions {
		sess := gtpv2.NewSession(dummyAddr, &gtpv2.Subscriber{IMSI: testIMSIPrefix + strconv.Itoa(i+1)})
		if err := sess.Activate(); err != nil {
			t.Fatal(err)
		}
		sess.AddTEID(gtpv2.IFTypeS11MMEGTPC, uint32(i+1))
		conn.RegisterSession(uint32(i+1), sess)
		sessions[i] = sess
	}
	return conn, sessions
}

// checkRemaining asserts that TEID 1 no longer resolves and TEIDs 2 to 4 still
// resolve to their own sessions.
func checkRemaining(t *testing.T, conn *gtpv2.Conn) {
	t.Helper()
	if got := conn.SessionCount(); got != 3 {
		t.Fatalf("SessionCount after removal = %d, want 3", got)
	}
	if sess, err := conn.GetSessionByTEID(1); err == nil {
		t.Errorf("removed TEID 1 still resolves to %s", sess.IMSI)
	}
	if _, err := conn.GetSessionByIMSI(testIMSIPrefix + "1"); err == nil {
		t.Error("removed IMSI still resolves")
	}
	for i := 2; i <= 4; i++ {
		sess, err := conn.GetSessionByTEID(uint32(i))
		if err != nil {
			t.Fatalf("TEID %d: %v", i, err)
		}
		if want := testIMSIPrefix + strconv.Itoa(i); sess.IMSI != want {
			t.Errorf("TEID %d resolves to %s, want %s", i, sess.IMSI, want)
		}
	}
}

func TestSessionCount(t *testing.T) {
	conn, sessions := newSessionFixture(t)
	if got, want := conn.SessionCount(), len(sessions); got != want {
		t.Errorf("SessionCount = %d, want %d", got, want)
	}
}

func TestGetSessionByIMSI_GetTEID(t *testing.T) {
	conn, sessions := newSessionFixture(t)
	for i := range sessions {
		sess, err := conn.GetSessionByIMSI(testIMSIPrefix + strconv.Itoa(i+1))
		if err != nil {
			t.Fatal(err)
		}
		teid, err := sess.GetTEID(gtpv2.IFTypeS11MMEGTPC)
		if err != nil {
			t.Fatal(err)
		}
		if teid != uint32(i+1) {
			t.Errorf("%s: TEID = %d, want %d", sess.IMSI, teid, i+1)
		}
	}
}

func BenchmarkAddSession(b *testing.B) {
	for k := 0.; k < 6; k++ {
		existingSessions := int(math.Pow(10, k))
		benchConn := gtpv2.NewConn(dummyAddr, gtpv2.IFTypeS11MMEGTPC, 0)
		for i := 0; i < existingSessions; i++ {
			imsi := fmt.Sprintf("%015d", i)
			benchConn.RegisterSession(0, gtpv2.NewSession(dummyAddr, &gtpv2.Subscriber{IMSI: imsi}))
		}
		b.Run(fmt.Sprintf("%d", existingSessions), func(b *testing.B) {
			for i := 1; i <= b.N; i++ {
				benchConn.RegisterSession(0, gtpv2.NewSession(dummyAddr, &gtpv2.Subscriber{IMSI: "001011234567891"}))
			}
		})
	}
}

func TestGetSessionByTEID(t *testing.T) {
	conn, sessions := newSessionFixture(t)
	for i := range sessions {
		sess, err := conn.GetSessionByTEID(uint32(i + 1))
		if err != nil {
			t.Fatal(err)
		}
		if want := testIMSIPrefix + strconv.Itoa(i+1); sess.IMSI != want {
			t.Errorf("TEID %d resolves to %s, want %s", i+1, sess.IMSI, want)
		}
	}
}

func TestGetIMSIByTEID(t *testing.T) {
	conn, sessions := newSessionFixture(t)
	for i := range sessions {
		imsi, err := conn.GetIMSIByTEID(uint32(i + 1))
		if err != nil {
			t.Fatal(err)
		}
		if want := testIMSIPrefix + strconv.Itoa(i+1); imsi != want {
			t.Errorf("TEID %d resolves to IMSI %s, want %s", i+1, imsi, want)
		}
	}
}

func TestRemoveSession(t *testing.T) {
	conn, sessions := newSessionFixture(t)
	conn.RemoveSession(sessions[0])
	checkRemaining(t, conn)
}

func TestRemoveSessionByIMSI(t *testing.T) {
	conn, _ := newSessionFixture(t)
	conn.RemoveSessionByIMSI(testIMSIPrefix + "1")
	checkRemaining(t, conn)
}
