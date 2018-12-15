package models

import (
	"errors"
	"time"

	"github.com/gosu-team/fptu-api/config"
)

// Confession ...
type Confession struct {
	ID        int        `json:"id" gorm:"primary_key"`
	CreatedAt *time.Time `json:"created_at, omitempty"`
	UpdatedAt *time.Time `json:"updated_at, omitempty"`
	DeletedAt *time.Time `json:"deleted_at, omitempty" sql:"index"`

	Content  string `json:"content" gorm:"not null; type:text;"`
	Sender   string `json:"sender" gorm:"not null; type:varchar(250);"`
	Status   int    `json:"status" gorm:"not null; type:int(11);"`
	Approver int    `json:"approver" gorm:"type:int(11);"`
	Reason   string `json:"reason" gorm:"type:varchar(250);"`
	CfsID    int    `json:"cfs_id" gorm:"type:int(11);"`
}

// TableName set Confession's table name to be `confessions`
func (Confession) TableName() string {
	return "confessions"
}

// FetchAll ...
func (c *Confession) FetchAll(numLoad int) []Confession {
	db := config.GetDatabaseConnection()

	var confessions []Confession
	db.Order("id desc").Limit(numLoad).Find(&confessions)

	return confessions
}

// FetchByID ...
func (c *Confession) FetchByID() error {
	db := config.GetDatabaseConnection()

	if err := db.Where("id = ?", c.ID).Find(&c).Error; err != nil {
		return errors.New("Could not find the confession")
	}

	return nil
}

// Create ...
func (c *Confession) Create() error {
	db := config.GetDatabaseConnection()

	// Validate record
	if !db.NewRecord(c) { // => returns `true` as primary key is blank
		return errors.New("New records can not have primary key id")
	}

	if err := db.Create(&c).Error; err != nil {
		return errors.New("Could not create confession")
	}

	return nil
}

// Save ...
func (c *Confession) Save() error {
	db := config.GetDatabaseConnection()

	if db.NewRecord(c) {
		if err := db.Create(&c).Error; err != nil {
			return errors.New("Could not create confessions")
		}
	} else {
		if err := db.Save(&c).Error; err != nil {
			return errors.New("Could not update confessions")
		}
	}
	return nil
}

// FetchBySender ...
func (c *Confession) FetchBySender(sender string, numLoad int) []Confession {
	db := config.GetDatabaseConnection()

	var confessions []Confession
	db.Order("id desc").Limit(numLoad).Where("sender = ?", sender).Find(&confessions)

	return confessions
}

// FetchOverview ...
func (c *Confession) FetchOverview() (int, int, int) {
	db := config.GetDatabaseConnection()

	totalCount, pendingCount, rejectedCount := 0, 0, 0
	db.Model(&Confession{}).Count(&totalCount)
	db.Model(&Confession{}).Where("status = ?", 0).Count(&pendingCount)
	db.Model(&Confession{}).Where("status = ?", 2).Count(&rejectedCount)

	return totalCount, pendingCount, rejectedCount
}

// GetNextConfessionID ...
func (c *Confession) GetNextConfessionID() int {
	db := config.GetDatabaseConnection()
	db.Order("cfs_id desc").Take(&c)
	return c.CfsID + 1
}

func (c *Confession) setConfessionApproved(status int, approver int, cfsID int) {
	c.Status = status
	c.Approver = approver
	c.CfsID = cfsID
}

// ApproveConfession ...
func (c *Confession) ApproveConfession(approverID int) error {
	if err := c.FetchByID(); err != nil {
		return errors.New("Could not find the confession")
	}

	if c.Status != 0 {
		return errors.New("Status of confession must be pending to be approved")
	}

	confessions := new(Confession)

	c.setConfessionApproved(1, approverID, confessions.GetNextConfessionID())

	if err := c.Save(); err != nil {
		return errors.New("Unable to update approved confession`")
	}

	return nil
}

func (c *Confession) setConfessionUnapproved() {
	c.Status = 0
	c.Approver = 0
	c.CfsID = 0
}

// RollbackApproveConfession ...
func (c *Confession) RollbackApproveConfession(approverID int) error {
	if err := c.FetchByID(); err != nil {
		return errors.New("Could not find the confession")
	}

	c.setConfessionUnapproved()

	if err := c.Save(); err != nil {
		return errors.New("Unable to rollback approved confession`")
	}

	return nil
}

func (c *Confession) setConfessionRejected(status int, approver int, reason string) {
	c.Status = status
	c.Approver = approver
	c.Reason = reason
}

// RejectConfession ...
func (c *Confession) RejectConfession(approverID int, reason string) error {
	if err := c.FetchByID(); err != nil {
		return errors.New("Could not find the confession")
	}

	if c.Status != 0 {
		return errors.New("Status of confession must be pending to be rejected")
	}

	c.setConfessionRejected(2, approverID, reason)

	if err := c.Save(); err != nil {
		return errors.New("Unable to update approved confession`")
	}

	return nil
}
