package agent

import (
"encoding/json"
)

// DynamicUISchema define a estrutura gerada sob demanda.
type DynamicUISchema struct {
Component string                 `json:"component"`
Props     map[string]interface{} `json:"props"`
Layout    string                 `json:"layout"`
}

// GenerateUISchema analisa o contexto e decide o layout do dashboard.
func (s *MasterTeamService) GenerateUISchema(taskType string) string {
var schema DynamicUISchema

switch taskType {
case "code_debug":
namicUISchema{
ent: "CodeAnalyzer",
g]interface{}{
tax-errors",
"split-screen",
":
namicUISchema{
ent: "ReactFlowGraph",
g]interface{}{
imated": true,
"full-canvas",
namicUISchema{
ent: "GeneralStatus",
g]interface{}{"status": "idle"},
dard",
json.Marshal(schema)
return string(res)
}
