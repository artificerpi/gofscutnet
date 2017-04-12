package main

// +build windows,!cgo
import (
	"errors"
	"os"
	"runtime"
	"syscall"
	"unsafe"
)

// get device adapter in windows
func getDeviceAdapterName(Index int) (string, error) {
	if runtime.GOOS != "windows" {
		err := errors.New("You can not use the method while not on windows")
		return "", err
	}
	b := make([]byte, 1000)
	l := uint32(len(b))
	aList := (*syscall.IpAdapterInfo)(unsafe.Pointer(&b[0]))
	// TODO(mikio): GetAdaptersInfo returns IP_ADAPTER_INFO that
	// contains IPv4 address list only. We should use another API
	// for fetching IPv6 stuff from the kernel.
	err := syscall.GetAdaptersInfo(aList, &l)
	if err == syscall.ERROR_BUFFER_OVERFLOW {
		b = make([]byte, l)
		aList = (*syscall.IpAdapterInfo)(unsafe.Pointer(&b[0]))
		err = syscall.GetAdaptersInfo(aList, &l)
	}
	if err != nil {
		return "", os.NewSyscallError("GetAdaptersInfo", err)
	}

	// get right adapter of the device
	for ai := aList; ai != nil; ai = ai.Next {
		if int(ai.Index) == Index {
			return string(ai.AdapterName[:]), nil
		}
	}

	err = errors.New("invalid index as parameter")
	return "", err
}
