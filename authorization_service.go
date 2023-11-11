package yauthorization

import (
	"gorm.io/gorm"
)

type AuthorizeService struct {
	db *gorm.DB
}

func NewAuthorizeService(
	db *gorm.DB,
) *AuthorizeService {
	return &AuthorizeService{
		db: db,
	}
}

func (auth *AuthorizeService) RoleCreate(identity Identity, role *RoleIdentity) error {
	return auth.db.Transaction(func(tx *gorm.DB) error {
		secq := NewSecQuery(identity, tx)
		return secq.Save(role).Error

	})
}

type PermissionPrepload struct{}

// GetEntityID implements RawQuery.
func (*PermissionPrepload) GetEntityID() string {
	panic("unimplemented")
}

// GetDomainID implements RawQuery.
func (*PermissionPrepload) GetDomainID() uint {
	return 0
}

// Permission implements RawQuery.
func (pre *PermissionPrepload) Permission(identity Identity, action Action) *EntityPermission {
	return &EntityPermission{
		IdentityID: identity.IdentityID(),
		DomainID:   pre.GetDomainID(),
		EntityID:   "EntityPermission",
		Policy:     Deny,
		Action:     action,
	}
}

// Raw implements RawQuery.
func (*PermissionPrepload) Raw() string {
	return "Permissions"
}

type RoleListQuery struct {
	DomainID uint `json:"domain_id" form:"domain_id" schema:"domain_id"`
}

func (auth *AuthorizeService) RoleList(identity Identity, query *RoleListQuery) (RoleListResult, error) {
	hasil := RoleListResult{}
	err := NewSecQuery(identity, auth.db).
		Model(&RoleIdentity{
			DomainID: query.DomainID,
		}).
		Preload(&PermissionPrepload{}).
		Find(&hasil).Error

	return hasil, err
}
func (auth *AuthorizeService) RoleDelete(identity Identity, role *RoleIdentity) error {
	err := NewSecQuery(identity, auth.db).Delete(role).Error
	return err
}

func (auth *AuthorizeService) RoleUpdatePermission(identity Identity, role *RoleIdentity, inheritRole *RoleIdentity, perms []*EntityPermission) error {

	return auth.db.Transaction(func(tx *gorm.DB) error {
		err := NewSecQuery(identity, auth.db).Model(&EntityPermission{
			DomainID:   role.DomainID,
			IdentityID: role.ID,
		}).Delete(&EntityPermission{}).Error

		if err != nil {
			return err
		}

		if inheritRole != nil {
			inheritperms := EntityPermissionList{}

			err := NewSecQuery(identity, tx).Model(&EntityPermission{
				DomainID:   inheritRole.DomainID,
				IdentityID: inheritRole.ID,
			}).Find(&inheritperms).Error

			if err != nil {
				return err
			}

			for _, perm := range inheritperms {
				perms = append(perms, perm)
			}
		}

		for _, perm := range perms {
			err := NewSecQuery(identity, tx).Save(perm).Error
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (auth *AuthorizeService) createAdminRole(updater Identity, domainID uint) (*RoleIdentity, error) {
	role := RoleIdentity{
		Key:      "admin",
		DomainID: domainID,
	}

	err := auth.RoleCreate(updater, &role)
	if err != nil {
		return &role, err
	}

	ent := &RoleIdentity{
		DomainID: domainID,
	}
	permsent := &EntityPermission{
		DomainID: domainID,
	}

	perms := []*EntityPermission{
		ent.Permission(updater, Create),
		ent.Permission(updater, Read),
		ent.Permission(updater, Delete),
		ent.Permission(updater, Update),

		permsent.Permission(updater, Create),
		permsent.Permission(updater, Read),
		permsent.Permission(updater, Delete),
		permsent.Permission(updater, Update),
	}

	for _, perm := range perms {
		perm.Policy = Allow
	}

	err = auth.RoleUpdatePermission(updater, &role, nil, perms)
	if err != nil {
		return &role, err
	}

	return &role, nil
}

// ----------------------- user ------------------

func (auth *AuthorizeService) UserAddAdminDomain(updater Identity, user Identity, domainID uint) error {
	role, err := auth.createAdminRole(updater, domainID)

	if err != nil {
		return err
	}

	err = user.SetRole(auth.db, role)
	if err != nil {
		return err
	}

	return nil

}
func (auth *AuthorizeService) UserRemoveAdminDomain(updater Identity, user Identity, roleID uint) error {
	err := user.DeleteRole(auth.db, roleID)
	if err != nil {
		return err
	}

	return nil
}
func (auth *AuthorizeService) UserListRole()   { panic("not implemented") }
func (auth *AuthorizeService) UserAddRole()    { panic("not implemented") }
func (auth *AuthorizeService) UserDeleteRole() { panic("not implemented") }

// ----------------------- entity ------------------

func (auth *AuthorizeService) ListEntity() ([]*EntityInfo, error) {
	hasil := []*EntityInfo{}

	err := auth.db.Model(&EntityInfo{}).Find(&hasil).Error

	return hasil, err
}

// ----------------------- domain ------------------

func (auth *AuthorizeService) ListDomain()     { panic("not implemented") }
func (auth *AuthorizeService) ListDomainRole() { panic("not implemented") }

// -------------------- model ----------------------------

type RoleListResult []*RoleIdentity

// GetEntityID implements Entity.
func (*RoleListResult) GetEntityID() string {
	return "RoleIdentity"
}

// GetDomainID implements Entity.
func (*RoleListResult) GetDomainID() uint {
	return 0
}

// Permission implements Entity.
func (list *RoleListResult) Permission(identity Identity, action Action) *EntityPermission {
	return &EntityPermission{
		IdentityID: identity.IdentityID(),
		DomainID:   list.GetDomainID(),
		EntityID:   list.GetEntityID(),
		Policy:     Deny,
		Action:     action,
	}
}

type EntityPermissionList []*EntityPermission

// GetEntityID implements Entity.
func (*EntityPermissionList) GetEntityID() string {
	return "EntityPermission"
}

// GetDomainID implements Entity.
func (*EntityPermissionList) GetDomainID() uint {
	return 0
}

// Permission implements Entity.
func (list *EntityPermissionList) Permission(identity Identity, action Action) *EntityPermission {
	return &EntityPermission{
		IdentityID: identity.IdentityID(),
		DomainID:   list.GetDomainID(),
		EntityID:   list.GetEntityID(),
		Policy:     Deny,
		Action:     action,
	}
}
