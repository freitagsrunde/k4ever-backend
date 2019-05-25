package utils

import (
	"fmt"
	"strings"

	"github.com/freitagsrunde/k4ever-backend/internal/k4ever"
	ldap "gopkg.in/ldap.v3"
)

// Connect to ldap and return connection object
func LdapConnect(config k4ever.Config) (*ldap.Conn, error) {
	conn, err := ldap.DialURL(config.LdapHost())

	if err != nil {
		return nil, fmt.Errorf("Failed to connect to ldap server: %s", config.LdapHost())
	}

	if err := conn.Bind(config.LdapBind(), config.LdapPassword()); err != nil {
		fmt.Println(err.Error())
		return nil, fmt.Errorf("Failed to bind to ldap server: %s", config.LdapBind())
	}
	return conn, nil
}

// try to authenticate user against ldap
func LdapAuth(user string, password string, conn *ldap.Conn, config k4ever.Config) error {
	result, err := conn.Search(ldap.NewSearchRequest(
		config.LdapBaseDN(),
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		filter(user, config),
		[]string{"dn"},
		nil,
	))

	if err != nil {
		return fmt.Errorf("Failed to find user: %s", user)
	}

	if len(result.Entries) < 1 {
		return fmt.Errorf("User does not exist: %s", user)
	}

	if len(result.Entries) > 1 {
		return fmt.Errorf("Too many entries returned")
	}

	if err := conn.Bind(result.Entries[0].DN, password); err != nil {
		fmt.Errorf("Failed to auth. %s", err)
	}
	return nil
}

func filter(needle string, config k4ever.Config) string {
	res := strings.Replace(
		config.LdapFilterDN(),
		"{username}",
		needle,
		-1,
	)
	return res
}
