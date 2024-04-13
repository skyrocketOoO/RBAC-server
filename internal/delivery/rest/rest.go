package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/skyrocketOoO/RBAC-server/domain"
)

type RestDelivery struct {
	usecase domain.Usecase
}

func NewDelivery(usecase domain.Usecase) *RestDelivery {
	return &RestDelivery{usecase: usecase}
}

// @Summary Check the server started
// @Accept json
// @Produce json
// @Success 200 {obj} domain.Response
// @Router /ping [get]
func (d *RestDelivery) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, domain.Response{Msg: "pong"})
}

// @Summary Check the server healthy
// @Accept json
// @Produce json
// @Success 200 {obj} domain.Response
// @Failure 503 {obj} domain.Response
// @Router /healthy [get]
func (d *RestDelivery) Healthy(c *gin.Context) {
	// do something check
	if err := d.usecase.Healthy(c.Request.Context()); err != nil {
		c.JSON(http.StatusServiceUnavailable, domain.Response{Msg: err.Error()})
		return
	}

	c.JSON(http.StatusOK, domain.Response{Msg: "healthy"})
}

func (d *RestDelivery) DeleteUser(c *gin.Context) {}

func (d *RestDelivery) UserGetPermissions(c *gin.Context) {}

func (d *RestDelivery) UserGetRoles(c *gin.Context)           {}
func (d *RestDelivery) UserCheck(c *gin.Context)              {}
func (d *RestDelivery) UserAddPermission(c *gin.Context)      {}
func (d *RestDelivery) UserRemovePermission(c *gin.Context)   {}
func (d *RestDelivery) UserAddRole(c *gin.Context)            {}
func (d *RestDelivery) UserRemoveRole(c *gin.Context)         {}
func (d *RestDelivery) DeleteRole(c *gin.Context)             {}
func (d *RestDelivery) RoleGetUsers(c *gin.Context)           {}
func (d *RestDelivery) RoleGetPermissions(c *gin.Context)     {}
func (d *RestDelivery) RoleAddPermission(c *gin.Context)      {}
func (d *RestDelivery) RoleRemovePermission(c *gin.Context)   {}
func (d *RestDelivery) RoleInheritRole(c *gin.Context)        {}
func (d *RestDelivery) RoleUnInheritRole(c *gin.Context)      {}
func (d *RestDelivery) RoleGetChildRole(c *gin.Context)       {}
func (d *RestDelivery) RoleGetParentRole(c *gin.Context)      {}
func (d *RestDelivery) DeleteObject(c *gin.Context)           {}
func (d *RestDelivery) WhichRoleHasPermission(c *gin.Context) {}
func (d *RestDelivery) WhichUserHasPermission(c *gin.Context) {}
