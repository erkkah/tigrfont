module github.com/erkkah/tigrfont

go 1.16

require (
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0
	github.com/nfnt/resize v0.0.0-20180221191011-83c6a9932646
	github.com/zachomedia/go-bdf v0.0.0-20210522061406-1a147053be95
	golang.org/x/image v0.0.0-20210628002857-a66eb6448b8d
	golang.org/x/text v0.3.7
)

//replace github.com/zachomedia/go-bdf => ../go-bdf
