package openapidb

// TODO: https://github.com/deepmap/oapi-codegen/issues/542
//go:generate oapi-codegen --old-config-style --package=openapidb --generate=types -alias-types -o ./openapi_db.gen.go ./openapi_0.0.37_21.08.8.json

// debug 2022/10/13 19:01:02 Error unmarshall: "json: cannot unmarshal number into Go struct field Dbv0037Job.Jobs.allocation_nodes of type string"
// debug 2022/10/13 19:04:06 Error unmarshall: "json: cannot unmarshal number into Go struct field .Jobs.het.job_id of type map[string]interface {}"

// debug 2022/10/13 19:08:37 Error unmarshall: "json: cannot unmarshal number into Go struct field .Jobs.steps.step.id of type string"
