# Distributed Commit Log

Distributed Commit Log implementation using Golang.

## Terminology

- Record—the data stored in our log.
- Store—the file we store records in.
- Index—the file we store index entries in.
- Segment—the abstraction that ties a store and an index together.
- Log—the abstraction that ties all the segments together.

## Features

- Create data classes using protobuf
- Authenticate client/server using certificates
- Authorize with access control lists

Todo: complete section

## Libraries

- protoc (data models)
- gRPC (communication)
- cfssl (generating certs)
- crypto (encryption)
- casbin (access control list)
- testify (testing module)

## Installation

Golang v1.15+ is recommended.
Todo: complete section