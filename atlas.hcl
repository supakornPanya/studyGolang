data "external_schema" "gorm" {
  program = [
    "go",
    "run",
    "-mod=mod",
    "./cmd/atlas-loader/main.go"
  ]
}

env "gorm" {
  src = data.external_schema.gorm.url
  url = "postgres://postgres:mysecretpassword@localhost:5432/postgres?sslmode=disable"
  dev = "postgres://postgres:mysecretpassword@localhost:5432/postgres_dev?sslmode=disable"
  
  migration {
    dir = "file://internal/db/migrations?format=goose"
  }
}