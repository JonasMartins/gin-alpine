// Package role ...
package role

import (
	"strings"

	base "gin-alpine/src/pkg/models"
)

type PositionType uint8

const (
	CUSTOMER PositionType = iota + 1
	MANAGER
	ADMIN
	DEV
)

func (p PositionType) String() string {
	return [...]string{"", "CUSTOMER", "MANAGER", "ADMIN", "DEV"}[p]
}

func IsValidPositionType(p PositionType) bool {
	return p < 4
}

func ParsePositionType(p string) (PositionType, bool) {
	switch strings.ToLower(p) {
	case "customer":
		return CUSTOMER, true
	case "admin":
		return ADMIN, true
	case "manager":
		return MANAGER, true
	case "dev":
		return DEV, true
	}
	return 0, false
}

func GetPositionTypeFromString(s string) *PositionType {
	p, ok := ParsePositionType(s)
	if !ok {
		return nil
	}
	return &p
}

type Role struct {
	Base        *base.Base    `json:"base,omitempty"`
	Role        *PositionType `json:"role,omitempty"`
	Description *string       `json:"description,omitempty"`
}
