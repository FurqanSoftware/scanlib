// Copyright 2020 Furqan Software Ltd. All rights reserved.

package py3

import (
	"git.furqansoftware.net/toph/scanlib/ast"
)

func analyzeSource(ctx *Context, n *ast.Source) error {
	return analyzeBlock(ctx, &n.Block)
}

func analyzeBlock(ctx *Context, n *ast.Block) error {
	analyzeBlockScanSame(ctx, n)
	analyzeBlockScanOne(ctx, n)
	return nil
}
