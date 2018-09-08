//go:generate go get -v github.com/mjibson/esc
//go:generate esc -o cmd/static_content.go -pkg cmd -private LICENSE.TXT
//go:generate esc -o pkg/static_content.go -pkg srcmanager -private repositories

package main

import (
	_ "github.com/mjibson/esc/embed"
)
