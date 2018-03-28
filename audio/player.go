package audio

import (
	"syscall"
	"unsafe"
	"fmt"
	"strconv"
)

var (
	winmm         = syscall.MustLoadDLL("winmm.dll")
	mciSendString = winmm.MustFindProc("mciSendStringW")
	mciGetErrorString = winmm.MustFindProc("mciGetErrorStringW")
)

type Player struct {
	fileName string
	isOpen bool
	mediaName string
	buf []uint16
}

func New() *Player {
	return &Player{mediaName: "media", buf: make([]uint16, 64)}
}

func command(lpstrCommand string, b []uint16) int {
	mPtr := *(*uintptr)(unsafe.Pointer(&b))
	i, _, _ := mciSendString.Call(uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(lpstrCommand))),
		mPtr,
		uintptr(len(b)), uintptr(0))
	return int(i)
}

func getErrorCommand(err int, b []uint16) int {
	mPtr := *(*uintptr)(unsafe.Pointer(&b))
	i, _, _ := mciGetErrorString.Call(uintptr((err)),
		mPtr,
		uintptr(len(b)))
	return int(i)
}

func (this *Player) Play() error {
	if (this.isOpen) {
		playCommand := "Play " + this.mediaName// + " notify";
		if this.IsPause(){
			playCommand = "Resume " + this.mediaName
		}
		_, err := this.cmd(playCommand)
		return err
	}else if this.fileName != "" {
		if result := this.openMediaFile(); result != nil {
			return result
		}
		return this.Play()
	}
	return nil
}

func (this *Player) Pause() error {
	if (this.isOpen) {
		playCommand := "Pause " + this.mediaName
		if _, err := this.cmd(playCommand); err != nil {
			return err
		}
	}
	return nil
}


func (this *Player) Stop() error {
	return this.Close()
}

func (this *Player) Close() error {
	if (this.isOpen) {
		playCommand := "Close " + this.mediaName//+ " wait"
		if _, err := this.cmd(playCommand); err != nil {
			return err
		}
		this.isOpen = false
	}
	return nil
}

func (this *Player) openMediaFile() error {
	if result := this.Close(); result != nil {
		return result
	}
	playCommand := "Open \"" + this.fileName + "\" type mpegvideo alias " + this.mediaName + " wait"
	//playCommand := "Open \"" + this.fileName + "\" alias " + this.mediaName// + " wait"


	if _, err := this.cmd(playCommand); err != nil {
		return err
	}
	this.isOpen = true
	return nil
}

func (this *Player) OpenFile(fileName string) error {
	this.fileName = fileName
	if err := this.openMediaFile(); err != nil {
		return err
	}
	return this.Play()
}

func (this *Player) IsPaying() bool {
	playCommand := "Status " + this.mediaName + " mode"
	status, err := this.cmd(playCommand)
	if err != nil {
		return false
	}
	return status == "playing"
}

func (this *Player) IsStop() bool {
	playCommand := "Status " + this.mediaName + " mode"
	status, err := this.cmd(playCommand)
	if err != nil {
		return false
	}
	return status == "stopped"
}

func (this *Player) IsPause() bool {
	playCommand := "Status " + this.mediaName + " mode"
	status, err := this.cmd(playCommand)
	if err != nil {
		return false
	}
	return status == "paused"
}

func (this *Player) SongTime() int {
	if this.IsPaying() {
		playCommand := "Status " + this.mediaName + " length wait"
		status, err := this.cmd(playCommand)
		if err != nil {
			return 0
		}
		r,_ := strconv.Atoi(status)
		return r
	}
	return 0
}

func (this *Player) CurrentPos() int {
	if this.IsPaying() {
		playCommand := "Status " + this.mediaName + " position wait"
		status, err := this.cmd(playCommand)
		if err != nil {
			return 0
		}
		r,_ := strconv.Atoi(status)
		return r
	}
	return 0
}

func (this *Player) cmd(lpstrCommand string) (str string, err error) {
	i := command(lpstrCommand, this.buf)
	str = syscall.UTF16ToString(this.buf)
	if i != 0 {
		str = ""
		getErrorCommand(i, this.buf)
		err = fmt.Errorf("command: %s; error: %s", lpstrCommand,  syscall.UTF16ToString(this.buf))
	}
	return str, err
}



