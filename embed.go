package raven

import "embed"

//go:embed web/dist/*
var FrontendDist embed.FS
