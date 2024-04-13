package usecase

import (
	"context"
	"math"

	"github.com/skyrocketOoO/RBAC-server/domain"
	"go.mongodb.org/mongo-driver/mongo"
)

type Usecase struct {
	mongoClient *mongo.Client
	graphInfra  domain.GraphInfra
	dbRepo      domain.DbRepository
}

func NewUsecase(mongoCli *mongo.Client, graphInfra domain.GraphInfra,
	dbRepo domain.DbRepository) *Usecase {
	return &Usecase{
		mongoClient: mongoCli,
		graphInfra:  graphInfra,
	}
}

func (u *Usecase) Healthy(c context.Context) error {
	// do something check like db connection is established
	if err := u.dbRepo.Ping(c); err != nil {
		return err
	}

	return nil
}

func (u *Usecase) DeleteUser(c context.Context, name string) error {
	return u.dbRepo.Delete(c, domain.Edge{UNs: "user", UName: name}, true)
}

func (u *Usecase) UserGetPermissions(c context.Context, name string) (
	[]domain.Permission, error) {
	return u.graphInfra.SearchPermissions(
		c,
		domain.Vertex{
			Ns:   "user",
			Name: name,
		},
		true,
		domain.SearchCond{},
		domain.CollectCond{
			NotIn: domain.Compare{
				Nses: []string{"role", "user"},
			},
		},
		math.MaxInt)
}

func (u *Usecase) UserGetRoles(c context.Context, name string) (
	[]string, error) {
	vertices, err := u.graphInfra.SearchVertices(
		c,
		domain.Vertex{
			Ns:   "user",
			Name: name,
		},
		true,
		domain.SearchCond{
			In: domain.Compare{
				Nses: []string{"role"},
			},
		},
		domain.CollectCond{
			In: domain.Compare{
				Nses: []string{"role"},
			},
		},
		math.MaxInt)
	if err != nil {
		return nil, err
	}
	roles := make([]string, len(vertices))
	for i, rel := range vertices {
		roles[i] = rel.Name
	}
	return roles, nil
}

func (u *Usecase) UserCheck(c context.Context, username string, objNs string,
	relation string, objName string) (bool, error) {
	return u.graphInfra.Check(
		c,
		domain.Vertex{
			Ns:   "user",
			Name: username,
		},
		domain.Vertex{
			Ns:   objNs,
			Name: objName,
		},
		relation,
		domain.SearchCond{},
	)
}

func (u *Usecase) UserAddPermission(c context.Context, username string,
	permission domain.Permission) error {
	return u.dbRepo.Create(c, domain.Edge{
		UNs:   "user",
		UName: username,
		Rel:   permission.Rel,
		VNs:   permission.Ns,
		VName: permission.Name,
	})
}

func (u *Usecase) UserRemovePermission(c context.Context, username string,
	permission domain.Permission) error {
	return u.dbRepo.Delete(c, domain.Edge{
		UNs:   "user",
		UName: username,
		Rel:   permission.Rel,
		VNs:   permission.Ns,
		VName: permission.Name,
	}, false)
}

func (u *Usecase) UserAddRole(c context.Context, username string,
	roleName string) error {
	return u.dbRepo.Create(c, domain.Edge{
		UNs:   "user",
		UName: username,
		Rel:   "member",
		VNs:   "role",
		VName: roleName,
	})
}

func (u *Usecase) UserRemoveRole(c context.Context, username string,
	roleName string) error {
	return u.dbRepo.Delete(c, domain.Edge{
		UNs:   "user",
		UName: username,
		Rel:   "member",
		VNs:   "role",
		VName: roleName,
	}, false)
}

func (u *Usecase) DeleteRole(c context.Context, name string) error {
	err := u.dbRepo.Delete(c, domain.Edge{
		UNs:   "role",
		UName: name,
	}, true)
	if err != nil {
		return err
	}
	return u.dbRepo.Delete(c, domain.Edge{
		VNs:   "role",
		VName: name,
	}, true)
}

func (u *Usecase) RoleGetUsers(c context.Context, name string) ([]string, error) {
	edges, err := u.dbRepo.Get(c, domain.Edge{Rel: "member", VNs: "role", VName: name},
		true)
	if err != nil {
		return nil, err
	}
	users := make([]string, len(edges))
	for i, edge := range edges {
		users[i] = edge.UName
	}
	return users, nil
}

