package consts

const (
	AssetStatusReady            = "ready"
	AssetStatusUnderMaintenance = "under_maintenance"
	AssetStatusNeedMaintenance  = "need_maintenance"
)

var AllAssetStatuses = []string{
	AssetStatusReady,
	AssetStatusUnderMaintenance,
	AssetStatusNeedMaintenance,
}
