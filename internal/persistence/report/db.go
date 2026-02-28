package report

import (
	"github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/domain"
)

// ReportRepository implements domain.ReportRepo using JSON file storage
type ReportRepository struct {
	storage StorageMgr
	reports []domain.Report
}

// NewReportRepository creates a new report repository
func NewReportRepository() *ReportRepository {
	repo := &ReportRepository{
		storage: NewStorageMgr(),
		reports: []domain.Report{},
	}
	// Load existing reports from file
	repo.storage.ReadJSON(ReportsFile, &repo.reports)
	return repo
}

// Create saves a pre-aggregated report
func (r *ReportRepository) Create(report *domain.Report) error {
	r.storage.mu.Lock()
	defer r.storage.mu.Unlock()

	// Generate ID if not provided
	if report.ReportType == "" {
		report.ReportType = "custom"
	}

	r.reports = append(r.reports, *report)
	return r.storage.WriteJSON(ReportsFile, r.reports)
}

// GetByID retrieves a report by ID
func (r *ReportRepository) GetByID(id string) (*domain.Report, error) {
	r.storage.mu.RLock()
	defer r.storage.mu.RUnlock()

	for i := range r.reports {
		if r.reports[i].CampaignID == id {
			return &r.reports[i], nil
		}
	}
	return nil, ErrNotFound
}

// GetByCampaignAndDateRange retrieves report for a campaign and date range
func (r *ReportRepository) GetByCampaignAndDateRange(campaignID, startDate, endDate string) (*domain.Report, error) {
	r.storage.mu.RLock()
	defer r.storage.mu.RUnlock()

	for i := range r.reports {
		if r.reports[i].CampaignID == campaignID &&
			r.reports[i].StartDate == startDate &&
			r.reports[i].EndDate == endDate {
			return &r.reports[i], nil
		}
	}
	return nil, ErrNotFound
}

// GetDaily retrieves daily report for a campaign
func (r *ReportRepository) GetDaily(campaignID, date string) (*domain.Report, error) {
	r.storage.mu.RLock()
	defer r.storage.mu.RUnlock()

	for i := range r.reports {
		if r.reports[i].CampaignID == campaignID &&
			r.reports[i].ReportType == "daily" &&
			r.reports[i].StartDate == date {
			return &r.reports[i], nil
		}
	}
	return nil, ErrNotFound
}

// GetWeekly retrieves weekly report for a campaign
func (r *ReportRepository) GetWeekly(campaignID, year, week string) (*domain.Report, error) {
	r.storage.mu.RLock()
	defer r.storage.mu.RUnlock()

	for i := range r.reports {
		if r.reports[i].CampaignID == campaignID &&
			r.reports[i].ReportType == "weekly" {
			// Match by week (simple implementation)
			return &r.reports[i], nil
		}
	}
	return nil, ErrNotFound
}

// GetMonthly retrieves monthly report for a campaign
func (r *ReportRepository) GetMonthly(campaignID, year, month string) (*domain.Report, error) {
	r.storage.mu.RLock()
	defer r.storage.mu.RUnlock()

	for i := range r.reports {
		if r.reports[i].CampaignID == campaignID &&
			r.reports[i].ReportType == "monthly" {
			// Match by month (simple implementation)
			return &r.reports[i], nil
		}
	}
	return nil, ErrNotFound
}

// GetAll retrieves all reports
func (r *ReportRepository) GetAll() ([]domain.Report, error) {
	r.storage.mu.RLock()
	defer r.storage.mu.RUnlock()

	return r.reports, nil
}
