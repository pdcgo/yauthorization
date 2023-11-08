package yauthorization

import "time"

type Policy int

const (
	Allow Policy = 1
	Deny  Policy = 0
)

type Action string

const (
	Create Action = "create"
	Update Action = "update"
	Read   Action = "read"
	Delete Action = "delete"
)

type EntityPermission struct {
	IdentityID uint   `gorm:"primaryKey" json:"identity_id"`
	DomainID   uint   `gorm:"primaryKey" json:"domain_id"`
	EntityID   string `gorm:"primaryKey" json:"entity_id"`
	Action     Action `gorm:"primaryKey" json:"action"`
	Policy     Policy `json:"policy"`
}

// GetDomainID implements Entity.
func (perm *EntityPermission) GetDomainID() uint {
	return perm.DomainID
}

// Permission implements Entity.
func (perm *EntityPermission) Permission(identity Identity, action Action) *EntityPermission {
	return &EntityPermission{
		IdentityID: identity.IdentityID(),
		DomainID:   perm.GetDomainID(),
		EntityID:   "EntityPermission",
		Policy:     Deny,
		Action:     action,
	}
}

type RoleIdentity struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Key       string    `json:"key" gorm:"index:domain_key,unique"`
	DomainID  uint      `json:"domain_id" gorm:"index:domain_key,unique"`
	CreatedAt time.Time `json:"create_at"`
	UpdatedAt time.Time `json:"update_at"`
}

// IdentityID implements Identity.
func (r *RoleIdentity) IdentityID() uint {
	return r.ID
}

// IsSuperUser implements Identity.
func (RoleIdentity) IsSuperUser() bool {
	return false
}

// GetDomainID implements Entity.
func (role *RoleIdentity) GetDomainID() uint {
	return role.DomainID
}

// Permission implements Entity.
func (role *RoleIdentity) Permission(identity Identity, action Action) *EntityPermission {
	return &EntityPermission{
		IdentityID: identity.IdentityID(),
		DomainID:   role.GetDomainID(),
		EntityID:   "RoleIdentity",
		Policy:     Deny,
		Action:     action,
	}
}
