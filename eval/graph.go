package eval

func init() {
	Functions["graphs.simple"] = graphSimple
	Functions["graphs.connected"] = graphConnected
	Functions["graphs.acyclic"] = graphAcyclic
	Functions["graphs.tree"] = graphTree
}

func graphArgs(args []interface{}) (n int, u, v []int, err error) {
	if len(args) != 3 {
		return 0, nil, nil, ErrInvalidArgument{}
	}
	var ok bool
	n, ok = toInt(args[0])
	if !ok {
		return 0, nil, nil, ErrInvalidArgument{}
	}
	u, ok = args[1].([]int)
	if !ok {
		return 0, nil, nil, ErrInvalidArgument{}
	}
	v, ok = args[2].([]int)
	if !ok {
		return 0, nil, nil, ErrInvalidArgument{}
	}
	if len(u) != len(v) {
		return 0, nil, nil, ErrInvalidArgument{}
	}
	return n, u, v, nil
}

// graphSimple returns true if the graph has no self-loops or duplicate edges.
func graphSimple(args ...interface{}) (interface{}, error) {
	n, u, v, err := graphArgs(args)
	if err != nil {
		return nil, err
	}
	_ = n
	type edge struct{ a, b int }
	seen := map[edge]bool{}
	for i := range u {
		if u[i] == v[i] {
			return false, nil
		}
		a, b := u[i], v[i]
		if a > b {
			a, b = b, a
		}
		e := edge{a, b}
		if seen[e] {
			return false, nil
		}
		seen[e] = true
	}
	return true, nil
}

// graphConnected returns true if all n nodes are connected by the edges.
func graphConnected(args ...interface{}) (interface{}, error) {
	n, u, v, err := graphArgs(args)
	if err != nil {
		return nil, err
	}
	if n == 0 {
		return true, nil
	}
	parent := make([]int, n+1)
	rank := make([]int, n+1)
	for i := range parent {
		parent[i] = i
	}
	for i := range u {
		union(parent, rank, u[i], v[i])
	}
	root := find(parent, 1)
	for i := 2; i <= n; i++ {
		if find(parent, i) != root {
			return false, nil
		}
	}
	return true, nil
}

// graphAcyclic returns true if the graph has no cycles.
func graphAcyclic(args ...interface{}) (interface{}, error) {
	n, u, v, err := graphArgs(args)
	if err != nil {
		return nil, err
	}
	parent := make([]int, n+1)
	rank := make([]int, n+1)
	for i := range parent {
		parent[i] = i
	}
	for i := range u {
		if find(parent, u[i]) == find(parent, v[i]) {
			return false, nil
		}
		union(parent, rank, u[i], v[i])
	}
	return true, nil
}

// graphTree returns true if the edges form a tree on n nodes.
func graphTree(args ...interface{}) (interface{}, error) {
	n, u, v, err := graphArgs(args)
	if err != nil {
		return nil, err
	}
	if len(u) != n-1 {
		return false, nil
	}
	parent := make([]int, n+1)
	rank := make([]int, n+1)
	for i := range parent {
		parent[i] = i
	}
	for i := range u {
		if find(parent, u[i]) == find(parent, v[i]) {
			return false, nil
		}
		union(parent, rank, u[i], v[i])
	}
	return true, nil
}

func find(parent []int, x int) int {
	for parent[x] != x {
		parent[x] = parent[parent[x]]
		x = parent[x]
	}
	return x
}

func union(parent, rank []int, a, b int) {
	a = find(parent, a)
	b = find(parent, b)
	if a == b {
		return
	}
	if rank[a] < rank[b] {
		a, b = b, a
	}
	parent[b] = a
	if rank[a] == rank[b] {
		rank[a]++
	}
}
