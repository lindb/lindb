package terminal

import (
	"syscall"
	"unsafe"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

func GetTerminalWidth() int {
	var ws struct {
		Row    uint16
		Col    uint16
		Xpixel uint16
		Ypixel uint16
	}
	_, _, err := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(&ws)))
	if err != 0 {
		return 80 // default width
	}
	return int(ws.Col)
}

func TableSylte() table.Style {
	return table.Style{
		Name: "Custom",
		Box: table.BoxStyle{
			BottomLeft:       "+",
			BottomRight:      "+",
			BottomSeparator:  "+",
			Left:             "|",
			LeftSeparator:    "+",
			MiddleHorizontal: "-",
			MiddleSeparator:  "+",
			MiddleVertical:   "|",
			PaddingLeft:      "",
			PaddingRight:     "",
			Right:            "|",
			RightSeparator:   "+",
			TopLeft:          "+",
			TopRight:         "+",
			TopSeparator:     "+",
			UnfinishedRow:    "",
		},
		Color: table.ColorOptions{
			Header: text.Colors{text.Bold},
		},
		Options: table.Options{
			DrawBorder:      true,
			SeparateColumns: true,
			SeparateHeader:  true,
			SeparateRows:    true,
		},
	}
}
