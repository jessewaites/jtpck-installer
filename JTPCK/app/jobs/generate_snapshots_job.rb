class GenerateSnapshotsJob < ApplicationJob
  queue_as :default

  def perform(target_date = Date.yesterday)
    SnapshotAggregator.new(target_date: target_date).run!
  end
end
