# Task Completion Checklist

## Before Committing Code

### Code Quality Checks
- [ ] Run `go fmt ./...` to format all Go files
- [ ] Run `go vet ./...` for static analysis
- [ ] Run `go mod tidy` to clean dependencies
- [ ] Ensure no unused imports or variables

### Testing Requirements
- [ ] Run `go test ./...` and ensure all tests pass
- [ ] Run `go test -race ./...` to check for race conditions
- [ ] Add tests for new functionality
- [ ] Update existing tests if behavior changes
- [ ] Ensure test coverage for critical paths

### Mock Generation
- [ ] Run `mockery --all` if interfaces were added/modified
- [ ] Verify generated mocks are properly committed

### Database Changes
- [ ] Create appropriate migrations in `migration/` directory
- [ ] Test migrations up and down
- [ ] Update seed data if necessary
- [ ] Verify database schema changes work correctly

### Documentation
- [ ] Update function/method comments for public APIs
- [ ] Update README.md if new features affect usage
- [ ] Document any new configuration options

### Build Verification
- [ ] Verify `go build .` succeeds without errors
- [ ] Test the built binary runs correctly
- [ ] Verify all dependencies are properly declared in go.mod

### Environment Configuration
- [ ] Update `.env.example` if new configuration variables added
- [ ] Verify application works with default configuration