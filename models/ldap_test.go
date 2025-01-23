package models

import (
	"github.com/jinzhu/gorm"
	"gopkg.in/check.v1"
)

func (s *ModelsSuite) TestPostLDAP(c *check.C) {
	ldap := LDAP{
		UserId:     1,
		Name:       "Default",
		Protocol:   "ldap",
		Host:       "1.1.1.1:389",
		BaseDn:     "ou=people,dc=example,dc=com",
		Query:      "(objectClass=inetOrgPerson)",
		Attributes: "givenName,cn,email,description",
	}
	err := PostLDAP(&ldap)
	c.Assert(err, check.Equals, nil)
	ls, err := GetLDAPs(1)
	c.Assert(err, check.Equals, nil)
	c.Assert(len(ls), check.Equals, 1)
}

func (s *ModelsSuite) TestPostLDAPNoProtocol(c *check.C) {
	ldap := LDAP{
		UserId:     1,
		Name:       "Default",
		Host:       "1.1.1.1:389",
		BaseDn:     "ou=people,dc=example,dc=com",
		Query:      "(objectClass=inetOrgPerson)",
		Attributes: "givenName,cn,email,description",
	}
	err := PostLDAP(&ldap)
	c.Assert(err, check.Equals, ErrLdapProtocolNotSpecified)
}

func (s *ModelsSuite) TestPostLDAPInvalidProtocol(c *check.C) {
	ldap := LDAP{
		UserId:     1,
		Name:       "Default",
		Protocol:   "http",
		Host:       "1.1.1.1:389",
		BaseDn:     "ou=people,dc=example,dc=com",
		Query:      "(objectClass=inetOrgPerson)",
		Attributes: "givenName,cn,email,description",
	}
	err := PostLDAP(&ldap)
	c.Assert(err, check.Equals, ErrLdapProtocolInvalid)
}

func (s *ModelsSuite) TestPostLDAPNoHost(c *check.C) {
	ldap := LDAP{
		UserId:     1,
		Name:       "Default",
		Protocol:   "ldap",
		BaseDn:     "ou=people,dc=example,dc=com",
		Query:      "(objectClass=inetOrgPerson)",
		Attributes: "givenName,cn,email,description",
	}
	err := PostLDAP(&ldap)
	c.Assert(err, check.Equals, ErrLdapHostNotSpecified)
}

func (s *ModelsSuite) TestPostLDAPNoBaseDn(c *check.C) {
	ldap := LDAP{
		UserId:     1,
		Name:       "Default",
		Protocol:   "ldap",
		Host:       "1.1.1.1:389",
		Query:      "(objectClass=inetOrgPerson)",
		Attributes: "givenName,cn,email,description",
	}
	err := PostLDAP(&ldap)
	c.Assert(err, check.Equals, ErrLdapBaseDnNotSpecified)
}

func (s *ModelsSuite) TestPostLDAPNoQuery(c *check.C) {
	ldap := LDAP{
		UserId:     1,
		Name:       "Default",
		Protocol:   "ldap",
		Host:       "1.1.1.1:389",
		BaseDn:     "ou=people,dc=example,dc=com",
		Attributes: "givenName,cn,email,description",
	}
	err := PostLDAP(&ldap)
	c.Assert(err, check.Equals, ErrLdapQueryNotSpecified)
}

func (s *ModelsSuite) TestPostLDAPNoAttributes(c *check.C) {
	ldap := LDAP{
		UserId:   1,
		Name:     "Default",
		Protocol: "ldap",
		Host:     "1.1.1.1:389",
		BaseDn:   "ou=people,dc=example,dc=com",
		Query:    "(objectClass=inetOrgPerson)",
	}
	err := PostLDAP(&ldap)
	c.Assert(err, check.Equals, ErrLdapAttributesNotSpecified)
}

func (s *ModelsSuite) TestInvalidAttributes(c *check.C) {
	ldap := LDAP{
		UserId:     1,
		Name:       "Default",
		Protocol:   "ldap",
		Host:       "1.1.1.1:389",
		BaseDn:     "ou=people,dc=example,dc=com",
		Query:      "(objectClass=inetOrgPerson)",
		Attributes: "invalid;attributes",
	}

	valid, attributes := ldap.ValidateAttributes()
	c.Assert(valid, check.Equals, false)
	c.Assert(attributes, check.Equals, ldap.Attributes)
}

func (s *ModelsSuite) TestValidAttributesAreTrimmed(c *check.C) {
	ldap := LDAP{
		UserId:     1,
		Name:       "Default",
		Protocol:   "ldap",
		Host:       "1.1.1.1:389",
		BaseDn:     "ou=people,dc=example,dc=com",
		Query:      "(objectClass=inetOrgPerson)",
		Attributes: " givenName ,  cn, email  ,  description",
	}

	valid, attributes := ldap.ValidateAttributes()
	c.Assert(valid, check.Equals, true)
	c.Assert(attributes, check.Equals, "givenName,cn,email,description")
}

func (s *ModelsSuite) TestLDAPConnect(c *check.C) {
	l := LDAP{
		UserId:     1,
		Name:       "Default",
		Protocol:   "ldap",
		Host:       "ldap.forumsys.com:389",
		Username:   "cn=read-only-admin,dc=example,dc=com",
		Password:   "password",
		BaseDn:     "dc=example,dc=com",
		Query:      "(objectClass=inetOrgPerson)",
		Attributes: "cn,sn,email",
	}

	conn, err := l.Connect()
	c.Assert(err, check.Equals, nil)
	c.Assert(conn, check.NotNil)
}

func (s *ModelsSuite) TestLDAPImportUsers(c *check.C) {
	l := LDAP{
		UserId:     1,
		Name:       "Default",
		Protocol:   "ldap",
		Host:       "ldap.forumsys.com:389",
		Username:   "cn=read-only-admin,dc=example,dc=com",
		Password:   "password",
		BaseDn:     "dc=example,dc=com",
		Query:      "(objectClass=inetOrgPerson)",
		Attributes: "cn,sn,mail",
	}

	expectedResultsLength := 15

	conn, err := l.Connect()
	c.Assert(err, check.Equals, nil)
	c.Assert(conn, check.NotNil)

	targets, err := l.ImportLDAPTargets(conn)
	c.Assert(err, check.Equals, nil)
	c.Assert(targets, check.HasLen, expectedResultsLength)

	// Albert Einstein should be in there
	einsteinFound := false
	for _, t := range targets {
		if t.LastName == "Einstein" { // This asserts also that 'sn' is 'Einstein'
			einsteinFound = true
			// The first attribute 'cn' is the whole name in this server
			c.Assert(t.FirstName, check.Equals, "Albert Einstein")
			c.Assert(t.Email, check.Equals, "einstein@ldap.forumsys.com")
			// Position should be empty since not in attributes
			c.Assert(t.Position, check.Equals, "")
			break
		}
	}
	c.Assert(einsteinFound, check.Equals, true)
}

func (s *ModelsSuite) TestGetInvalidLDAP(ch *check.C) {
	_, err := GetLDAP(-1, 1)
	ch.Assert(err, check.Equals, gorm.ErrRecordNotFound)
}
