class CreateSnapshots < ActiveRecord::Migration[8.1]
  def change
    create_table :org_snapshots do |t|
      t.references :organization, null: false, foreign_key: true
      t.date :captured_on, null: false
      t.integer :schema_version, null: false, default: 1
      t.integer :event_count, null: false, default: 0
      t.integer :success_count, null: false, default: 0
      t.integer :failure_count, null: false, default: 0
      t.bigint :total_tokens, null: false, default: 0
      t.bigint :input_tokens, null: false, default: 0
      t.bigint :output_tokens, null: false, default: 0
      t.bigint :total_latency_ms, null: false, default: 0
      t.integer :avg_latency_ms
      t.timestamps
    end

    add_index :org_snapshots, [:organization_id, :captured_on], unique: true, name: "index_org_snapshots_on_org_and_date"

    create_table :team_snapshots do |t|
      t.references :organization, null: false, foreign_key: true
      t.references :team, null: false, foreign_key: true
      t.date :captured_on, null: false
      t.integer :schema_version, null: false, default: 1
      t.integer :event_count, null: false, default: 0
      t.integer :success_count, null: false, default: 0
      t.integer :failure_count, null: false, default: 0
      t.bigint :total_tokens, null: false, default: 0
      t.bigint :input_tokens, null: false, default: 0
      t.bigint :output_tokens, null: false, default: 0
      t.bigint :total_latency_ms, null: false, default: 0
      t.integer :avg_latency_ms
      t.timestamps
    end

    add_index :team_snapshots, [:team_id, :captured_on], unique: true, name: "index_team_snapshots_on_team_and_date"
    add_index :team_snapshots, [:organization_id, :captured_on], name: "index_team_snapshots_on_org_and_date"

    create_table :user_snapshots do |t|
      t.references :organization, null: false, foreign_key: true
      t.references :team, null: false, foreign_key: true
      t.references :user, null: false, foreign_key: true
      t.date :captured_on, null: false
      t.integer :schema_version, null: false, default: 1
      t.integer :event_count, null: false, default: 0
      t.integer :success_count, null: false, default: 0
      t.integer :failure_count, null: false, default: 0
      t.bigint :total_tokens, null: false, default: 0
      t.bigint :input_tokens, null: false, default: 0
      t.bigint :output_tokens, null: false, default: 0
      t.bigint :total_latency_ms, null: false, default: 0
      t.integer :avg_latency_ms
      t.timestamps
    end

    add_index :user_snapshots, [:user_id, :captured_on], unique: true, name: "index_user_snapshots_on_user_and_date"
    add_index :user_snapshots, [:team_id, :captured_on], name: "index_user_snapshots_on_team_and_date"
    add_index :user_snapshots, [:organization_id, :captured_on], name: "index_user_snapshots_on_org_and_date"
  end
end
