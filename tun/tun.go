/* SPDX-License-Identifier: MIT
 *
 * Copyright (C) 2017-2023 WireGuard LLC. All Rights Reserved.
 */

package tun

import (
	"os"
)

type Event int

const (
	EventUp = 1 << iota
	EventDown
	EventMTUUpdate
)

const DefaultMTU = 1500

type Device interface {
	// File returns the file descriptor of the device.
	File() *os.File

	// ReadPackets one or more packets from the Device (without any additional headers).
	// On a successful read it returns the number of packets read, and sets
	// packet lengths within the sizes slice. len(sizes) must be >= len(bufs).
	// A nonzero offset can be used to instruct the Device on where to begin
	// reading into each element of the bufs slice.
	ReadPackets(bufs [][]byte, sizes []int, offset int) (n int, err error)

	// WritePackets one or more packets to the device (without any additional headers).
	// On a successful write it returns the number of packets written. A nonzero
	// offset can be used to instruct the Device on where to begin writing from
	// each packet contained within the bufs slice.
	WritePackets(bufs [][]byte, offset int) (int, error)

	// MTU returns the MTU of the Device.
	MTU() (int, error)

	// Name returns the current name of the Device.
	Name() (string, error)

	// Events returns a channel of type Event, which is fed Device events.
	Events() <-chan Event

	// Close stops the Device and closes the Event channel.
	Close() error

	// BatchSize returns the preferred/max number of packets that can be read or
	// written in a single read/write call. BatchSize must not change over the
	// lifetime of a Device.
	BatchSize() int

	Read(p []byte) (n int, err error)
	Write(p []byte) (n int, err error)
}

func (tun *NativeTun) Read(p []byte) (n int, err error) {
	var (
		bufs  = make([][]byte, 1)
		sizes = make([]int, 1)
	)

	bufs[0] = make([]byte, len(p))
	n, err = tun.ReadPackets(bufs, sizes, 0)
	if err != nil {
		return 0, err
	}
	if sizes[0] < 1 {
		return 0, nil
	}

	copy(p, bufs[0][:sizes[0]])

	return sizes[0], nil
}

func (tun *NativeTun) Write(p []byte) (n int, err error) {
	return tun.WritePackets([][]byte{p}, 0)
}
