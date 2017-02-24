// +build windows,!gendocs

package sshfs

import (
	"fmt"
	"syscall"
)

func CreateSftpDaemon() {
	RunSftpDaemon()
}

func createPipes() (syscall.Handle, syscall.Handle){
	var reader syscall.Handle
	var writer syscall.Handle

	syscall.CreatePipe(&reader, &writer, nil, 0)
	return reader, writer
}

func createWindowsProcess() {
	_ , _ = createPipes()

	var startupInfo syscall.StartupInfo
	var processInfo syscall.ProcessInformation

	argv := syscall.StringToUTF16Ptr("c:\\windows\\system32\\calc.exe")
	//argv := syscall.StringToUTF16Ptr("echo hello world")

	err := syscall.CreateProcess(
		nil,
		argv,
		nil,
		nil,
		false,
		0,
		nil,
		nil,
		&startupInfo,
		&processInfo)

	fmt.Printf("Return: %d\n", err)
}
