package graph

import (
	"context"

	"github.com/skyrocketOoO/RBAC-server/domain"
	"github.com/skyrocketOoO/go-utility/queue"
	"github.com/skyrocketOoO/go-utility/set"
)

type GraphInfra struct {
	dbRepo domain.DbRepository
}

func NewGraphInfra(dbRepo domain.DbRepository) *GraphInfra {
	return &GraphInfra{
		dbRepo: dbRepo,
	}
}

func (g *GraphInfra) Check(c context.Context, start domain.Vertex, target domain.Vertex,
	relation string, searchCond domain.SearchCond) (bool, error) {
	visited := set.NewSet[domain.Vertex]()
	q := queue.NewQueue[domain.Vertex]()
	visited.Add(start)
	q.Push(start)
	for !q.IsEmpty() {
		qLen := q.Len()
		for i := 0; i < qLen; i++ {
			vertex, _ := q.Pop()
			edges, err := g.dbRepo.Get(c, domain.Edge{
				UNs:   vertex.Ns,
				UName: vertex.Name,
			}, true)
			if err != nil {
				return false, err
			}

			for _, edge := range edges {
				if edge.VNs == target.Ns && edge.VName == target.Name &&
					edge.Rel == relation {
					return true, nil
				}
				child := domain.Vertex{
					Ns:   edge.VNs,
					Name: edge.VName,
				}
				if !visited.Exist(child) {
					visited.Add(child)
					q.Push(child)
				}
			}
		}
	}

	return false, nil
}

func (g *GraphInfra) SearchPermissions(c context.Context, start domain.Vertex,
	isU bool, searchCond domain.SearchCond, collectCond domain.CollectCond,
	maxDepth int) ([]domain.Permission, error) {
	if isU {
		depth := 0
		pSet := set.NewSet[domain.Permission]()
		visited := set.NewSet[domain.Vertex]()
		q := queue.NewQueue[domain.Vertex]()
		visited.Add(start)
		q.Push(start)
		for !q.IsEmpty() {
			qLen := q.Len()
			for i := 0; i < qLen; i++ {
				vertex, _ := q.Pop()
				query := domain.Edge{
					UNs:   vertex.Ns,
					UName: vertex.Name,
				}
				qEdges, err := g.dbRepo.Get(c, query, true)
				if err != nil {
					return nil, err
				}

				for _, edge := range qEdges {
					child := domain.Vertex{
						Ns:   edge.VNs,
						Name: edge.VName,
					}
					if collectCond.ShouldCollect(child) {
						pSet.Add(domain.Permission{
							Ns:   edge.VNs,
							Name: edge.VName,
							Rel:  edge.Rel,
						})
					}
					if !searchCond.ShouldStop(child) &&
						!visited.Exist(child) {
						visited.Add(child)
						q.Push(child)
					}
				}
			}
			depth++
			if depth >= maxDepth {
				break
			}
		}

		return pSet.ToSlice(), nil
	} else {
		depth := 0
		pSet := set.NewSet[domain.Permission]()
		visited := set.NewSet[domain.Vertex]()
		q := queue.NewQueue[domain.Vertex]()
		visited.Add(start)
		q.Push(start)
		for !q.IsEmpty() {
			qLen := q.Len()
			for i := 0; i < qLen; i++ {
				vertex, _ := q.Pop()
				query := domain.Edge{
					VNs:   vertex.Ns,
					VName: vertex.Name,
				}
				qEdges, err := g.dbRepo.Get(c, query, true)
				if err != nil {
					return nil, err
				}

				for _, edge := range qEdges {
					parent := domain.Vertex{
						Ns:   edge.UNs,
						Name: edge.UName,
					}
					if collectCond.ShouldCollect(parent) {
						pSet.Add(domain.Permission{
							Ns:   edge.UNs,
							Name: edge.UName,
							Rel:  edge.Rel,
						})
					}
					if !searchCond.ShouldStop(parent) &&
						!visited.Exist(parent) {
						visited.Add(parent)
						q.Push(parent)
					}
				}
			}
			depth++
			if depth >= maxDepth {
				break
			}
		}

		return pSet.ToSlice(), nil
	}
}

