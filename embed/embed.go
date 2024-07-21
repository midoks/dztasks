package embed

import (
	"embed"
)

//go:embed static/*
var Static embed.FS


//go:embed templates
var Templates embed.FS


//go:embed conf/*
var Conf embed.FS