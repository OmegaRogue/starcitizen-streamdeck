package main

import (
	"io"
	"runtime"

	"github.com/klauspost/compress/zstd"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"starcitizen-streamdeck/cmd"
	"starcitizen-streamdeck/internal/logger"
	"starcitizen-streamdeck/pkg/zip"
)

func ZstdWrapper(reader io.Reader) io.ReadCloser {
	r, err := zstd.NewReader(reader)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to read ZStd")
	}
	return r.IOReadCloser()
}

const location string = "/mnt/fastbulk/Games/star-citizen/drive_c/Program Files/Roberts Space Industries/StarCitizen/LIVE/Data.p4k"

const prefix string = "/mnt/fastbulk/Games/star-citizen"

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())
	zerolog.SetGlobalLevel(zerolog.TraceLevel)
	logger.SetupLogger()
	zip.RegisterDecompressor(100, ZstdWrapper)

	cmd.Execute()
}
