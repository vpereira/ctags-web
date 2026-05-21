# ctags-web

Small web UI for browsing a `universal-ctags` index stored in FerretDB.

It has three pieces:

- `index/` reads `ctags.json` and stores tags
- `import/` reads source files and stores code lines
- `web/` serves search and file views on `:8080`

The UI is simple:

- `/token?token=main` searches tags
- `/show?file=/path/to/file.go&linecount=33` shows a file

## Run with Docker

Start everything:

```bash
docker compose up -d
```

Open:

```text
http://localhost:8080/
```

The database data is stored under `data/postgres`.

If you want a clean database:

```bash
docker compose down
rm -rf data/postgres
docker compose up -d
```

## Build locally

```bash
make build
```

That builds:

- `index/ctags-index`
- `import/ctags-import`
- `web/ctags-web`

## Index a project

Generate the ctags JSON:

```bash
./scripts/run_ctags.sh .
```

Load tags:

```bash
./index/ctags-index "mongodb://test:test@localhost:27017/" ctags.json
```

Load source lines:

```bash
./import/ctags-import "mongodb://test:test@localhost:27017/" ctags code .
```

Run the web server locally:

```bash
./web/ctags-web "mongodb://test:test@localhost:27017/" ctags ctags
```

If Docker is already running, you can also do the import inside the `web` container:

```bash
docker compose exec web sh -lc '
ctags --recurse=yes --fields=* --output-format=json -f /tmp/ctags.json /app &&
/app/index "mongodb://test:test@ferretdb:27017/" /tmp/ctags.json &&
/app/import "mongodb://test:test@ferretdb:27017/" ctags code /app &&
rm -f /tmp/ctags.json
'
```
