// Copyright 2020 Furqan Software Ltd. All rights reserved.

package py3

import "git.furqansoftware.net/toph/scanlib/gen/code"

type Context struct {
	types map[string]string
	cw    *code.Writer

	scan int
}
