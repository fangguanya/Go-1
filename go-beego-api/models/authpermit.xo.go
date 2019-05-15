// Package models contains the types for schema 'dbo'.
package models

// Code generated by xo. DO NOT EDIT.

import (
	"errors"
	"time"
)

// Authpermit represents a row from 'dbo.AuthPermit'.
type Authpermit struct {
	Code        string    `json:"Code"`        // Code
	Name        string    `json:"Name"`        // Name
	Revision    int       `json:"Revision"`    // Revision
	Createdby   string    `json:"CreatedBy"`   // CreatedBy
	Createdtime time.Time `json:"CreatedTime"` // CreatedTime
	Updatedby   string    `json:"UpdatedBy"`   // UpdatedBy
	Updatedtime time.Time `json:"UpdatedTime"` // UpdatedTime

	// xo fields
	_exists, _deleted bool
}

// Exists determines if the Authpermit exists in the database.
func (a *Authpermit) Exists() bool {
	return a._exists
}

// Deleted provides information if the Authpermit has been deleted from the database.
func (a *Authpermit) Deleted() bool {
	return a._deleted
}

// Insert inserts the Authpermit to the database.
func (a *Authpermit) Insert(db XODB) error {
	var err error

	// if already exist, bail
	if a._exists {
		return errors.New("insert failed: already exists")
	}

	// sql insert query, primary key must be provided
	const sqlstr = `INSERT INTO dbo.AuthPermit (` +
		`Code, Name, Revision, CreatedBy, CreatedTime, UpdatedBy, UpdatedTime` +
		`) VALUES (` +
		`$1, $2, $3, $4, $5, $6, $7` +
		`)`

	// run query
	XOLog(sqlstr, a.Code, a.Name, a.Revision, a.Createdby, a.Createdtime, a.Updatedby, a.Updatedtime)
	_, err = db.Exec(sqlstr, a.Code, a.Name, a.Revision, a.Createdby, a.Createdtime, a.Updatedby, a.Updatedtime)
	if err != nil {
		return err
	}

	// set existence
	a._exists = true

	return nil
}

// Update updates the Authpermit in the database.
func (a *Authpermit) Update(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !a._exists {
		return errors.New("update failed: does not exist")
	}

	// if deleted, bail
	if a._deleted {
		return errors.New("update failed: marked for deletion")
	}

	// sql query
	const sqlstr = `UPDATE dbo.AuthPermit SET ` +
		`Name = $1, Revision = $2, CreatedBy = $3, CreatedTime = $4, UpdatedBy = $5, UpdatedTime = $6` +
		` WHERE Code = $7`

	// run query
	XOLog(sqlstr, a.Name, a.Revision, a.Createdby, a.Createdtime, a.Updatedby, a.Updatedtime, a.Code)
	_, err = db.Exec(sqlstr, a.Name, a.Revision, a.Createdby, a.Createdtime, a.Updatedby, a.Updatedtime, a.Code)
	return err
}

// Save saves the Authpermit to the database.
func (a *Authpermit) Save(db XODB) error {
	if a.Exists() {
		return a.Update(db)
	}

	return a.Insert(db)
}

// Delete deletes the Authpermit from the database.
func (a *Authpermit) Delete(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !a._exists {
		return nil
	}

	// if deleted, bail
	if a._deleted {
		return nil
	}

	// sql query
	const sqlstr = `DELETE FROM dbo.AuthPermit WHERE Code = $1`

	// run query
	XOLog(sqlstr, a.Code)
	_, err = db.Exec(sqlstr, a.Code)
	if err != nil {
		return err
	}

	// set deleted
	a._deleted = true

	return nil
}

// AuthpermitByCode retrieves a row from 'dbo.AuthPermit' as a Authpermit.
//
// Generated from index 'PK_AuthPermit'.
func AuthpermitByCode(db XODB, code string) (*Authpermit, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`Code, Name, Revision, CreatedBy, CreatedTime, UpdatedBy, UpdatedTime ` +
		`FROM dbo.AuthPermit ` +
		`WHERE Code = $1`

	// run query
	XOLog(sqlstr, code)
	a := Authpermit{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, code).Scan(&a.Code, &a.Name, &a.Revision, &a.Createdby, &a.Createdtime, &a.Updatedby, &a.Updatedtime)
	if err != nil {
		return nil, err
	}

	return &a, nil
}