func (u *Usecase) RoleGetPermissions(c context.Context, name string) (
	[]domain.Permission, error) {
	return u.graphInfra.SearchPermissions(
		c,
		domain.Vertex{
			Ns:   "role",
			Name: name,
		},
		true,
		domain.SearchCond{},
		domain.CollectCond{
			NotIn: domain.Compare{
				Nses: []string{"role", "user"},
			},
		},
		math.MaxInt)
}

func (u *Usecase) RoleAddPermission(c context.Context, roleName string,
	permission domain.Permission) error {
	return u.dbRepo.Create(c, domain.Edge{
		UNs:   "role",
		UName: roleName,
		Rel:   permission.Rel,
		VNs:   permission.Ns,
		VName: permission.Name,
	})
}

func (u *Usecase) RoleRemovePermission(c context.Context, roleName string,
	permission domain.Permission) error {
	return u.dbRepo.Delete(c, domain.Edge{
		UNs:   "role",
		UName: roleName,
		Rel:   permission.Rel,
		VNs:   permission.Ns,
		VName: permission.Name,
	}, false)
}

func (u *Usecase) RoleInheritRole(c context.Context, parentName string,
	childName string) error {
	return u.dbRepo.Create(c, domain.Edge{
		UNs:   "role",
		UName: parentName,
		Rel:   "parent",
		VNs:   "role",
		VName: childName,
	})
}

func (u *Usecase) RoleUnInheritRole(c context.Context, parentName string,
	childName string) error {
	return u.dbRepo.Delete(c, domain.Edge{
		UNs:   "role",
		UName: parentName,
		Rel:   "parent",
		VNs:   "role",
		VName: childName,
	}, false)
}

func (u *Usecase) RoleGetChildRole(c context.Context, name string) (
	[]string, error) {
	edges, err := u.dbRepo.Get(c, domain.Edge{
		UNs:   "role",
		UName: name,
		Rel:   "parent",
		VNs:   "role",
	}, true)
	if err != nil {
		return nil, err
	}

	roles := make([]string, len(edges))
	for i, edge := range edges {
		roles[i] = edge.VName
	}
	return roles, nil
}

func (u *Usecase) RoleGetParentRole(c context.Context, name string) (
	[]string, error) {
	edges, err := u.dbRepo.Get(c, domain.Edge{
		UNs:   "role",
		Rel:   "parent",
		VNs:   "role",
		VName: name,
	}, true)
	if err != nil {
		return nil, err
	}

	roles := make([]string, len(edges))
	for i, edge := range edges {
		roles[i] = edge.VName
	}
	return roles, nil
}

func (u *Usecase) DeleteObject(c context.Context, ns string,
	name string) error {
	return u.dbRepo.Delete(c, domain.Edge{VNs: ns, VName: name}, true)
}

func (u *Usecase) WhichRoleHasPermission(c context.Context, objNs string,
	objName string) ([]string, error) {
	vertices, err := u.graphInfra.SearchVertices(
		c,
		domain.Vertex{
			Ns:   objNs,
			Name: objName,
		},
		false,
		domain.SearchCond{},
		domain.CollectCond{
			NotIn: domain.Compare{
				Nses: []string{"role"},
			},
		},
		math.MaxInt)
	if err != nil {
		return nil, err
	}
	roles := make([]string, len(vertices))
	for i, v := range vertices {
		roles[i] = v.Name
	}
	return roles, nil
}

func (u *Usecase) WhichUserHasPermission(c context.Context, objNs string,
	objName string) ([]string, error) {
	vertices, err := u.graphInfra.SearchVertices(
		c,
		domain.Vertex{
			Ns:   objNs,
			Name: objName,
		},
		false,
		domain.SearchCond{},
		domain.CollectCond{
			NotIn: domain.Compare{
				Nses: []string{"user"},
			},
		},
		math.MaxInt)
	if err != nil {
		return nil, err
	}
	users := make([]string, len(vertices))
	for i, v := range vertices {
		users[i] = v.Name
	}
	return users, nil
}

// func (u *Usecase) DeletePermission(c context.Context, name string) {
// 	u.dbRepo.Delete(c, domain.Edge{Rel: "permission", VName: permissionName},
// 		true)
// }
