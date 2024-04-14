package api

import (
	"github.com/gin-gonic/gin"
	"github.com/skyrocketOoO/RBAC-server/internal/delivery/rest"
)

func Binding(r *gin.Engine, d *rest.RestDelivery) {
	r.GET("/ping", d.Ping)
	r.GET("/healthy", d.Healthy)
	// r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/ping", d.Ping)
	r.GET("/healthy", d.Healthy)

	userR := r.Group("/user")
	{
		userR.DELETE("/", d.DeleteUser)
		userR.GET("/permission", d.UserGetPermissions)
		userR.GET("/role", d.UserGetRoles)
		userR.GET("/check", d.UserCheck)
		userR.POST("/permission", d.UserAddPermission)
		userR.DELETE("/permission", d.UserRemovePermission)
		userR.POST("/role", d.UserAddRole)
		userR.DELETE("/role", d.UserRemoveRole)
	}
	roleR := r.Group("/role")
	{
		roleR.DELETE("/", d.DeleteRole)
		roleR.GET("/user", d.RoleGetUsers)
		roleR.GET("/permission", d.RoleGetPermissions)
		roleR.POST("/permission", d.RoleAddPermission)
		roleR.DELETE("/permission", d.RoleRemovePermission)
		roleR.POST("/inherit", d.RoleInheritRole)
		roleR.DELETE("/inherit", d.RoleUnInheritRole)
		roleR.GET("/child", d.RoleGetChildRole)
		roleR.GET("/parent", d.RoleGetParentRole)
	}
	objectR := r.Group("/object")
	{
		objectR.DELETE("/", d.DeleteObject)
		objectR.GET("/role", d.WhichRoleHasPermission)
		objectR.GET("/user", d.WhichUserHasPermission)
	}
}
