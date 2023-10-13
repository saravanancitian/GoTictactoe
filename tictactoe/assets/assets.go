package assets

import (
	"embed"
)

//go:embed images/*
var Assets embed.FS

//go:embed fonts/*
var Fonts embed.FS

//go:embed audio/*
var Audios embed.FS
