package domain

import (
	"context"
)

type DbRepository interface {
	Ping(c context.Context) error
	Get(c context.Context, edge Edge, queryMode bool) (edges []Edge, err error)
	Create(c context.Context, edge Edge) error
	Delete(c context.Context, edge Edge, queryMode bool) error
	ClearAll(c context.Context) error
}

type GraphInfra interface {
	Check(c context.Context, start Vertex, target Vertex, relation string,
		searchCond SearchCond) (found bool, err error)
	SearchPermissions(c context.Context, start Vertex, isSbj bool,
		searchCond SearchCond, collectCond CollectCond, maxDepth int) (
		permissions []Permission, err error)
	SearchVertices(c context.Context, start Vertex, isSbj bool,
		searchCond SearchCond, collectCond CollectCond, maxDepth int) (
		permissions []Vertex, err error)
	GetTree(c context.Context, sbj Vertex, maxDepth int) (*TreeNode, error)
}

type Usecase interface {
	Healthy(c context.Context) error
	DeleteUser(c context.Context, name string) error
	UserGetPermissions(c context.Context, name string) ([]Permission, error)
	UserGetRoles(c context.Context, name string) ([]string, error)
	UserCheck(c context.Context, username string, objNs string, relation string,
		objName string) (bool, error)
	UserAddPermission(c context.Context, username string, permission Permission) error
	UserRemovePermission(c context.Context, username string,
		permission Permission) error
	UserAddRole(c context.Context, username string, roleName string) error
	UserRemoveRole(c context.Context, username string, roleName string) error
	DeleteRole(c context.Context, name string) error
	RoleGetUsers(c context.Context, name string) ([]string, error)
	RoleGetPermissions(c context.Context, name string) ([]Permission, error)
	RoleAddPermission(c context.Context, roleName string, permission Permission) error
	RoleRemovePermission(c context.Context, roleName string,
		permission Permission) error
	RoleInheritRole(c context.Context, parentName string, childName string) error
	RoleUnInheritRole(c context.Context, parentName string, childName string) error
	RoleGetChildRole(c context.Context, name string) ([]string, error)
	RoleGetParentRole(c context.Context, name string) ([]string, error)
	DeleteObject(c context.Context, ns string, name string) error
	WhichRoleHasPermission(c context.Context, objNs string, objName string) (
		[]string, error)
	WhichUserHasPermission(c context.Context, objNs string, objName string) (
		[]string, error)
}
