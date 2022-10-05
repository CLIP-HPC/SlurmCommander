package openapi

// TODO: come up with a workflow to maintain and (auto)-import currated list of
// openapi_version.json files and/or openapi.gen.go for different slurm versions
//

//xxxgo:generate oapi-codegen --old-config-style --package=openapi --generate=types -alias-types -o ./openapi.gen.go ./openapi_0_0_39.json

// debug 2022/10/03 22:20:02 Error unmarshall: "json: cannot unmarshal object into Go struct field V0038JobResources.Jobs.job_resources.allocated_nodes of type []openapi.V0038NodeAllocation"
// v38 and 39 work, when commented out line 483 from openapi.gen.go:
//
// array of allocated nodes
//AllocatedNodes *[]V0038NodeAllocation `json:"allocated_nodes,omitempty"`
