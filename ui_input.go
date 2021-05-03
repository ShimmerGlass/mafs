package main

type uiInput struct {
	buffer []rune
	Cursor int
}

func (i *uiInput) Add(r rune) {
	var left, right []rune

	left = i.buffer[:i.Cursor]
	if i.Cursor <= len(i.buffer)-1 {
		right = i.buffer[i.Cursor:]
	}

	i.buffer = append(left, append([]rune{r}, right...)...)
	i.Cursor++
}

func (i *uiInput) Del() {
	if i.Cursor == 0 {
		return
	}
	copy(i.buffer[i.Cursor-1:], i.buffer[i.Cursor:])
	i.buffer = i.buffer[:len(i.buffer)-1]
	i.Cursor--
}

func (i *uiInput) DelFwd() {
	if i.Cursor >= len(i.buffer) {
		return
	}
	copy(i.buffer[i.Cursor:], i.buffer[i.Cursor+1:])
	i.buffer = i.buffer[:len(i.buffer)-1]
}

func (i *uiInput) MoveCursor(n int) {
	i.Cursor += n
	if i.Cursor < 0 {
		i.Cursor = 0
	}
	if i.Cursor > len(i.buffer) {
		i.Cursor = len(i.buffer)
	}
}

func (i *uiInput) Clear() {
	i.buffer = nil
	i.Cursor = 0
}

func (i *uiInput) Text() string {
	return string(i.buffer)
}

func (i *uiInput) Runes() []rune {
	return i.buffer
}
