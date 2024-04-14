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

func (d *RestDelivery) DeleteUser(c *gin.Context) {
	err := d.usecase.DeleteUser(c.Request.Context(), c.Param("name"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.Response{Msg: err.Error()})
		return
	}
}

func (d *RestDelivery) UserGetPermissions(c *gin.Context) {
	pers, err := d.usecase.UserGetPermissions(c.Request.Context(), c.Param("name"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.Response{Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, pers)
}

func (d *RestDelivery) UserGetRoles(c *gin.Context) {
	roles, err := d.usecase.UserGetRoles(c.Request.Context(), c.Param("name"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.Response{Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, roles)
}

func (d *RestDelivery) UserCheck(c *gin.Context) {
	ok, err := d.usecase.UserCheck(c, c.Param("name"), c.Param("objns"),
		c.Param("rel"), c.Param("objname"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.Response{Msg: err.Error()})
		return
	}
	if !ok {
		c.Status(http.StatusForbidden)
		return
	}
}

func (d *RestDelivery) UserAddPermission(c *gin.Context) {
	var requestBody struct {
		Relation string `json:"relation"`
		ObjNs    string `json:"obj_ns"`
		ObjName  string `json:"obj_name"`
	}
	if err := c.ShouldBindUri(&requestBody); err != nil {
		c.JSON(400, gin.H{"msg": err})
		return
	}
	if err := d.usecase.UserAddPermission(c.Request.Context(), c.Param("name"),
		domain.Permission{
			Rel:  requestBody.Relation,
			Ns:   requestBody.ObjNs,
			Name: requestBody.ObjName,
		}); err != nil {
		c.JSON(http.StatusInternalServerError, domain.Response{Msg: err.Error()})
		return
	}
}
func (d *RestDelivery) UserRemovePermission(c *gin.Context) {
	var requestBody struct {
		Relation string `json:"relation"`
		ObjNs    string `json:"obj_ns"`
		ObjName  string `json:"obj_name"`
	}
	if err := c.ShouldBindUri(&requestBody); err != nil {
		c.JSON(400, gin.H{"msg": err})
		return
	}
	if err := d.usecase.UserRemovePermission(c.Request.Context(), c.Param("name"),
		domain.Permission{
			Rel:  requestBody.Relation,
			Ns:   requestBody.ObjNs,
			Name: requestBody.ObjName,
		}); err != nil {
		c.JSON(http.StatusInternalServerError, domain.Response{Msg: err.Error()})
		return
	}
}

func (d *RestDelivery) UserAddRole(c *gin.Context) {
	var requestBody struct {
		RoleName string `json:"role_name"`
	}
	if err := c.ShouldBindUri(&requestBody); err != nil {
		c.JSON(400, gin.H{"msg": err})
		return
	}
	if err := d.usecase.UserAddRole(c.Request.Context(), c.Param("name"),
		requestBody.RoleName); err != nil {
		c.JSON(http.StatusInternalServerError, domain.Response{Msg: err.Error()})
		return
	}
}

func (d *RestDelivery) UserRemoveRole(c *gin.Context) {
	var requestBody struct {
		RoleName string `json:"role_name"`
	}
	if err := c.ShouldBindUri(&requestBody); err != nil {
		c.JSON(400, gin.H{"msg": err})
		return
	}
	if err := d.usecase.UserRemoveRole(c.Request.Context(), c.Param("name"),
		requestBody.RoleName); err != nil {
		c.JSON(http.StatusInternalServerError, domain.Response{Msg: err.Error()})
		return
	}
}

func (d *RestDelivery) DeleteRole(c *gin.Context) {
	if err := d.usecase.DeleteRole(c.Request.Context(), c.Param("name")); err != nil {
		c.JSON(http.StatusInternalServerError, domain.Response{Msg: err.Error()})
		return
	}
}

func (d *RestDelivery) RoleGetUsers(c *gin.Context) {
	users, err := d.usecase.RoleGetUsers(c.Request.Context(), c.Param("name"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.Response{Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (d *RestDelivery) RoleGetPermissions(c *gin.Context) {
	pers, err := d.usecase.RoleGetPermissions(c.Request.Context(), c.Param("name"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.Response{Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, pers)
}

func (d *RestDelivery) RoleAddPermission(c *gin.Context) {
	var requestBody struct {
		Relation string `json:"relation"`
		ObjNs    string `json:"obj_ns"`
		ObjName  string `json:"obj_name"`
	}
	if err := c.ShouldBindUri(&requestBody); err != nil {
		c.JSON(400, gin.H{"msg": err})
		return
	}
	if err := d.usecase.RoleAddPermission(c.Request.Context(), c.Param("name"),
		domain.Permission{
			Rel:  requestBody.Relation,
			Ns:   requestBody.ObjNs,
			Name: requestBody.ObjName,
		}); err != nil {
		c.JSON(http.StatusInternalServerError, domain.Response{Msg: err.Error()})
		return
	}
}

func (d *RestDelivery) RoleRemovePermission(c *gin.Context) {
	var requestBody struct {
		Relation string `json:"relation"`
		ObjNs    string `json:"obj_ns"`
		ObjName  string `json:"obj_name"`
	}
	if err := c.ShouldBindUri(&requestBody); err != nil {
		c.JSON(400, gin.H{"msg": err})
		return
	}
	if err := d.usecase.RoleRemovePermission(c.Request.Context(), c.Param("name"),
		domain.Permission{
			Rel:  requestBody.Relation,
			Ns:   requestBody.ObjNs,
			Name: requestBody.ObjName,
		}); err != nil {
		c.JSON(http.StatusInternalServerError, domain.Response{Msg: err.Error()})
		return
	}
}

func (d *RestDelivery) RoleInheritRole(c *gin.Context) {
	var requestBody struct {
		Name string `json:"name"`
	}
	if err := c.ShouldBindUri(&requestBody); err != nil {
		c.JSON(400, gin.H{"msg": err})
		return
	}
	if err := d.usecase.RoleInheritRole(c.Request.Context(), c.Param("name"),
		requestBody.Name); err != nil {
		c.JSON(http.StatusInternalServerError, domain.Response{Msg: err.Error()})
		return
	}
}

func (d *RestDelivery) RoleUnInheritRole(c *gin.Context) {
	var requestBody struct {
		Name string `json:"name"`
	}
	if err := c.ShouldBindUri(&requestBody); err != nil {
		c.JSON(400, gin.H{"msg": err})
		return
	}
	if err := d.usecase.RoleUnInheritRole(c.Request.Context(), c.Param("name"),
		requestBody.Name); err != nil {
		c.JSON(http.StatusInternalServerError, domain.Response{Msg: err.Error()})
		return
	}
}

func (d *RestDelivery) RoleGetChildRole(c *gin.Context) {
	roles, err := d.usecase.RoleGetChildRole(c.Request.Context(), c.Param("name"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.Response{Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, roles)
}

func (d *RestDelivery) RoleGetParentRole(c *gin.Context) {
	roles, err := d.usecase.RoleGetParentRole(c.Request.Context(), c.Param("name"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.Response{Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, roles)
}

func (d *RestDelivery) DeleteObject(c *gin.Context) {
	if err := d.usecase.DeleteObject(c.Request.Context(), c.Param("ns"),
		c.Param("name")); err != nil {
		c.JSON(http.StatusInternalServerError, domain.Response{Msg: err.Error()})
		return
	}
}

func (d *RestDelivery) WhichRoleHasPermission(c *gin.Context) {
	roles, err := d.usecase.WhichRoleHasPermission(c.Request.Context(), c.Param("ns"),
		c.Param("name"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.Response{Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, roles)
}

func (d *RestDelivery) WhichUserHasPermission(c *gin.Context) {
	users, err := d.usecase.WhichUserHasPermission(c.Request.Context(), c.Param("ns"),
		c.Param("name"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.Response{Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}
