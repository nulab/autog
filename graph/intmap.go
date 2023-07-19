package graph

type intmap[T any] map[*T]int

type NodeMap = intmap[Node]

type EdgeMap = intmap[Edge]

// func (m intmap[T]) String() string {
// 	type pair struct {
// 		n *Node
// 		i int
// 	}
// 	var kvPairs []pair
// 	for k, v := range m {
// 		kvPairs = append(kvPairs, pair{k, v})
// 	}
// 	slices.SortFunc(kvPairs, func(a, b pair) bool {
// 		return a.i > b.i
// 	})
// 	bld := strings.Builder{}
// 	for _, p := range kvPairs {
// 		bld.WriteRune('[')
// 		bld.WriteString(p.n.ID)
// 		bld.WriteRune(':')
// 		bld.WriteString(strconv.Itoa(p.i))
// 		bld.WriteRune(']')
// 		bld.WriteRune(' ')
// 	}
// 	return bld.String()
// }
