package report

import (
	"context"
	"encoding/json"

	"github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/domain"
	"github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/persistence/db"
)

const reportPrefix = "reports"

type reportDBModel struct {
	CampaignID       string                    `json:"campaign_id"`
	CampaignName     string                    `json:"campaign_name"`
	Period           string                    `json:"period"`
	StartDate        string                    `json:"start_date"`
	EndDate          string                    `json:"end_date"`
	ReportType       string                    `json:"report_type"`
	TotalImpressions int64                     `json:"total_impressions"`
	TotalClicks      int64                     `json:"total_clicks"`
	TotalConversions int64                     `json:"total_conversions"`
	TotalRevenue     float64                   `json:"total_revenue"`
	TotalSpend       float64                   `json:"total_spend"`
	Metrics          domain.Metrics            `json:"metrics"`
	ByChannel        []domain.ChannelBreakdown `json:"by_channel"`
	ByDevice         []domain.DeviceBreakdown  `json:"by_device"`
}

// ReportRepository implements domain.ReportRepo using file-backed persistence.
type ReportRepository struct {
	storage db.PersistenceDB
}

// NewReportRepository creates a new report repository.
func NewReportRepository(storage ...db.PersistenceDB) *ReportRepository {
	selected := db.PersistenceDB(db.NewStorageMgr())
	if len(storage) > 0 && storage[0] != nil {
		selected = storage[0]
	}
	return &ReportRepository{storage: selected}
}

// Create saves a pre-aggregated report.
func (r *ReportRepository) Create(report *domain.Report) error {
	if report.ReportType == "" {
		report.ReportType = "custom"
	}
	return r.storage.Create(
		context.Background(),
		reportKey(report.CampaignID, report.StartDate, report.EndDate, report.ReportType),
		toReportDBModel(*report),
	)
}

func (r *ReportRepository) getAll() ([]domain.Report, error) {
	stored, err := r.storage.List(context.Background(), reportPrefix)
	if err != nil {
		return nil, err
	}

	reports := make([]domain.Report, 0, len(stored))
	for _, item := range stored {
		model, err := decodeReportDBModel(item)
		if err != nil {
			return nil, err
		}
		reports = append(reports, model.toDomain())
	}
	return reports, nil
}

// GetByID retrieves a report by ID (mapped to CampaignID for current domain contract).
func (r *ReportRepository) GetByID(id string) (*domain.Report, error) {
	reports, err := r.getAll()
	if err != nil {
		return nil, err
	}
	for i := range reports {
		if reports[i].CampaignID == id {
			return &reports[i], nil
		}
	}
	return nil, db.ErrNotFound
}

// GetByCampaignAndDateRange retrieves report for a campaign and date range.
func (r *ReportRepository) GetByCampaignAndDateRange(campaignID, startDate, endDate string) (*domain.Report, error) {
	reports, err := r.getAll()
	if err != nil {
		return nil, err
	}
	for i := range reports {
		if reports[i].CampaignID == campaignID &&
			reports[i].StartDate == startDate &&
			reports[i].EndDate == endDate {
			return &reports[i], nil
		}
	}
	return nil, db.ErrNotFound
}

// GetDaily retrieves daily report for a campaign.
func (r *ReportRepository) GetDaily(campaignID, date string) (*domain.Report, error) {
	reports, err := r.getAll()
	if err != nil {
		return nil, err
	}
	for i := range reports {
		if reports[i].CampaignID == campaignID && reports[i].ReportType == "daily" && reports[i].StartDate == date {
			return &reports[i], nil
		}
	}
	return nil, db.ErrNotFound
}

// GetWeekly retrieves weekly report for a campaign.
func (r *ReportRepository) GetWeekly(campaignID, year, week string) (*domain.Report, error) {
	reports, err := r.getAll()
	if err != nil {
		return nil, err
	}
	for i := range reports {
		if reports[i].CampaignID == campaignID && reports[i].ReportType == "weekly" {
			return &reports[i], nil
		}
	}
	return nil, db.ErrNotFound
}

// GetMonthly retrieves monthly report for a campaign.
func (r *ReportRepository) GetMonthly(campaignID, year, month string) (*domain.Report, error) {
	reports, err := r.getAll()
	if err != nil {
		return nil, err
	}
	for i := range reports {
		if reports[i].CampaignID == campaignID && reports[i].ReportType == "monthly" {
			return &reports[i], nil
		}
	}
	return nil, db.ErrNotFound
}

// GetAll retrieves all reports.
func (r *ReportRepository) GetAll() ([]domain.Report, error) {
	return r.getAll()
}

func toReportDBModel(r domain.Report) reportDBModel {
	return reportDBModel{
		CampaignID:       r.CampaignID,
		CampaignName:     r.CampaignName,
		Period:           r.Period,
		StartDate:        r.StartDate,
		EndDate:          r.EndDate,
		ReportType:       r.ReportType,
		TotalImpressions: r.TotalImpressions,
		TotalClicks:      r.TotalClicks,
		TotalConversions: r.TotalConversions,
		TotalRevenue:     r.TotalRevenue,
		TotalSpend:       r.TotalSpend,
		Metrics:          r.Metrics,
		ByChannel:        r.ByChannel,
		ByDevice:         r.ByDevice,
	}
}

func (m reportDBModel) toDomain() domain.Report {
	return domain.Report{
		CampaignID:       m.CampaignID,
		CampaignName:     m.CampaignName,
		Period:           m.Period,
		StartDate:        m.StartDate,
		EndDate:          m.EndDate,
		ReportType:       m.ReportType,
		TotalImpressions: m.TotalImpressions,
		TotalClicks:      m.TotalClicks,
		TotalConversions: m.TotalConversions,
		TotalRevenue:     m.TotalRevenue,
		TotalSpend:       m.TotalSpend,
		Metrics:          m.Metrics,
		ByChannel:        m.ByChannel,
		ByDevice:         m.ByDevice,
	}
}

func decodeReportDBModel(value interface{}) (reportDBModel, error) {
	b, err := json.Marshal(value)
	if err != nil {
		return reportDBModel{}, err
	}
	var model reportDBModel
	if err := json.Unmarshal(b, &model); err != nil {
		return reportDBModel{}, err
	}
	return model, nil
}

func reportKey(campaignID, startDate, endDate, reportType string) string {
	return reportPrefix + "/" + campaignID + "_" + startDate + "_" + endDate + "_" + reportType
}
