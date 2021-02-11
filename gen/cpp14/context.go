// Copyright 2020 Furqan Software Ltd. All rights reserved.

package cpp14

import "git.furqansoftware.net/toph/scanlib/gen/code"

type Context struct {
	types    map[string]string
	includes map[string]bool
	cw       *code.Writer
}
