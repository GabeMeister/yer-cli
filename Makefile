# Watch the templ files, and re-generate go code if they change
live/templ:
	templ generate --watch --proxy="http://localhost:4000" --open-browser=false -v

# Run the golang web server, reload if any regular golang code changes
live/server:
	air

# Watch for css changes in any *.templ file
live/tailwind:
	go run github.com/cosmtrek/air@v1.51.0 \
	--build.cmd "npx tailwindcss --input ./presentation/styles/globals.css --output ./presentation/static/css/styles.css" \
	--build.bin "true" \
	--build.delay "100" \
	--build.include_ext "templ,globals.css"

# Watch for any changes to the styles.css file, and notify the proxy that
# styling has changed
live/sync_assets:
	go run github.com/cosmtrek/air@v1.51.0 \
	--build.cmd "templ generate --notify-proxy" \
	--build.bin "true" \
	--build.delay "100" \
	--build.include_dir "presentation/static/css" \
	--build.include_ext "css"

# Start all 5 watch processes in parallel.
live: 
	make -j4 live/templ live/server live/tailwind live/sync_assets