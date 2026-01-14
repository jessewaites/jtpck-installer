class SnapshotAggregator
  SCHEMA_VERSION = 1

  attr_reader :target_date, :logger

  def initialize(target_date: Date.yesterday, logger: Rails.logger)
    @target_date = target_date
    @logger = logger || Logger.new($stdout)
  end

  def run!
    logger.info("[SnapshotAggregator] Building snapshots for #{target_date}")
    range = time_range_for_date(target_date)
    upsert_org_snapshots(range)
    upsert_team_snapshots(range)
    upsert_user_snapshots(range)
  end

  private

  def time_range_for_date(date)
    date = date.to_date
    Time.zone.local(date.year, date.month, date.day).all_day
  end

  def base_select
    <<~SQL.squish
      COUNT(*) AS event_count,
      SUM(CASE WHEN success THEN 1 ELSE 0 END) AS success_count,
      SUM(CASE WHEN success THEN 0 ELSE 1 END) AS failure_count,
      SUM(total_tokens) AS total_tokens,
      SUM(input_tokens) AS input_tokens,
      SUM(output_tokens) AS output_tokens,
      SUM(latency_ms) AS total_latency_ms,
      AVG(latency_ms) AS avg_latency_ms
    SQL
  end

  def upsert_org_snapshots(range)
    rows = UsageEvent
      .where(occurred_at: range)
      .group(:organization_id)
      .select(:organization_id, Arel.sql(base_select))
      .map { |record| build_snapshot_row(record, captured_on: target_date) }

    return if rows.empty?

    OrgSnapshot.upsert_all(rows, unique_by: :index_org_snapshots_on_org_and_date)
  end

  def upsert_team_snapshots(range)
    rows = UsageEvent
      .where(occurred_at: range)
      .group(:organization_id, :team_id)
      .select(:organization_id, :team_id, Arel.sql(base_select))
      .map { |record| build_snapshot_row(record, captured_on: target_date) }

    return if rows.empty?

    TeamSnapshot.upsert_all(rows, unique_by: :index_team_snapshots_on_team_and_date)
  end

  def upsert_user_snapshots(range)
    rows = UsageEvent
      .where(occurred_at: range)
      .where.not(user_id: nil)
      .group(:organization_id, :team_id, :user_id)
      .select(:organization_id, :team_id, :user_id, Arel.sql(base_select))
      .map { |record| build_snapshot_row(record, captured_on: target_date) }

    return if rows.empty?

    UserSnapshot.upsert_all(rows, unique_by: :index_user_snapshots_on_user_and_date)
  end

  def build_snapshot_row(record, captured_on:)
    {
      organization_id: record.try(:organization_id),
      team_id: record.try(:team_id),
      user_id: record.try(:user_id),
      captured_on: captured_on,
      schema_version: SCHEMA_VERSION,
      event_count: record.event_count.to_i,
      success_count: record.success_count.to_i,
      failure_count: record.failure_count.to_i,
      total_tokens: record.total_tokens.to_i,
      input_tokens: record.input_tokens.to_i,
      output_tokens: record.output_tokens.to_i,
      total_latency_ms: record.total_latency_ms.to_i,
      avg_latency_ms: record.avg_latency_ms&.to_i,
      created_at: Time.current,
      updated_at: Time.current
    }.compact
  end
end
