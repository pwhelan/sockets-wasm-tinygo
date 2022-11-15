package socket

import (
	"fmt"
	"unsafe"
)

type AddressFamily int

const (
	AF_INET AddressFamily = 0
)

type socket_type int

const (
	SOCK_STREAM socket_type = 1
)

type errno int

type SocketAddressInet struct {
	Port    uint16
	Address string
}

type in_addr struct {
	s_addr uint32
}

type __wasi_addr_ip4_t struct {
	n0 uint8
	n1 uint8
	n2 uint8
	n3 uint8
}

// Implements __wasi_addr_t
type __wasi_addr_t struct {
	kind    uint32
	addr    __wasi_addr_ip4_t
	port    uint16 /* host byte order */
	__pad__ [4096]uint16
}

// Implements __wasi_iovec_t.
type __wasi_iovec_t struct {
	buf unsafe.Pointer
	len int
}

//go:wasm-module wasi_snapshot_preview1
//export sock_open
func sock_open(poolfd int, family AddressFamily, stype socket_type, fd *int) errno

//go:wasm-module wasi_snapshot_preview1
//export sock_connect
func sock_connect(fd int, addr *__wasi_addr_t) errno

//go:wasm-module wasi_snapshot_preview1
//export sock_close
func sock_close(fd int) errno

//go:wasm-module wasi_snapshot_preview1
//export sock_send
func sock_send(fd int, data *__wasi_iovec_t, data_len uint32, si_flags uint32, solen *int) errno

//go:wasm-module wasi_snapshot_preview1
//export sock_recv
func sock_recv(fd int, data *__wasi_iovec_t, data_len int, ri_flags int, rolen *int, roflags *int) errno

func Open(family AddressFamily, stype socket_type) (int, error) {
	var fd int

	errno := sock_open(0, family, stype, &fd)
	if errno != 0 {
		return -1, fmt.Errorf("ERRNO: %d", errno)
	}
	return fd, nil
}

func Connect(sockfd int, addr SocketAddressInet) error {
	sin := __wasi_addr_t{
		kind: 0,
		addr: __wasi_addr_ip4_t{127, 0, 0, 1},
		port: addr.Port,
	}
	errno := sock_connect(sockfd, &sin)
	if errno != 0 {
		return fmt.Errorf("ERRNO: %d", errno)
	}
	return nil
}

func Send(fd int, buf []byte) (int, error) {
	var sent int
	iovec := __wasi_iovec_t{
		buf: unsafe.Pointer(&buf[0]),
		len: len(buf),
	}

	errno := sock_send(fd, &iovec, uint32(iovec.len), 0, &sent)
	if errno != 0 {
		return -1, fmt.Errorf("ERRNO: %d", errno)
	}

	return sent, nil
}

func Recv(fd int, buf []byte) (int, error) {
	var recv int
	var flags int
	iovec := __wasi_iovec_t{
		buf: unsafe.Pointer(&buf[0]),
		len: len(buf),
	}

	errno := sock_recv(fd, &iovec, 1, 0, &recv, &flags)
	if errno != 0 {
		return 0, fmt.Errorf("ERRNO: %d", errno)
	}

	return recv, nil
}

func Close(fd int) error {
	errno := sock_close(fd)
	if errno != 0 {
		return fmt.Errorf("ERRNO: %d", errno)
	}
	return nil
}
