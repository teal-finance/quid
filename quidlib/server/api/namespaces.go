package api

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/teal-finance/quid/quidlib/server"
	db "github.com/teal-finance/quid/quidlib/server/db"
	"github.com/teal-finance/quid/quidlib/tokens"
)

// AllNamespaces : get all namespaces
func AllNamespaces(c echo.Context) error {
	data, err := db.SelectAllNamespaces()
	if err != nil {
		log.Fatal(err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "error selecting namespaces",
		})
	}
	return c.JSON(http.StatusOK, &data)
}

// SetNamespaceRefreshTokenMaxTTL : set a max refresh token ttl for a namespace
func SetNamespaceRefreshTokenMaxTTL(c echo.Context) error {
	m := echo.Map{}
	if err := c.Bind(&m); err != nil {
		return err
	}
	ID := int64(m["id"].(float64))
	refreshMxTTL := m["refresh_max_ttl"].(string)

	err := db.UpdateNamespaceRefreshTokenMaxTTL(ID, refreshMxTTL)
	if err != nil {
		log.Fatal(err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "error updating tokens max ttl in namespace",
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"ok": true,
	})
}

// SetNamespaceTokenMaxTTL : set a max access token ttl for a namespace
func SetNamespaceTokenMaxTTL(c echo.Context) error {
	m := echo.Map{}
	if err := c.Bind(&m); err != nil {
		return err
	}
	ID := int64(m["id"].(float64))
	maxTTL := m["max_ttl"].(string)

	err := db.UpdateNamespaceTokenMaxTTL(ID, maxTTL)
	if err != nil {
		log.Fatal(err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "error updating tokens max ttl in namespace",
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"ok": true,
	})
}

// NamespaceInfo : info about a namespace
func NamespaceInfo(c echo.Context) error {
	m := echo.Map{}
	if err := c.Bind(&m); err != nil {
		return err
	}
	ID := int64(m["id"].(float64))

	nu, err := db.CountUsersForNamespace(ID)
	if err != nil {
		log.Fatal(err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "error counting users in namespace",
		})
	}

	g, err := db.SelectGroupsForNamespace(ID)
	if err != nil {
		log.Fatal(err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "error counting users in namespace",
		})
	}

	data := server.NamespaceInfo{
		NumUsers: nu,
		Groups:   g,
	}

	return c.JSON(http.StatusOK, &data)
}

// GetNamespaceKey : get the key for a namespace
func GetNamespaceKey(c echo.Context) error {
	m := echo.Map{}
	if err := c.Bind(&m); err != nil {
		return err
	}
	ID := int64(m["id"].(float64))

	found, data, err := db.SelectNamespaceKey(ID)
	if err != nil {
		log.Fatal(err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "error finding namespace key",
		})
	}
	if !found {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "namespace not found",
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"key": data,
	})

}

// FindNamespace : namespace creation http handler
func FindNamespace(c echo.Context) error {
	m := echo.Map{}
	if err := c.Bind(&m); err != nil {
		return err
	}
	name := m["name"].(string)

	data, err := db.SelectNamespaceStartsWith(name)
	if err != nil {
		log.Fatal(err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "error finding namespace",
		})
	}
	return c.JSON(http.StatusOK, &data)

}

// DeleteNamespace : namespace creation http handler
func DeleteNamespace(c echo.Context) error {
	m := echo.Map{}
	if err := c.Bind(&m); err != nil {
		return err
	}
	ID := int64(m["id"].(float64))

	qRes := db.DeleteNamespace(ID)
	if qRes.HasError {
		emo.QueryError(qRes.Error.Message)
		if qRes.Error.HasUserMessage {
			return c.JSON(http.StatusConflict, echo.Map{
				"error": "error deleting namespace: " + qRes.Error.Message,
			})
		}
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "ok",
	})

}

// SetNamespaceEndpointAvailability :
func SetNamespaceEndpointAvailability(c echo.Context) error {
	m := echo.Map{}
	if err := c.Bind(&m); err != nil {
		return err
	}
	ID := int64(m["id"].(float64))
	enable := m["enable"].(bool)

	err := db.SetNamespaceEndpointAvailability(ID, enable)
	if err != nil {
		return c.JSON(http.StatusConflict, echo.Map{
			"error": "error updating namespace",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "ok",
	})
}

// CreateNamespace : namespace creation http handler
func CreateNamespace(c echo.Context) error {
	m := echo.Map{}
	if err := c.Bind(&m); err != nil {
		return err
	}
	name := m["name"].(string)
	maxTTL := m["max_ttl"].(string)
	refreshMaxTTL := m["refresh_max_ttl"].(string)
	enableEndpoint := m["enable_endpoint"].(bool)

	key := tokens.GenKey()
	refreshKey := tokens.GenKey()
	ns, exists, err := createNamespace(name, key, refreshKey, maxTTL, refreshMaxTTL, enableEndpoint)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "error creating namespace",
		})
	}
	if exists {
		return c.JSON(http.StatusConflict, echo.Map{
			"error": "namespace already exists",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"namespace_id": ns.ID,
	})
}

// createNamespace : create a namespace
func createNamespace(name, key, refreshKey, ttl, refreshMaxTTL string, endpoint bool) (server.Namespace, bool, error) {
	ns := server.Namespace{}

	exists, err := db.NamespaceExists(name)
	if err != nil {
		return ns, false, err
	}
	if exists {
		return ns, true, nil
	}

	uid, err := db.CreateNamespace(name, key, refreshKey, ttl, refreshMaxTTL, endpoint)
	if err != nil {
		return ns, false, err
	}
	ns.ID = uid
	ns.Name = name
	return ns, false, nil
}
