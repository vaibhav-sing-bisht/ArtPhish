package models

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/go-ldap/ldap/v3"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type LDAP struct {
	Id       int64  `json:"id" gorm:"column:id; primary_key:yes"`
	UserId   int64  `json:"-" gorm:"column:user_id"`
	Name     string `json:"name"`
	Protocol string `json:"protocol"`
	Host     string `json:"host"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	BaseDn   string `json:"base_dn" gorm:"column:base_dn"`
	Query    string `json:"query"`
	// Attributes should be a comma separated list of the ldap server attributes
	// representing in this order first name, last name, email, and optionally position
	// E.g: firstName,lastName,mail,function
	Attributes       string    `json:"attributes"`
	ModifiedDate     time.Time `json:"modified_date"`
	IgnoreCertErrors bool      `json:"ignore_cert_errors"`
}

// ErrLdapProtocolNotSpecified is thrown when there is no protocol specified
var ErrLdapProtocolNotSpecified = errors.New("No LDAP Protocol specified")

// ErrLdapProtocolInvalid is thrown when protocol is not ldap or ldaps
var ErrLdapProtocolInvalid = errors.New("LDAP protocol must be ldap or ldaps")

// ErrLdapHostNotSpecified is thrown when there is no Host specified
var ErrLdapHostNotSpecified = errors.New("No LDAP Host specified")

// ErrInvalidLdapHost indicates that the SMTP server string is invalid
var ErrInvalidLdapHost = errors.New("Invalid LDAP server address")

// ErrLdapBaseDnNotSpecified is thrown when there is no Base search DN specified
var ErrLdapBaseDnNotSpecified = errors.New("No LDAP Base DN specified")

// ErrLdapQueryNotSpecified is thrown when there is no Root DN specified
var ErrLdapQueryNotSpecified = errors.New("No LDAP Query specified")

// ErrLdapAttributesNotSpecified is thrown when there is no Attributes
var ErrLdapAttributesNotSpecified = errors.New("No LDAP Attributes specified")

// ErrLdapAttributesInvalid is thrown when attributes format is invalid
var ErrLdapAttributesInvalid = errors.New("LDAP Attributes Invalid")

// TableName specifies the database tablename for Gorm to use
func (l LDAP) TableName() string {
	return "ldap"
}

func (l *LDAP) ValidateAttributes() (bool, string) {
	parts := strings.Split(l.Attributes, ",")

	if len(parts) < 3 || len(parts) > 4 {
		return false, l.Attributes
	}

	for i, part := range parts {
		parts[i] = strings.TrimSpace(part)
	}

	return true, strings.Join(parts, ",")
}

// Validate ensures that LDAP configs/connections are valid
func (l *LDAP) Validate() error {
	switch {
	case l.Protocol == "":
		return ErrLdapProtocolNotSpecified
	case l.Protocol != "ldap" && l.Protocol != "ldaps":
		return ErrLdapProtocolInvalid
	case l.Host == "":
		return ErrLdapHostNotSpecified
	case l.BaseDn == "":
		return ErrLdapBaseDnNotSpecified
	case l.Query == "":
		return ErrLdapQueryNotSpecified
	case l.Attributes == "":
		return ErrLdapAttributesNotSpecified
	}

	valid, attributes := l.ValidateAttributes()
	if !valid {
		return ErrLdapAttributesInvalid
	}
	l.Attributes = attributes

	// Make sure addr is in host:port format
	hp := strings.Split(l.Host, ":")
	if len(hp) > 2 {
		return ErrInvalidLdapHost
	} else if len(hp) < 2 {
		hp = append(hp, "389")
	}

	_, err := strconv.Atoi(hp[1])
	if err != nil {
		return ErrInvalidLdapHost
	}
	return err
}

// Connect returns a connection for the given LDAP profile
func (l *LDAP) Connect() (*ldap.Conn, error) {
	hp := strings.Split(l.Host, ":")
	if len(hp) < 2 {
		hp = append(hp, "389")
	}
	host := hp[0]
	// Any issues should have been caught in validation, but we'll
	// double check here.
	port, err := strconv.Atoi(hp[1])
	if err != nil {
		log.Error(err)
		return nil, err
	}

	var conn *ldap.Conn
	url := fmt.Sprintf("%s://%s:%d", l.Protocol, host, port)

	log.Infof("Connecting to LDAP Server `%s`", url)

	if l.Protocol == "ldaps" {
		conn, err = ldap.DialURL(url, ldap.DialWithTLSConfig(&tls.Config{
			InsecureSkipVerify: true,
		}))
	} else {
		conn, err = ldap.DialURL(url)
	}

	if err != nil {
		log.Error("Failed to connect to ldap server", err)
		return nil, err
	}

	if l.Username != "" {
		err = conn.Bind(l.Username, l.Password)
		if err != nil {
			log.Errorf("Failed to bind LDAP user `%s`", l.Username)
			return nil, err
		}
	}

	return conn, nil
}

func (l *LDAP) ImportLDAPTargets(conn *ldap.Conn) ([]Target, error) {
	ts := []Target{}
	attributes := strings.Split(l.Attributes, ",")
	for i, attribute := range attributes {
		attributes[i] = strings.TrimSpace(attribute)
	}

	request := ldap.NewSearchRequest(
		l.BaseDn,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		l.Query,
		attributes,
		nil,
	)

	results, err := conn.Search(request)
	if err != nil {
		log.Error("Error importing users", err)
		return ts, err
	}

	for _, entry := range results.Entries {
		fn := entry.GetAttributeValue(attributes[0])
		ln := entry.GetAttributeValue(attributes[1])
		ea := entry.GetAttributeValue(attributes[2])
		ps := ""
		if len(attributes) == 4 {
			ps = entry.GetAttributeValue(attributes[3])
		}

		t := Target{
			BaseRecipient: BaseRecipient{
				FirstName: fn,
				LastName:  ln,
				Email:     ea,
				Position:  ps,
			},
		}
		ts = append(ts, t)
	}

	return ts, nil
}

// GetLDAPs returns the LDAPs owned by the given user.
func GetLDAPs(uid int64) ([]LDAP, error) {
	ls := []LDAP{}
	err := db.Where("user_id=?", uid).Find(&ls).Error
	if err != nil {
		log.Error(err)
		return ls, err
	}

	return ls, nil
}

// GetLDAP returns the LDAP, if it exists, specified by the given id and user_id.
func GetLDAP(id int64, uid int64) (LDAP, error) {
	l := LDAP{}
	err := db.Where("user_id=? and id=?", uid, id).Find(&l).Error
	if err != nil {
		log.Error(err)
		return l, err
	}

	return l, err
}

// GetLDAPByName returns the LDAP, if it exists, specified by the given name and user_id.
func GetLDAPByName(n string, uid int64) (LDAP, error) {
	l := LDAP{}
	err := db.Where("user_id=? and name=?", uid, n).Find(&l).Error
	if err != nil {
		log.Error(err)
		return l, err
	}

	return l, err
}

// PostLDAP creates a new LDAP in the database.
func PostLDAP(l *LDAP) error {
	err := l.Validate()
	if err != nil {
		log.Error(err)
		return err
	}
	// Insert into the DB
	err = db.Save(l).Error
	if err != nil {
		log.Error(err)
	}

	return err
}

// PutLDAP edits an existing LDAP in the database.
// Per the PUT Method RFC, it presumes all data for a LDAP is provided.
func PutLDAP(l *LDAP) error {
	err := l.Validate()
	if err != nil {
		log.Error(err)
		return err
	}
	err = db.Where("id=?", l.Id).Save(l).Error
	if err != nil {
		log.Error(err)
	}

	return err
}

// DeleteLDAP deletes an existing LDAP in the database.
// An error is returned if a LDAP with the given user id and LDAP id is not found.
func DeleteLDAP(id int64, uid int64) error {
	err := db.Where("user_id=?", uid).Delete(LDAP{Id: id}).Error
	if err != nil {
		log.Error(err)
	}

	return err
}
