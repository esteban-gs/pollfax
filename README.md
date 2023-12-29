# Pollfax

## Summary
This is an educational project for the continued learning of:
- Data ingestion pipelines
- ETL
- Go(lang)

## Architecture
Golang as the server language paired with sqlx package for interfacing
with the database.

### Data ingestion
The server starts a cron job to ingest data from api.congress.gov. The main
data points to ingest are:
- Current congress number. e.g 118th congress, etc.
- Get first ~250 bills belonging to the 118th congress.
- Persist the bills to a transient bills table.

### Roadmap
The pollfax MVP will accomplish the following
- Let users browse the most recently updated bills.
- Users can up/downvote bills -> this is persisted in a database
- Dashboard TBD


## Database Schema Version Control

### Create New Migration
- run `migrate create -ext sqlx -dir db/migrations -seq create_bills_sentiment_table`
- add the necessary sql
