package openapidb

// TODO: https://github.com/deepmap/oapi-codegen/issues/542
//go:generate oapi-codegen --old-config-style --package=openapidb --generate=types -alias-types -o ./openapi_db.gen.go ./openapi_0_0_39_db.json

// v39 bugs

// line: 3181 has typo:
// has:
// #/components/schemas/dbv0.0.36_tres_list
// should have:
// #/components/schemas/dbv0.0.39_tres_list

// line: 365
// same as with openapi pkg, comment out AllocationNodes

// line: 661
// can't unmarshall Task
