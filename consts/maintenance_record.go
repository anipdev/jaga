package consts

const (
	RecordStatusPending    = "pending"
	RecordStatusInProgress = "in_progress"
	RecordStatusOnHold     = "on_hold"
	RecordStatusFinished   = "finished"
	RecordStatusFailed     = "failed"
	RecordStatusCanceled   = "canceled"
)

var AllMaintenanceRecordStatuses = []string{
	RecordStatusPending,
	RecordStatusInProgress,
	RecordStatusOnHold,
	RecordStatusFinished,
	RecordStatusFailed,
	RecordStatusCanceled,
}
