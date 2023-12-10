package tun

import (
	"bytes"
	"fmt"
	"os/exec"
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	d, err := CreateTUN(fmt.Sprintf("tun-%d", time.Now().Unix()), "My", DefaultMTU)
	if err != nil {
		t.Fatalf("failed create device: %v", err)
	}

	deviceName, err := d.Name()
	if err != nil {
		t.Fatalf("failed to get device name: %v", err)
	}

	if err := execute(fmt.Sprintf("Disable-NetAdapterBinding -Name \"%s\" -ComponentID ms_tcpip6", deviceName)); err != nil {
		t.Fatalf("failed to set interface address: %v", err)
	}
	if err := execute(fmt.Sprintf("netsh interface ip set address name=\"%s\" source=static addr=192.168.100.1 mask=255.255.255.0 gateway=none", deviceName)); err != nil {
		t.Fatalf("failed to set interface address: %v", err)
	}
	if err := execute(fmt.Sprintf("netsh interface ip set interface \"%s\" mtu=%d", deviceName, DefaultMTU)); err != nil {
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

	cmd := exec.Command("powershell.exe", append([]string{"-c"}, command)...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	return cmd.Run()
}
