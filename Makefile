BINARY_NAME= search
GO=go

.PHONY: build clean run

build:
	@echo "Compilando todos los archivos .go..."
	$(GO) build -o $(BINARY_NAME) *.go

run:
	@echo "Ejecutando la aplicación..."
	./$(BINARY_NAME)

clean:
	@echo "Limpiando..."
	rm -f $(BINARY_NAME)
	rm matrix_states.json