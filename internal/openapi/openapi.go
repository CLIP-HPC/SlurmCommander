package openapi

// TODO: come up with a workflow to maintain and (auto)-import currated list of
// openapi_version.json files and/or openapi.gen.go for different slurm versions
//

//go:generate oapi-codegen --old-config-style --package=openapi --generate=types -alias-types -o ./openapi.gen.go ./openapi_0.0.39_master.json
