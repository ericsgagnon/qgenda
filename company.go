package main

import (
	"time"

	"github.com/google/uuid"
)

// Company contains basic company info
type Company struct {
	ID                 uuid.UUID `json:"CompanyKey"`
	Name               string    `json:"CompanyName"`
	Abbreviation       string    `json:"CompanyAbbr"`
	CreatedTime        time.Time `json:"DateCreatedUtc"`
	CompanyLocation    `json:""`
	CompanyPhoneNumber `json:""`
	Profiles           `json:""`
	ProfileName        `json:""`
	ProfileKey         `json:""`
	IsAdmin            `json:""`
	Organizations      `json:""`
	OrgName            `json:""`
	OrgKey             `json:""`
}
