# Generate Go types from all schema files in schema/
generate-types:
    #!/usr/bin/env bash
    set -euo pipefail
    for schema in schema/*.json; do
        version=$(basename "$schema" .json)
        pkg="v${version//[.-]/_}"
        dir="internal/$pkg"
        mkdir -p "$dir"
        echo "Generating $dir/types_gen.go from $schema"
        go-jsonschema -p "$pkg" \
            --tags json \
            --schema-root-type 'OpenCLI.json=Document' \
            -o "$dir/types_gen.go" \
            "$schema"
    done
