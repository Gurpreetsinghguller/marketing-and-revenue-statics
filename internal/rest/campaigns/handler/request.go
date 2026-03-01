package handler

import "time"

type CreateCampaignRequest struct {
	Name        string  `json:"name" validate:"required,min=1"`
	Description string  `json:"description" validate:"omitempty"`
	Budget      float64 `json:"budget" validate:"gte=0"`
	Channel     string  `json:"channel" validate:"required,min=1"`
	IsPublic    bool    `json:"is_public"`
}

type UpdateCampaignRequest struct {
	Name        *string  `json:"name,omitempty" validate:"required,min=1"`
	Description *string  `json:"description,omitempty" validate:"omitempty"`
	Budget      *float64 `json:"budget,omitempty" validate:"omitempty,gte=0"`
	IsPublic    *bool    `json:"is_public,omitempty"`
	Status      *string  `json:"status,omitempty" validate:"required,oneof=active paused completed"`
}

type DateRange struct {
	Start *time.Time `json:"start,omitempty" validate:"required"`
	End   *time.Time `json:"end,omitempty" validate:"required,gtefield=Start"`
}

type PatchCampaignStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=active paused completed"`
}

type SearchCampaignsRequest struct {
	Query string `json:"query" validate:"required,min=1"`
}
