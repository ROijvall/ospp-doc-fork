GOPATH := ./backend/src/
PYPATH := /frontend/src
GOBUILD = go build
GOCLEAN	= go clean
GOTEST = go clean

run: build
	./gameserver
	gnome-terminal -e cd frontend/src/ || python3 graphics.py

build: 
	$(GOBUILD) src/super-tanks/backend/src/ -o gameserver

client:
	python3 $(PYPATH)/graphics.py

test:
	go test ./backend/src/tank_test.go

clean-py:
	find . -name '*.pyc' -exec rm --force {} +
	find . -name '*.pyo' -exec rm --force {} +
	name '*~' -exec rm --force  {}

clean-go:
	rm -f *.exe

clean: clean-py clean-go

.PHONY: clean
