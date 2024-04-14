package api

import (
	"github.com/gin-gonic/gin"
	"github.com/skyrocketOoO/RBAC-server/internal/delivery/rest"
)

func Binding(r *gin.Engine, d *rest.RestDelivery) {
	r.GET("/ping", d.Ping)
	r.GET("/healthy", d.Healthy)
	// r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	userR := r.Group("/user")
	{
		userR.DELETE("/:name", d.DeleteUser)
		userR.GET("/:name/permission", d.UserGetPermissions)
		userR.GET("/:name/role", d.UserGetRoles)
		userR.GET("/:name/check/:rel/:objns/:objname", d.UserCheck)
		userR.POST("/:name/permission", d.UserAddPermission)
		userR.DELETE("/:name/permission", d.UserRemovePermission)
		userR.POST("/:name/role", d.UserAddRole)
		userR.DELETE("/:name/role", d.UserRemoveRole)
	}
	roleR := r.Group("/role")
	{
		roleR.DELETE("/:name", d.DeleteRole)
		roleR.GET("/:name/user", d.RoleGetUsers)
		roleR.GET("/:name/permission", d.RoleGetPermissions)
		roleR.POST("/:name/permission", d.RoleAddPermission)
		roleR.DELETE("/:name/permission", d.RoleRemovePermission)
		roleR.POST("/:name/inherit", d.RoleInheritRole)
		roleR.DELETE("/:name/inherit", d.RoleUnInheritRole)
		roleR.GET("/:name/child", d.RoleGetChildRole)
		roleR.GET("/:name/parent", d.RoleGetParentRole)
	}
	objectR := r.Group("/object")
	{
		objectR.DELETE("/:ns/:name", d.DeleteObject)
		objectR.GET("/:ns/:name/role", d.WhichRoleHasPermission)
		objectR.GET("/:ns/:name/user", d.WhichUserHasPermission)
	}
}
