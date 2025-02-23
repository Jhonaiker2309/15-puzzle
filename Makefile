BINARY_NAME= solver
GO=go

.PHONY: build clean run

build:
	@echo "Compilando todos los archivos .go..."
	$(GO) build -o $(BINARY_NAME) *.go

run:
	@echo "Ejecutando la aplicación..."
	./$(BINARY_NAME)

extra:
	@echo "Ejecutando la aplicación..."
	./$(BINARY_NAME) -extra_heuristic

clean:
	@echo "Limpiando..."
	rm -f $(BINARY_NAME)
	rm matrix_states.json