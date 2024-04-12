package usecase

import (
	"context"
	"math"

	"github.com/skyrocketOoO/RBAC-server/domain"
)

type Usecase struct {
	dbRepo     domain.DbRepository
	graphInfra domain.GraphInfra
}

func NewUsecase(dbRepo domain.DbRepository, graphInfra domain.GraphInfra) *Usecase {
	return &Usecase{
		dbRepo:     dbRepo,
		graphInfra: graphInfra,
	}
}

func (u *Usecase) DeleteUser(c context.Context, name string) error {
	return u.dbRepo.Delete(c, domain.Edge{UNs: "user", UName: name}, true)
}

func (u *Usecase) UserGetPermissions(c context.Context, name string) (
	[]domain.Permission, error) {
	return u.graphInfra.GetPaths(
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
	rels, err := u.graphInfra.GetPaths(
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
	roles := make([]string, len(rels))
	for i, rel := range rels {
		roles[i] = rel.Name
	}
	return roles, nil
}

func (u *Usecase) UserCheck(c context.Context, username string, objNs string,
	objName string) (bool, error) {
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
	return u.graphInfra.GetPaths(
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
	permissions, err := u.graphInfra.GetPaths(
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
	for
}

func (u *Usecase) WhichUserHasPermission(c context.Context, objNs string,
	objName string) {

}

func (u *Usecase) DeletePermission(c context.Context, permissionName string) {}
