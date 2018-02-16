package info

import (
	"os/exec"
	"strings"
	"bytes"
	"runtime"
	"os"
	"fmt"
)

type GoInfoObject struct {
	GoOS     string
	Kernel   string
	Core     string
	Platform string
	OS       string
	Hostname string
	CPUs     int
}

func GetInfo() *GoInfoObject {
	cmd := exec.Command("cmd", "ver")
	cmd.Stdin = strings.NewReader("some input")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
	osStr := strings.Replace(out.String(), "\n", "", -1)
	osStr = strings.Replace(osStr, "\r\n", "", -1)
	tmp1 := strings.Index(osStr, "[Version")
	tmp2 := strings.Index(osStr, "]")
	var ver string
	if tmp1 == -1 || tmp2 == -1 {
		ver = "unknown"
	} else {
		ver = osStr[tmp1+9 : tmp2]
	}
	gio := &GoInfoObject{Kernel: "windows", Core: ver, Platform: "unknown", OS: "windows", GoOS: runtime.GOOS, CPUs: runtime.NumCPU()}
	gio.Hostname, _ = os.Hostname()
	return gio
}

func (gi *GoInfoObject) VarDump() {
	fmt.Println("GoOS:", gi.GoOS)
	fmt.Println("Kernel:", gi.Kernel)
	fmt.Println("Core:", gi.Core)
	fmt.Println("Platform:", gi.Platform)
	fmt.Println("OS:", gi.OS)
	fmt.Println("Hostname:", gi.Hostname)
	fmt.Println("CPUs:", gi.CPUs)
}

func (gi *GoInfoObject) GetPcInfo() string {
	return fmt.Sprintf("GoOS:%v\nKernel:%v\nCore:%v\nPlatform:%v\nOS:%v\nHostname:%v\nCPUs:%v", gi.GoOS, gi.Kernel, gi.Core, gi.Platform, gi.OS, gi.Hostname, gi.CPUs)
}
