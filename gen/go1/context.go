// Copyright 2020 Furqan Software Ltd. All rights reserved.

package go1

import "git.furqansoftware.net/toph/scanlib/gen/code"

type Context struct {
	types   map[string]string
	imports map[string]bool
	cw      *code.Writer
}
