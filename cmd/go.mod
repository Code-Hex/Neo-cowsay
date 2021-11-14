module github.com/Code-Hex/Neo-cowsay/v2/cmd

replace github.com/Code-Hex/Neo-cowsay/v2 => ../

go 1.17

require (
	github.com/Code-Hex/Neo-cowsay/v2 v2.0.0-00010101000000-000000000000
	github.com/Code-Hex/go-wordwrap v1.0.0
	github.com/jessevdk/go-flags v1.5.0
	github.com/ktr0731/go-fuzzyfinder v0.5.1
	github.com/mattn/go-colorable v0.1.11
	github.com/mattn/go-runewidth v0.0.13
	github.com/rivo/uniseg v0.2.0
	golang.org/x/crypto v0.0.0-20211108221036-ceb1ce70b4fa
)

require (
	github.com/gdamore/encoding v1.0.0 // indirect
	github.com/gdamore/tcell/v2 v2.4.0 // indirect
	github.com/lucasb-eyer/go-colorful v1.0.3 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/nsf/termbox-go v0.0.0-20201124104050-ed494de23a00 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	golang.org/x/sys v0.0.0-20210927094055-39ccf1dd6fa6 // indirect
	golang.org/x/term v0.0.0-20201210144234-2321bbc49cbf // indirect
	golang.org/x/text v0.3.6 // indirect
)
