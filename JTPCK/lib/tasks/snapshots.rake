namespace :snapshots do
  desc "Generate daily snapshots for yesterday (usage: rake snapshots:daily)"
  task daily: :environment do
    date = Date.yesterday
    puts "Generating snapshots for #{date}..."
    GenerateSnapshotsJob.perform_now(date)
    puts "Done."
  end

  desc "Backfill snapshots for a date range (usage: rake snapshots:backfill[2025-12-01,2025-12-31])"
  task :backfill, %i[start_date end_date] => :environment do |_, args|
    unless args[:start_date] && args[:end_date]
      abort "Usage: rake snapshots:backfill[START_DATE,END_DATE]"
    end

    range = Date.parse(args[:start_date])..Date.parse(args[:end_date])
    range.each do |date|
      puts "Generating snapshots for #{date}..."
      GenerateSnapshotsJob.perform_now(date)
    end

    puts "Backfill complete for #{range}"
  end
end
