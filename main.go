package main

import (
	"WinAMP/audio"
	"WinAMP/ui"
	"path/filepath"
	"strconv"
	"runtime"
)

var (
	par *ui.Par
	par2 *ui.Par
	g *ui.Gauge
)

type playlist struct{
	arr []song
	list []string
	select_song int
	play_song int
	time_song int
}

type song struct {
	file   string
	name string
}

var NumSelectSong = 0

func main() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	pl := playlist{
		arr: make([]song, 0, 5),
		list: make([]string, 0, 5),
		select_song: -1,
		play_song: -1,
	}
	pl.AddFile("1.mp3")
	pl.AddFile("2.mp3")
	pl.AddFile("3.mp3")



	player := audio.New()
	defer player.Close()

	err := ui.Init(80, 27)
	if err != nil {
		panic(err)
	}
	defer ui.Close()
	ui.SetTitle("Player WinAMP")


	//strs := make([]string, 0, len(songs))
	//for i, s :=range songs {
	//	f := filepath.Base(s.file)
	//	strs = append(strs, "["+strconv.Itoa(i+1)+"] "+f)
	//}
	pl.SetActive(0)

	ls := ui.NewList()
	ls.Items = pl.list
	ls.ItemFgColor = ui.ColorYellow
	ls.BorderLabel = " Список песен "
	ls.Height = 20
	ls.Width = 80
	ls.Y = 0

	par = ui.NewPar("Пусто")
	par.Height = 3
	par.Width = 80
	par.Y = 20
	par.BorderLabel = " Сейчас играет песня "

	par2 = ui.NewPar(strconv.Itoa(NumSelectSong))
	par2.Height = 3
	par2.Width = 9
	par2.Y = 24
	par2.BorderLabel = " Кол-во "

	g = ui.NewGauge()
	g.Percent = 0
	g.Width = 70
	g.Height = 3
	g.Y = 24
	g.X = 10
	g.BorderLabel = ""
	g.Label = "{{percent}}% (милисекунд)"
	g.LabelAlign = ui.AlignRight

	ui.Render(ls, par, par2, g)



	//ui.Body.AddRows(
	//	ui.NewRow(
	//		ui.NewCol(13, 0, ls)))


	//ui.Body.Align()
	//ui.Render(ui.Body)


	ui.Handle(ui.EventResize, func(e ui.Event) {
		//w, _ := ui.Size()
		//ui.Body.Width = w
		//ui.Body.Align()
		////ui.Clear()
		//ui.Render(ui.Body)
	})

	ui.Handle(ui.EventTimer, func(e ui.Event) {
		if player.IsPaying() {
			pos := player.CurrentPos()
			pecent := 0
			if pl.time_song != 0 {
				pecent = int(float32(pos) / float32(pl.time_song) * 100)
			}
			if g.Percent != pecent {
				g.Percent = pecent
				ui.Render(g)
			}
		}else if !player.IsPause(){
			pl.PlayNext(player)
			ui.Render(ls, par, par2, g)
		}
	})

	ui.Handle(ui.EventKey, func(ev ui.Event) {
		if ev.Key == ui.KeyEsc {
			ui.StopLoop()
		}else if ev.Key == ui.KeyArrowDown {
			pl.SetActive(pl.select_song + 1)
			ui.Render(ls)
		}else if ev.Key == ui.KeyArrowUp{
			pl.SetActive(pl.select_song-1)
			ui.Render(ls)
		}else if ev.Key == ui.KeyEnter{
			if player.IsPause() {
				player.Play()
			}else {
				pl.Play(player, pl.select_song)
				ui.Render(ls, par, par2, g)
			}
		}else if ev.Key == ui.KeyCtrlP {
			if player.IsPause() {
				player.Play()
			}else{
				player.Pause()
			}
		}
	})
	ui.Loop()


	//err := audio.OpenFile("D:\\C++Project\\WinAMP\\Debug\\m.mp3")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println("Sleep")
	//time.Sleep(time.Second * 10)
	//
	//if audio.IsPaying() {
	//	log.Println("IsPaying")
	//}else{
	//	log.Println("No IsPaying")
	//}
	//
	//audio.Pause()
	//time.Sleep(time.Second * 3)
	//
	//if audio.IsPause() {
	//	log.Println("IsPause")
	//}else{
	//	log.Println("No IsPause")
	//}
	//
	//audio.Play()
	//time.Sleep(time.Second * 10)

}

func (pl *playlist) SetActive(i int)  {
	if i >=0 && i < len(pl.list) {
		oldActive := pl.select_song
		pl.select_song = i
		pl.paint_song(oldActive)
		pl.paint_song(pl.select_song)
	}
}

func (pl *playlist) AddFile(path string)  {
	pl.arr = append(pl.arr, song{file:path, name:filepath.Base(path)})
	id := len(pl.arr)
	pl.list = append(pl.list, "[" + strconv.Itoa(id) + "] "+pl.arr[id-1].name)
}

func (pl *playlist) Play(player *audio.Player, i int)  {
	if i >=0 && i < len(pl.arr) && pl.play_song != i {
		NumSelectSong ++
		par2.Text = strconv.Itoa(NumSelectSong)
		g.Percent = 0



		oldActive := pl.play_song
		pl.paint_song(oldActive)

		pl.play_song = i

		f:=pl.arr[i].file
		err := player.OpenFile(f)
		if err != nil {
			pl.play_song = -1
			par.Text = "Ошибка: "+err.Error()
			return
		}
		pl.paint_song(pl.play_song)
		par.Text = filepath.Base(f)

		pl.time_song = player.SongTime()
	}
}



func (pl *playlist) paint_song(i int)  {
	if i >=0 && i < len(pl.list) {
		//f := filepath.Base(songs[i].file)
		str := "[" + strconv.Itoa(i+1) + "] "
		if pl.select_song == i || pl.play_song == i {
			str += "[" + pl.arr[i].name + "]("
			if pl.play_song == i {
				str += "fg-red"
			} else {
				str += "fg-white"
			}
			if pl.select_song == i {
				str += ",bg-green"
			}
			str += ")"
		} else {
			str += pl.arr[i].name
		}
		pl.list[i]=str
	}
}

func (pl *playlist) PlayNext(player *audio.Player)  {
	newId := pl.play_song
	pl.play_song = -1
	pl.paint_song(newId)
	if len(pl.arr) == 0 {
		return
	}
	newId ++
	if newId >= len(pl.arr) {
		newId = 0
	}
	pl.Play(player, newId)
}
