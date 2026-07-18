package auth

import (
	"fmt"
	"strings"
)

type Scope string

const (
	ScopeServerRead  Scope = "server:read"
	ScopeServerWrite Scope = "server:write"
	ScopeBackupRead  Scope = "backup:read"
	ScopeBackupWrite Scope = "backup:write"
	ScopeAdmin       Scope = "admin"
)

var AllScopes = Scopes{ScopeServerRead, ScopeServerWrite, ScopeBackupRead, ScopeBackupWrite, ScopeAdmin}

var validScopes = map[Scope]bool{
	ScopeServerRead:  true,
	ScopeServerWrite: true,
	ScopeBackupRead:  true,
	ScopeBackupWrite: true,
	ScopeAdmin:       true,
}

type Scopes []Scope

func (s Scopes) Contains(scope Scope) bool {
	if scope == ScopeAdmin {
		for _, ss := range s {
			if ss == ScopeAdmin {
				return true
			}
		}
		return false
	}
	for _, ss := range s {
		if ss == ScopeAdmin || ss == scope {
			return true
		}
	}
	return false
}

func (s Scopes) HasAll(required ...Scope) bool {
	for _, r := range required {
		if !s.Contains(r) {
			return false
		}
	}
	return true
}

func (s Scopes) HasAny(required ...Scope) bool {
	for _, r := range required {
		if s.Contains(r) {
			return true
		}
	}
	return false
}

func (s Scopes) Intersect(other Scopes) Scopes {
	var result Scopes
	for _, scope := range s {
		if other.Contains(scope) {
			result = append(result, scope)
		}
	}
	return result
}

func (s Scopes) String() string {
	ss := make([]string, len(s))
	for i, scope := range s {
		ss[i] = string(scope)
	}
	return strings.Join(ss, ", ")
}

func (s Scopes) Validate() error {
	for _, scope := range s {
		if !validScopes[scope] {
			return fmt.Errorf("invalid scope: %s", scope)
		}
	}
	return nil
}
