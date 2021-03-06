path = bin/splitter-mac

build:
	cd src && GOOS=darwin GOARCH=amd64 go build -o ../$(path) main.go

install: build
	mv $(path) /usr/local/bin/splitter

install-libgit2:
	cd /tmp && \
	curl -sL https://github.com/libgit2/libgit2/archive/refs/tags/v1.3.0.tar.gz -o libgit2-1.3.0.tar.gz && \
	tar xzf libgit2-1.3.0.tar.gz && \
	cd libgit2-1.3.0/build && \
	cmake .. -DCMAKE_INSTALL_PREFIX=/usr/local/opt && \
	sudo cmake --build . --target install && \
	cd /tmp && rm -rf libgit2*

download:
	curl -sL https://github.com/myposter-de/monorepo-splitter/releases/download/0.0.3/splitter-mac -o /usr/local/bin/splitter
