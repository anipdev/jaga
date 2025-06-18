package consts

const (
	ScheduleTypePeriodic    = "periodic"
	ScheduleTypeConditional = "conditional"
)

var AllMaintenanceScheduleStatuses = []string{
	ScheduleTypePeriodic,
	ScheduleTypeConditional,
}
