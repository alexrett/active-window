//+build windows

package activeWindow

import (
	// "fmt"
	"syscall"
	"path/filepath"
	"unsafe"
	"github.com/Andoryuuta/kiwi/w32"
    "golang.org/x/sys/windows"
)

var (
	user32 					= windows.NewLazyDLL("user32.dll")
	psapi         		= windows.NewLazyDLL("psapi.dll")
	procGetWindowText   	= user32.NewProc("GetWindowTextW")
	procGetWindowTextLength = user32.NewProc("GetWindowTextLengthW")
	procGetWindowThreadProcessId = user32.NewProc("GetWindowThreadProcessId")
	procGetProcessImageFileNameA = psapi.NewProc("GetProcessImageFileNameA")
)

type (
	HANDLE uintptr
	HWND HANDLE
)

func GetWindowTextLength(hwnd HWND) int {
	ret, _, _ := procGetWindowTextLength.Call(
		uintptr(hwnd))

	return int(ret)
}

func getFileNameByPID(pid uint32) string {
	var fileName string = `<Unknown File>`

	//Open process
	hnd, ok := w32.OpenProcess(w32.PROCESS_QUERY_INFORMATION, false, pid)
	if !ok {
		return fileName
	}
	defer w32.CloseHandle(hnd)

	//Get file path
	path, ok := w32.GetProcessImageFileName(hnd)
	if !ok {
		return fileName
	}

	//Split file path to get file name
	_, fileName = filepath.Split(path)
	return fileName
}

func GetProcessImageFileName(hwnd HWND) string {
	textLen := GetWindowTextLength(hwnd) + 1

	buf := make([]uint16, textLen)
	procGetProcessImageFileNameA.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(textLen))

	return syscall.UTF16ToString(buf)
}

func GetWindowText(hwnd HWND) string {
	textLen := GetWindowTextLength(hwnd) + 1

	buf := make([]uint16, textLen)
	procGetWindowText.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(textLen))

	return syscall.UTF16ToString(buf)
}

func getWindow(funcName string) uintptr {
    proc := user32.NewProc(funcName)
    hwnd, _, _ := proc.Call()
    return hwnd
}

func GetWindowThreadProcessId(hwnd HWND) uint32 {
	var procId uint32
	procGetWindowThreadProcessId.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(&procId)))

	return procId
}


func (a *ActiveWindow) getActiveWindowTitle() (string, string) {

	if hwnd := getWindow("GetForegroundWindow") ; hwnd != 0 {
		text := GetWindowText(HWND(hwnd))
		pid := GetWindowThreadProcessId(HWND(hwnd))
		process := getFileNameByPID(uint32(pid))
		// splitted := strings.Split(" - ", text)
		// fmt.Println("window :", text, "# hwnd:", hwnd, "#pid:", pid, "# proc:", process)
		return process,text
	}

	return "", ""
}
