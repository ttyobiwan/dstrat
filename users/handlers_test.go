package users

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ttyobiwan/dstrat/internal/tests"
)

func TestCreateUser(t *testing.T) {
	db := tests.GetTestDB(t)
	h := NewUserHandler()

	data := `{
		"username": "chito"
	}`

	rec, c := tests.MakeRequest(http.MethodPost, "/api/users", "", data)
	if assert.NoError(t, h.CreateUser(&UserContext{c, db})) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		user := User{}
		err := json.Unmarshal(rec.Body.Bytes(), &user)
		assert.NoError(t, err)
		assert.NotEmpty(t, user.ID)
	}
}

func TestCreateNonUniqueUser(t *testing.T) {
	db := tests.GetTestDB(t)
	h := NewUserHandler()

	data := `{
		"username": "chito"
	}`

	rec, c := tests.MakeRequest(http.MethodPost, "/api/users", "", data)
	if assert.NoError(t, h.CreateUser(&UserContext{c, db})) {
		assert.Equal(t, http.StatusCreated, rec.Code)
	}

	rec, c = tests.MakeRequest(http.MethodPost, "/api/users", "", data)
	if assert.NoError(t, h.CreateUser(&UserContext{c, db})) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		resp := struct{ Detail string }{}
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, resp.Detail, "Username already taken")
	}
}
