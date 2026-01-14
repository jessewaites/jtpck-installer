class OrgSnapshot < ApplicationRecord
  belongs_to :organization

  scope :on_date, ->(date) { where(captured_on: date) }
  scope :for_range, ->(range) { where(captured_on: range) }
end
