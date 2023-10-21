go build -o raceweb cmd/rest/raceweb.go
go build -o placer cmd/default_placer/main.go
go build -o result-builder cmd/result_builder/main.go
go build -o cli cmd/manual_cli/main.go
go build -o overall-scorer cmd/scorer/main.go
go build -o generator cmd/event_generator/main.go