func (g *GraphInfra) SearchVertices(c context.Context, start domain.Vertex,
	isU bool, searchCond domain.SearchCond, collectCond domain.CollectCond,
	maxDepth int) ([]domain.Vertex, error) {
	if isU {
		depth := 0
		vSet := set.NewSet[domain.Vertex]()
		visited := set.NewSet[domain.Vertex]()
		q := queue.NewQueue[domain.Vertex]()
		visited.Add(start)
		q.Push(start)
		for !q.IsEmpty() {
			qLen := q.Len()
			for i := 0; i < qLen; i++ {
				vertex, _ := q.Pop()
				query := domain.Edge{
					UNs:   vertex.Ns,
					UName: vertex.Name,
				}
				qEdges, err := g.dbRepo.Get(c, query, true)
				if err != nil {
					return nil, err
				}

				for _, edge := range qEdges {
					child := domain.Vertex{
						Ns:   edge.VNs,
						Name: edge.VName,
					}
					if collectCond.ShouldCollect(child) {
						vSet.Add(child)
					}
					if !searchCond.ShouldStop(child) &&
						!visited.Exist(child) {
						visited.Add(child)
						q.Push(child)
					}
				}
			}
			depth++
			if depth >= maxDepth {
				break
			}
		}

		return vSet.ToSlice(), nil
	} else {
		depth := 0
		vSet := set.NewSet[domain.Vertex]()
		visited := set.NewSet[domain.Vertex]()
		q := queue.NewQueue[domain.Vertex]()
		visited.Add(start)
		q.Push(start)
		for !q.IsEmpty() {
			qLen := q.Len()
			for i := 0; i < qLen; i++ {
				vertex, _ := q.Pop()
				query := domain.Edge{
					VNs:   vertex.Ns,
					VName: vertex.Name,
				}
				qEdges, err := g.dbRepo.Get(c, query, true)
				if err != nil {
					return nil, err
				}

				for _, edge := range qEdges {
					parent := domain.Vertex{
						Ns:   edge.UNs,
						Name: edge.UName,
					}
					if collectCond.ShouldCollect(parent) {
						vSet.Add(parent)
					}
					if !searchCond.ShouldStop(parent) &&
						!visited.Exist(parent) {
						visited.Add(parent)
						q.Push(parent)
					}
				}
			}
			depth++
			if depth >= maxDepth {
				break
			}
		}

		return vSet.ToSlice(), nil
	}
}

// func (g *GraphInfra) GetPassedVertices(c context.Context, start domain.Vertex,
// 	isU bool, searchCond domain.SearchCond, collectCond domain.CollectCond,
// 	maxDepth int) ([]domain.Vertex, error) {
// 	if isU {
// 		// if err := utils.ValidateVertex(start, true); err != nil {
// 		// 	return nil, err
// 		// }
// 		depth := 0
// 		verticesSet := set.NewSet[domain.Vertex]()
// 		visited := set.NewSet[domain.Vertex]()
// 		q := queue.NewQueue[domain.Vertex]()
// 		visited.Add(start)
// 		q.Push(start)
// 		for !q.IsEmpty() {
// 			qLen := q.Len()
// 			for i := 0; i < qLen; i++ {
// 				vertex, _ := q.Pop()
// 				query := domain.Edge{
// 					UNs:   vertex.Ns,
// 					UName: vertex.Name,
// 					URel:  vertex.Rel,
// 				}
// 				qEdges, err := g.dbRepo.Get(c, query, true)
// 				if err != nil {
// 					return nil, err
// 				}

