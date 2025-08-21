package generated

/*
	Place all go generate directives here in each line.
	Use relative imports for any dependent files.
	Go generate runs your tools with current working directory as the package directory.
*/

//go:generate go tool oapi-codegen -version
//go:generate go tool oapi-codegen -package genRouter -generate "std-http,strict-server,skip-prune" -o router/server.gen.go ../../api/openapi.yml
//go:generate go tool oapi-codegen -package genRouter -generate "types,skip-prune" -o router/types.gen.go ../../api/openapi.yml
//go:generate go tool oapi-codegen -package genRouter -generate "spec,skip-prune" -o router/spec.gen.go ../../api/openapi.yml
//go:generate go tool mockgen -version
//go:generate go tool mockgen -package mockdb -source ../db/db.go -destination mockdb/db_mock.go
