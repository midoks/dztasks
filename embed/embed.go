package embed

import (
	"embed"
)

//go:embed static/*
var Static embed.FS


//go:embed tmplate/*
var Templates embed.FS