// 				for _, edge := range qEdges {
// 					child := domain.Vertex{
// 						Ns:   edge.VNs,
// 						Name: edge.VName,
// 						Rel:  edge.VRel,
// 					}
// 					if collectCond.ShouldCollect(child) {
// 						verticesSet.Add(child)
// 					}
// 					if !searchCond.ShouldStop(child) &&
// 						!visited.Exist(child) {
// 						visited.Add(child)
// 						q.Push(child)
// 					}
// 				}
// 			}
// 			depth++
// 			if depth >= maxDepth {
// 				break
// 			}
// 		}

// 		return verticesSet.ToSlice(), nil
// 	} else {
// 		// if err := utils.ValidateVertex(start, false); err != nil {
// 		// 	return nil, err
// 		// }
// 		depth := 0
// 		verticesSet := set.NewSet[domain.Vertex]()
// 		visited := set.NewSet[domain.Vertex]()
// 		q := queue.NewQueue[domain.Vertex]()
// 		visited.Add(start)
// 		q.Push(start)
// 		for !q.IsEmpty() {
// 			qLen := q.Len()
// 			for i := 0; i < qLen; i++ {
// 				vertex, _ := q.Pop()
// 				query := domain.Edge{
// 					VNs:   vertex.Ns,
// 					VName: vertex.Name,
// 					VRel:  vertex.Rel,
// 				}
// 				qEdges, err := g.dbRepo.Get(c, query, true)
// 				if err != nil {
// 					return nil, err
// 				}

// 				for _, edge := range qEdges {
// 					parent := domain.Vertex{
// 						Ns:   edge.UNs,
// 						Name: edge.UName,
// 						Rel:  edge.URel,
// 					}
// 					if collectCond.ShouldCollect(parent) {
// 						verticesSet.Add(parent)
// 					}
// 					if !searchCond.ShouldStop(parent) &&
// 						!visited.Exist(parent) {
// 						visited.Add(parent)
// 						q.Push(parent)
// 					}
// 				}
// 			}
// 			depth++
// 			if depth >= maxDepth {
// 				break
// 			}
// 		}

// 		return verticesSet.ToSlice(), nil
// 	}
// }

func (g *GraphInfra) GetTree(c context.Context, start domain.Vertex, maxDepth int) (
	*domain.TreeNode, error) {
	if res, err := g.dbRepo.Get(
		c,
		domain.Edge{
			UNs:   start.Ns,
			UName: start.Name,
		}, true); err != nil {
		return nil, err
	} else if len(res) == 0 {
		return nil, domain.ErrRecordNotFound
	}

	root := &domain.TreeNode{
		Ns:       start.Ns,
		Name:     start.Name,
		Children: map[string]*domain.TreeNode{},
	}
	visited := map[domain.Vertex]*domain.TreeNode{}
	visited[start] = root
	q := queue.NewQueue[*domain.TreeNode]()
	q.Push(root)
	for depth := 1; depth <= maxDepth && !q.IsEmpty(); depth++ {
		for i := 0; i < q.Len(); i++ {
			u, err := q.Pop()
			if err != nil {
				return nil, err
			}
			edges, err := g.dbRepo.Get(c,
				domain.Edge{
					UNs:   u.Ns,
					UName: u.Name,
				},
				true,
			)
			if err != nil {
				return nil, err
			}
			for _, edge := range edges {
				v := domain.Vertex{
					Ns:   edge.VNs,
					Name: edge.VName,
				}
				if node, ok := visited[v]; !ok {
					newNode := &domain.TreeNode{
						Ns:       v.Ns,
						Name:     v.Name,
						Children: map[string]*domain.TreeNode{},
					}
					q.Push(newNode)
					visited[v] = newNode
					visited[domain.Vertex{
						Ns: u.Ns, Name: u.Name}].Children[edge.Rel] = newNode
				} else {
					visited[domain.Vertex{
						Ns: u.Ns, Name: u.Name}].Children[edge.Rel] = node
				}
			}
		}
	}
	return root, nil
}
