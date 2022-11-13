package packages

import (
	"context"
	"fmt"

	"github.com/Nerzal/gocloak/v12"
)

type keycloackConfig struct {
	GoCloak           *gocloak.GoCloak
	GoCloakToken      *gocloak.JWT
	GocloakErr        error
	keycloackHost     string
	keycloackUser     string
	keycloackPassword string
	keycloackRealm    string
}

func NewKeycloak(realm string) *keycloackConfig {
	config := keycloackConfig{}
	config.keycloackHost = GetString("KC_HOST")
	config.keycloackUser = GetString("KC_USER")
	config.keycloackPassword = GetString("KC_PASSWORD")

	gocloakRes := gocloak.NewClient(config.keycloackHost)
	gocloakToken, gocloakTokenErr := gocloakRes.LoginAdmin(context.Background(), config.keycloackUser, config.keycloackPassword, realm)

	return &keycloackConfig{GoCloak: gocloakRes, GoCloakToken: gocloakToken, GocloakErr: gocloakTokenErr, keycloackRealm: realm}
}

// CreateUser creates the given user in the given realm and returns it's userID Note: Keycloak has not documented what members of the User object are actually being accepted, when creating a user
func (h *keycloackConfig) Create(ctx context.Context, body gocloak.User) (string, error) {
	res, err := h.GoCloak.CreateUser(ctx, h.GoCloakToken.AccessToken, "realm", body)
	if err != nil {
		return res, err
	}

	return res, nil
}

// AddClientRoleToUser adds client-level role mappings
func (h *keycloackConfig) AddClientRoleToUser(ctx context.Context, idOfClient, userId string, roles []gocloak.Role) (interface{}, error) {
	err := h.GoCloak.AddClientRoleToUser(ctx, h.GoCloakToken.AccessToken, h.keycloackRealm, idOfClient, userId, roles)
	if err != nil {
		return nil, err
	}

	return fmt.Sprintf("Assign clientRoleUser to userId %s success", userId), nil
}

// AddClientRolesToGroup adds a client role to the group
func (h *keycloackConfig) AddClientRolesToGroup(ctx context.Context, idOfClient, groupId string, roles []gocloak.Role) (interface{}, error) {
	err := h.GoCloak.AddClientRolesToGroup(ctx, h.GoCloakToken.AccessToken, h.keycloackRealm, idOfClient, groupId, roles)
	if err != nil {
		return nil, err
	}

	return fmt.Sprintf("Assign clientRolesGroup to groupId %s success", groupId), nil
}

// AddRealmRoleToUser adds realm-level role mappings
func (h *keycloackConfig) AddRealmRoleToUser(ctx context.Context, userId string, roles []gocloak.Role) (interface{}, error) {
	err := h.GoCloak.AddRealmRoleToUser(ctx, h.GoCloakToken.AccessToken, h.keycloackRealm, userId, roles)
	if err != nil {
		return nil, err
	}

	return fmt.Sprintf("Assign realmRoleUser to userId %s success", userId), nil
}

// AddRealmRoleToGroup adds realm-level role mappings
func (h *keycloackConfig) AddRealmRoleToGroup(ctx context.Context, groupId string, roles []gocloak.Role) (interface{}, error) {
	err := h.GoCloak.AddRealmRoleToGroup(ctx, h.GoCloakToken.AccessToken, h.keycloackRealm, groupId, roles)
	if err != nil {
		return nil, err
	}

	return fmt.Sprintf("Assign realmRoleGroup to groupId %s success", groupId), nil
}
