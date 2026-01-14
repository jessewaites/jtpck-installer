class TeamSnapshot < ApplicationRecord
  belongs_to :organization
  belongs_to :team

  scope :on_date, ->(date) { where(captured_on: date) }
  scope :for_range, ->(range) { where(captured_on: range) }
end
