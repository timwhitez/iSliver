package main

import (
	"bytes"
	"os/exec"
	"syscall"

	"golang.org/x/sys/windows"

	"log"
)

func main() {

	startProcess("C:\\windows\\system32\\notepad.exe", "C:\\windows\\system32", nil, 9700, nil, nil, false)

}

func startProcess(proc string, currentDir string, args []string, ppid uint32, stdout *bytes.Buffer, stderr *bytes.Buffer, suspended bool) (*exec.Cmd, error) {
	var CurrentToken windows.Token
	var cmd *exec.Cmd
	if len(args) > 0 {
		cmd = exec.Command(proc, args...)
	} else {
		cmd = exec.Command(proc)
	}
	cmd.SysProcAttr = &windows.SysProcAttr{
		Token:      syscall.Token(CurrentToken),
		HideWindow: true,
	}
	err := SpoofParent(ppid, cmd, currentDir)
	if err != nil {
		log.Printf("could not spoof parent PID: %v\n", err)
	}
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	if suspended {
		cmd.SysProcAttr.CreationFlags = windows.CREATE_SUSPENDED
	}
	err = cmd.Start()
	if err != nil {
		log.Println("Could not start process:", proc)
		return nil, err
	}
	return cmd, nil
}

func SpoofParent(ppid uint32, cmd *exec.Cmd, currentDir string) error {
	parentHandle, err := windows.OpenProcess(windows.PROCESS_CREATE_PROCESS|windows.PROCESS_DUP_HANDLE|windows.PROCESS_QUERY_INFORMATION, false, ppid)
	if err != nil {
		log.Printf("OpenProcess failed: %v\n", err)
		return err
	}
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &windows.SysProcAttr{}
	}
	cmd.SysProcAttr.ParentProcess = syscall.Handle(parentHandle)
	cmd.Dir = currentDir
	return nil
}
