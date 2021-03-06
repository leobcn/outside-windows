package winbase

import (
	. "github.com/tHinqa/outside"
	. "github.com/tHinqa/outside-windows/types"
	. "github.com/tHinqa/outside/types"
	"syscall"
	"testing"
)

import . "unsafe"

func init() {
	AddApis(WinBaseApis)
	AddApis(WinBaseANSIApis)
	//AddApis(WinBaseUnicodeApis)
	SI.Len = DWORD(Sizeof(SI))
}

var SI STARTUPINFO
var ST SYSTEMTIME

func Test(t *testing.T) {
	Buffer := new(MEMORYSTATUSEX)
	SetStructSize(Buffer)
	ret, _ := GlobalMemoryStatusEx(Buffer)
	t.Logf("%v %+v\n", ret, *Buffer)

	var a, i DWORD
	var b BOOL
	gsta, _ := GetSystemTimeAdjustment(&a, &i, &b)
	t.Logf("%v %v %v %v\n", gsta, a, i, b)

	GetStartupInfo(&SI)
	GetStartupInfo(&SI)
	t.Logf("%+v\n%s\n%s\n", SI, *SI.Desktop, *SI.Title)

	GetLocalTime(&ST)
	t.Logf("%+v\n", ST)
}

func BenchmarkSyscall(b *testing.B) {
	var d = syscall.MustLoadDLL("kernel32.dll")
	defer d.Release()
	var G = d.MustFindProc("GetStartupInfoA")
	for i := 0; i < b.N; i++ {
		syscall.Syscall(G.Addr(), 1, (uintptr)(Pointer(&SI)), 0, 0)
	}
}

func BenchmarkReflectStartupInfo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetStartupInfo(&SI)
	}
}

func BenchmarkVanillaStartupInfo(b *testing.B) {
	var d = syscall.MustLoadDLL("kernel32.dll")
	var G = d.MustFindProc("GetStartupInfoA")
	defer d.Release()
	for i := 0; i < b.N; i++ {
		syscall.Syscall(G.Addr(), 1, (uintptr)(Pointer(&SI)), 0, 0)
		s := OVString(CStrToString((uintptr)(Pointer(SI.Desktop))))
		SI.Desktop = &s
		s = OVString(CStrToString((uintptr)(Pointer(SI.Title))))
		SI.Title = &s
	}
}

func BenchmarkReflectSystemTime(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetLocalTime(&ST)
	}
}

func BenchmarkVanillaSystemTime(b *testing.B) {
	var d = syscall.MustLoadDLL("kernel32.dll")
	var G = d.MustFindProc("GetLocalTime")
	defer d.Release()
	for i := 0; i < b.N; i++ {
		syscall.Syscall(G.Addr(), 1, (uintptr)(Pointer(&ST)), 0, 0)
	}
}
