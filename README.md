Requirements:
    - go
    - npm
    - grunt
    - protobuf (https://github.com/google/protobuf)

You will need to create a symbolic link to client sources in the root directory
```
    ln -s ./src/GRts/client/ ./
```

Server:
```
# Activate gopath
. ./gopath.sh

# Build protobuf files
. ./protobuild.sh

# Run it !
go run src/GRts/server/*.go

# Build it...
go build src/GRts/server/*.go
```

Client:
```
# Client directory
cd src/GRts/client/

# Install dependencies
npm install

# Babelify sources files:
grunt build

# To watch files
grunt
```

Voila: http://127.0.0.1:8080/index
