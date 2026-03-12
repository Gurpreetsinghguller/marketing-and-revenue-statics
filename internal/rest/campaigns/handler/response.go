package handler

type CampaignPreviewResponse struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Status      string `json:"status"`
}

type DeleteCampaignResponse struct {
	Message string `json:"message"`
}
