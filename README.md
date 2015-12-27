Server:
```
# Activate gopath
. ./gopath.sh

# Build protobuf files
. ./protobuild.sh

# Run it !
go run src/GRts/*.go

# Build it...
go build src/GRts/*.go
```

Client:
```
# Go to client directory
cd src/client/

# Install dependencies
npm install

# Babelify sources files:
grunt build

# To watch files
grunt
```

Voila: http://127.0.0.1:8080/index
