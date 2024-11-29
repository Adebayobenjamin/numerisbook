.PHONY: mocks
## mocks: generate mock files for interfaces
mocks:
	@echo "Generating mocks..."
	@# Repository mocks
	@for file in pkg/repositories/interfaces/*.interface.go; do \
		base_name=$$(basename $$file .interface.go); \
		mockgen -source=$$file -destination=pkg/repositories/mocks/mock_$${base_name}.go -package=repository_mocks; \
	done
	@# Service mocks
	@for file in pkg/services/interfaces/*.interface.go; do \
		base_name=$$(basename $$file .interface.go); \
		mockgen -source=$$file -destination=pkg/services/mocks/mock_$${base_name}.go -package=services_mocks; \
	done
	@echo "Mocks generated successfully!"
	