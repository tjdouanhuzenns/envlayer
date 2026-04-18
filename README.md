# envlayer

A CLI tool for layering and merging environment variable files across dev/staging/prod contexts.

---

## Installation

```bash
go install github.com/yourusername/envlayer@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/envlayer.git
cd envlayer
go build -o envlayer .
```

---

## Usage

Define your environment files in layers and merge them into a single output:

```bash
# Merge base + production overrides into a resolved .env file
envlayer merge --base .env.base --layer .env.prod --output .env

# Preview the merged result without writing to disk
envlayer merge --base .env.base --layer .env.staging --dry-run

# Specify a context directly
envlayer apply --context prod
```

**Example layer files:**

`.env.base`
```
APP_NAME=myapp
LOG_LEVEL=info
DB_PORT=5432
```

`.env.prod`
```
LOG_LEVEL=warn
DB_HOST=prod-db.example.com
```

**Merged output:**
```
APP_NAME=myapp
LOG_LEVEL=warn
DB_PORT=5432
DB_HOST=prod-db.example.com
```

Values defined in a layer file take precedence over the base. Multiple layers can be chained — later layers override earlier ones.

---

## License

[MIT](LICENSE)