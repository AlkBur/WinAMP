package audio

import (
	"testing"
	"time"
	"fmt"
	"strconv"
)

var songs = struct {
	id int
	arr []string
}{
	id: -1,
	arr: []string{
	"1.mp3",
	"2.mp3",
	"3.mp3"},
}

func PlayNext(p *Player, t *testing.T)  {
	songs.id++
	if songs.id >= len(songs.arr){
		songs.id = 0;
	}
	err := p.OpenFile(songs.arr[songs.id])
	if err != nil {
		t.Fatal(err)
	}
	t.Log(songs.arr[songs.id])
}

func TestPlay(t *testing.T)  {
	p := New()
	defer p.Close()

	PlayNext(p, t)

	str,_:=p.cmd("status "+p.mediaName+" length wait")
	tm,_ := strconv.Atoi(str)
	td := float32(tm)/60000
	fmt.Println(td)


	stopPlay := false
	for !stopPlay {
		select {
		case <-time.After(time.Second):
			if !p.IsPaying() {
				t.Log("Not play")
				stopPlay = true
			}
			str,_=p.cmd("status "+p.mediaName+" position wait")


			fmt.Println(str)

		}
	}

}
