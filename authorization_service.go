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

func (auth *AuthorizeService) RoleList(identity Identity) (RoleListResult, error) {
	hasil := RoleListResult{}
	err := NewSecQuery(identity, auth.db).Model(&RoleIdentity{}).Find(&hasil).Error

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

// ----------------------- user ------------------

func (auth *AuthorizeService) UserAddAdminDomain(identity Identity, domainID uint) error {
	panic("not implemented")
}
func (auth *AuthorizeService) UserRemoveAdminDomain(identity Identity, domainID uint) error {
	panic("not implemented")
}
func (auth *AuthorizeService) UserListRole()   { panic("not implemented") }
func (auth *AuthorizeService) UserAddRole()    { panic("not implemented") }
func (auth *AuthorizeService) UserDeleteRole() { panic("not implemented") }

// ----------------------- entity ------------------

func (auth *AuthorizeService) ListEntity() { panic("not implemented") }

// ----------------------- domain ------------------

func (auth *AuthorizeService) ListDomain()     { panic("not implemented") }
func (auth *AuthorizeService) ListDomainRole() { panic("not implemented") }

// -------------------- model ----------------------------

type RoleListResult []*RoleIdentity

// GetDomainID implements Entity.
func (*RoleListResult) GetDomainID() uint {
	return 0
}

// Permission implements Entity.
func (list *RoleListResult) Permission(identity Identity, action Action) *EntityPermission {
	return &EntityPermission{
		IdentityID: identity.IdentityID(),
		DomainID:   list.GetDomainID(),
		EntityID:   "RoleIdentity",
		Policy:     Deny,
		Action:     action,
	}
}

type EntityPermissionList []*EntityPermission

// GetDomainID implements Entity.
func (*EntityPermissionList) GetDomainID() uint {
	return 0
}

// Permission implements Entity.
func (list *EntityPermissionList) Permission(identity Identity, action Action) *EntityPermission {
	return &EntityPermission{
		IdentityID: identity.IdentityID(),
		DomainID:   list.GetDomainID(),
		EntityID:   "EntityPermission",
		Policy:     Deny,
		Action:     action,
	}
}
