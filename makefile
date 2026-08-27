.PHONY: dev build run

dev-api:
	cd apps/api && go run cmd/main.go

dev-back:
	cd apps/backend && go run main.go

dev-interface:
	cd apps/web && pnpm run dev
