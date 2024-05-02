unit-test:
	gotestsum -- ./... -failfast -race -coverprofile ./coverage.out

mocks:
	mockery --all --keeptree --inpackage