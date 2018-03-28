package ui

import "image"

type Buffer struct {
	Area    image.Rectangle // selected drawing area
	CellMap map[image.Point]Cell
}

func NewBuffer() Buffer {
	return Buffer{
		CellMap: make(map[image.Point]Cell),
		Area:    image.Rectangle{}}
}

func (b *Buffer) SetArea(r image.Rectangle) {
	b.Area.Max = r.Max
	b.Area.Min = r.Min
}

func (b Buffer) Fill(ch rune, fg, bg Attribute) {
	for x := b.Area.Min.X; x < b.Area.Max.X; x++ {
		for y := b.Area.Min.Y; y < b.Area.Max.Y; y++ {
			b.Set(x, y, Cell{ch, fg, bg})
		}
	}
}

func (b Buffer) Set(x, y int, c Cell) {
	b.CellMap[image.Pt(x, y)] = c
}

func (b *Buffer) Merge(bs ...Buffer) {
	for _, buf := range bs {
		for p, v := range buf.CellMap {
			b.Set(p.X, p.Y, v)
		}
		b.SetArea(b.Area.Union(buf.Area))
	}
}

func NewFilledBuffer(x0, y0, x1, y1 int, ch rune, fg, bg Attribute) Buffer {
	buf := NewBuffer()
	buf.Area.Min = image.Pt(x0, y0)
	buf.Area.Max = image.Pt(x1, y1)

	for x := buf.Area.Min.X; x < buf.Area.Max.X; x++ {
		for y := buf.Area.Min.Y; y < buf.Area.Max.Y; y++ {
			buf.Set(x, y, Cell{ch, fg, bg})
		}
	}
	return buf
}

