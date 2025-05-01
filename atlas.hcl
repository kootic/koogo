data "external_schema" "bun" {
  program = [
    "go",
    "run",
    "-mod=mod",
    "ariga.io/atlas-provider-bun",
    "load",
    "--path", "./internal/repo/postgres/bun",
    "--dialect", "postgres",
  ]
}

env "bun" {
  src = data.external_schema.bun.url
  dev = "postgres://postgres:postgres@localhost:5432/atlasdev?sslmode=disable"
  migration {
    dir = "file://internal/repo/postgres/migrations"
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}
