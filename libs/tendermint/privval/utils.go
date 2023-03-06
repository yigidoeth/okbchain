package privval

import (
	"fmt"
	"net"

	"github.com/pkg/errors"

	"github.com/okx/okbchain/libs/tendermint/crypto/ed25519"
	"github.com/okx/okbchain/libs/tendermint/libs/log"
	tmnet "github.com/okx/okbchain/libs/tendermint/libs/net"
)

// IsConnTimeout returns a boolean indicating whether the error is known to
// report that a connection timeout occurred. This detects both fundamental
// network timeouts, as well as ErrConnTimeout errors.
func IsConnTimeout(err error) bool {
	switch errors.Cause(err).(type) {
	case EndpointTimeoutError:
		return true
	case timeoutError:
		return true
	default:
		return false
	}
}

// NewSignerListener creates a new SignerListenerEndpoint using the corresponding listen address
func NewSignerListener(listenAddr string, logger log.Logger) (*SignerListenerEndpoint, error) {
	var listener net.Listener

	protocol, address := tmnet.ProtocolAndAddress(listenAddr)
	ln, err := net.Listen(protocol, address)
	if err != nil {
		return nil, err
	}
	switch protocol {
	case "unix":
		listener = NewUnixListener(ln)
	case "tcp":
		// TODO: persist this key so external signer can actually authenticate us
		listener = NewTCPListener(ln, ed25519.GenPrivKey())
	default:
		return nil, fmt.Errorf(
			"wrong listen address: expected either 'tcp' or 'unix' protocols, got %s",
			protocol,
		)
	}

	pve := NewSignerListenerEndpoint(logger.With("module", "privval"), listener)

	return pve, nil
}

// GetFreeLocalhostAddrPort returns a free localhost:port address
func GetFreeLocalhostAddrPort() string {
	port, err := tmnet.GetFreePort()
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("127.0.0.1:%d", port)
}
