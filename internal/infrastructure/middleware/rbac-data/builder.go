package rbac_data

// RBACMapBuilder represent the builder for the rbac mapping
type RBACMapBuilder interface {
	Build() RBACMap
	Add(r Route, m Method, p RBACPermission) RBACMapBuilder
}

type rbacMapBuilder RBACMap

// NewRBACMap create a builder to generate a RBACMap
func NewRBACMapBuilder() RBACMapBuilder {
	return &rbacMapBuilder{}
}

func (d *rbacMapBuilder) Build() RBACMap {
	return RBACMap(*d)
}

func (d *rbacMapBuilder) Add(route Route, method Method, permission RBACPermission) RBACMapBuilder {
	methods, ok := (*d)[route]
	if !ok {
		(*d)[route] = map[Method]RBACPermission{method: permission}
		return d
	}
	methods[method] = permission
	return d
}
