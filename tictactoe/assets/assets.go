package assets

import (
	"embed"
)

//go:embed images/*
var Assets embed.FS

//go:embed fonts/*
var Fonts embed.FS
