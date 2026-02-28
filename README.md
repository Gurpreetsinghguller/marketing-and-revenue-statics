# Marketing & Revenue Statistics

## Objective

Create scalable and performant backend REST APIs for Marketing & Revenue Statistics. The system should track campaign performance, user engagement, and revenue metrics, and provide aggregated analytics reports for different roles like marketers and analysts. Use the DRF (Django Rest Framework) or NodeJs (_based on the job role_) for the same.

You're encouraged to use tools like Cursor or GitHub Copilot. But your solution should show architectural decision-making.

## Part 1: Authentication and User Roles

- A user profile should have basic information, i.e., `Name`, `Role`, `Bio`, `Picture`, and `Phone number`.
- Create necessary APIs for the following use cases.
  - A user should be able to register & login using Email & Password.
    - You can use JWT to implement this.
    - No 3rd party ( Google, Github, … ) auth integration is required.
  - A user should be able to update profile information.
- Implement role-based access control:
  - Admin: Full access
  - Marketer: Access to own campaigns and performance reports
  - Analyst: Read-only access to reports

### Unauthenticated access (limited)

- Can view publicly published campaign stats summary (anonymized)
- Can view landing pages or campaign preview info

### All other API endpoints require authentication

## Part 2: Campaign Management

- Marketers can:
  - Create, update, delete, activate, pause campaigns.
- Anyone (unauthenticated):
  - Can view campaign preview info (read-only, if marked public).
- Authenticated users:
  - Can get a list of campaigns with filters:
    - status, date_range, created_by, channel etc
  - Can search campaigns by name/description

## Part 3: Event Tracking and Metrics Collection

- Event Types: impression, click, conversion
- Public/External systems can send tracking events via API (clicks, impressions, conversions).
- Events are saved to a time-series DB or table (EventLog).

## Part 4: Analytics & Reporting APIs

- Provide aggregated stats based on time period, campaign, or user-level data.
- Marketers/Analysts can get:
  - Daily/weekly/monthly performance summaries
- Users can filter by: campaign, event_type, date_range, channel etc

### Metrics Definitions

| Metric                     | Formula                                              |
| -------------------------- | ---------------------------------------------------- |
| CTR (Click Through Rate)   | (Total Clicks / Total Impressions) \* 100            |
| CPC (Cost Per Click)       | Total Spend / Total Clicks                           |
| ROI (Return on Investment) | ((Total Revenue - Total Spend) / Total Spend) \* 100 |
| Conversion Rate            | (Total Conversions / Total Clicks) \* 100            |

## Part 5: User Engagement Tracking

- Capture behavioral data across campaigns.
- Track:
  - Time spent on campaign pages
  - Click path (funnels)
  - Drop-off rates between stages (e.g., Ad → Landing → Signup → Purchase)

## Task-6: Unit Tests and API Documentation

- Write unit test cases for the developed functionality
- Create API documentation for the developed endpoints

## Bonus Enhancements (Optional)

- Email notifications should be triggered (_using Django signals, or NodeJS events_)
  - When CTR drops below X%
  - When the budget is exceeded
- Add rate limit in the API
  - Without an authenticated API, anonymous user-based
  - With an authenticated API, user-based

### 🧠 Think Piece

- How would you separate write-heavy workloads (event ingestion) from read-heavy workloads (analytics/reporting)?
- How would you design a pipeline to stream and aggregate campaign performance data in near real-time?
- Would you pre-aggregate daily metrics? If so, how would you store and update them?
- How would you structure your API to support exporting large datasets without timing out?
- How would you reduce cloud storage and compute costs for storing clickstream or impression logs?

## 📂 Submission

- Push to the assigned repository.
- Create `Architecture.md` with:
  - **Setup instructions**
  - **Architecture explanation** (1-2 paragraphs)
  - **Answers to Think Pieces** above
  - Mention any use of AI tooling (e.g., Cursor, GitHub Copilot)

## 📝 Evaluation Criteria

Your submission will be evaluated based on the following criteria:

- Think of this as an app you create during day-to-day operations. You should invest more time to make it rightly engineered than submitting it early.
- Ensure that the APIs meet all the core requirements and behave as expected.
- Code quality & best practices
- API design & developer experience
- Performance & scalability considerations



