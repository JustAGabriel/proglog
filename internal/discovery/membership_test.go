package discovery

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/serf/serf"
	"github.com/justagabriel/proglog/internal"
	"github.com/stretchr/testify/require"
)

func TestMembership(t *testing.T) {
	// act
	m, handler := setupMember(t, nil)
	m, _ = setupMember(t, m)
	m, _ = setupMember(t, m)

	// assert
	require.Eventually(t,
		func() bool {
			return len(handler.joins) == 2 &&
				len(m[0].Members()) == 3 &&
				len(handler.leaves) == 0
		},
		5*time.Second,
		250*time.Microsecond)

	require.NoError(t, m[2].Leave())
	require.Eventually(t,
		func() bool {
			return len(handler.joins) == 2 &&
				len(m[0].Members()) == 3 &&
				m[0].Members()[2].Status == serf.StatusLeft &&
				len(handler.leaves) == 1
		},
		3*time.Second,
		250*time.Microsecond)

	require.Equal(t, fmt.Sprintf("%d", 2), <-handler.leaves)
}

func setupMember(t *testing.T, members []*Membership) ([]*Membership, *handler) {
	id := len(members)
	port := internal.FreePort(t)
	addr := fmt.Sprintf("%s:%d", "127.0.0.1", port)
	tags := map[string]string{
		"rpc_addr": addr,
	}

	config := Config{
		NodeName: fmt.Sprintf("%d", id),
		BindAddr: addr,
		Tags:     tags,
	}

	handler := &handler{}
	if len(members) == 0 {
		handler.joins = make(chan map[string]string, 3)
		handler.leaves = make(chan string, 3)
	} else {
		config.StartJoinAddrs = []string{
			members[0].BindAddr,
		}
	}

	m, err := New(handler, config)
	require.NoError(t, err)
	members = append(members, m)
	return members, handler
}

type handler struct {
	joins  chan map[string]string
	leaves chan string
}

func (h *handler) Join(id, addr string) error {
	if h.joins != nil {
		h.joins <- map[string]string{
			"id":   id,
			"addr": addr,
		}
	}
	return nil
}

func (h *handler) Leave(id string) error {
	if h.leaves != nil {
		h.leaves <- id
	}
	return nil
}
