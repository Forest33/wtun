package tun

import (
	"bytes"
	"fmt"
	"os/exec"
	"testing"
)

func TestCreate(t *testing.T) {
	d, err := CreateTUN("", DefaultMTU)
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

	if err := execute(fmt.Sprintf("ifconfig %s 192.168.100.1 192.168.100.2 mtu %d up", deviceName, DefaultMTU)); err != nil {
		t.Fatalf("failed to set interface address: %v", err)
	}

	buf := make([]byte, DefaultMTU)
	n, err := d.Read(buf)
	if err != nil {
		t.Fatalf("failed reading from device: %v", err)
	}

	n, err = d.Write(buf[:n])
	if err != nil {
		t.Fatalf("failed to write to device: %v", err)
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
