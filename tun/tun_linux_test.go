package tun

import (
	"bytes"
	"fmt"
	"os/exec"
	"testing"
)

func TestCreate(t *testing.T) {
	//d, err := CreateTUN("", DefaultMTU, unix.IFF_VNET_HDR)
	d, err := CreateTUN("", DefaultMTU, 0)
	if err != nil {
		t.Fatalf("failed create device: %v", err)
	}

	if ev := <-d.Events(); ev != EventUp {
		t.Fatal("device in not up")
	}

	deviceName, err := d.Name()
	if err != nil {
		t.Fatalf("failed to get device name: %v", err)
	}

	if err := execute(fmt.Sprintf("ip addr add dev %s local 192.168.100.1 remote 192.168.100.2", deviceName)); err != nil {
		t.Fatalf("failed to set interface address: %v", err)
	}
	if err := execute(fmt.Sprintf("ip link set dev %s mtu %d up", deviceName, DefaultMTU)); err != nil {
		t.Fatalf("failed to set interface MTU: %v", err)
	}

	buf := make([]byte, DefaultMTU)
	nr, err := d.Read(buf)
	if err != nil {
		t.Fatalf("failed reading from device: %v", err)
	}

	nw, err := d.Write(buf[:nr])
	if err != nil {
		t.Fatalf("failed to write to device: %v", err)
	}

	if nw != nr {
		t.Fatalf("write failed nw=%d nr=%d", nw, nr)
	}
}

func execute(command string) error {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command("bash", append([]string{"-c"}, command)...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	return cmd.Run()
}
