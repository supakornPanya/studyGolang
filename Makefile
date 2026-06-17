# Create a new empty migration file (for manual SQL, seeding, etc.)
# Usage: make migrate-new-file name=add_seed_data
migrate-new-file:
	atlas migrate new $(name) --env gorm

# Run Atlas diff to generate a new goose migration
# Usage: make migrate-diff name=create_users_table
migrate-diff:
	atlas migrate diff $(name) --env gorm

# Check migration status in gorm
migrate-status:
	atlas migrate status --env gorm

# Apply new migrations in gorm
migrate-apply:
	atlas migrate apply --env gorm

# Rollback/revert the last migration in gorm
migrate-down:
	atlas migrate down --env gorm

# inspect schema
inspect-schema:
	atlas schema inspect -u "env://src" --env gorm

# re-hash migration files 
# - run at first/when edit migration files for create history migration file
migrate-hash:
	atlas migrate hash --env gorm