# protogen

`protogen` allows you to generate TypeScript types and [zod](https://github.com/colinhacks/zod) validators from [protobuf](https://protobuf.dev/) schemas.

## Quickstart

## Why?

Why build (yet) another type generator when more mature options like [ts-proto](https://github.com/stephenh/ts-proto) exist?

I wanted to learn Go and a CLI seems like a great starter project. Also I'm obsessed with [generators](https://github.com/natmfat/shitgen).

In terms of a real-world use case, I needed a typed communication system for a multiplayer space game I've been working on. Multiplayer games involve a lot of communication between the server and client (for player inputs and broadcasting physics updates, for example). From the start it was clear that I needed a compressed structure. Typed arrays were okay, but weren't general enough to work with many data types - encoding and decoding seemingly basic primitives (like variable length strings) involved a lot of code duplication and felt wrong.

Ultimately, I decided to use `protobuf` and build some kind of type generator. I'm not totally certain if this is even the right tool for the job, but I know raw JSON certainly isn't.

## Supported Syntax

`protogen` only supports a subset of the [proto3 language](https://protobuf.dev/programming-guides/proto3/). I've covered the basics, but anything more complicated should use mature and battle tested library instead.

| Syntax                | Implemented |
| --------------------- | ----------- |
| Message               | ⚠️ (wip)    |
| Reserved              | ✅          |
| Comments              | ✅          |
| Enums                 | ✅          |
| Import                | ❌          |
| Maps                  | ❌          |
| `repeated`            | ❌          |
| `optional`            | ❌          |
| `google.protobuf.Any` | ❌          |
| `oneof`               | ❌          |
| `option`              | ❌          |
| `extend`              | ❌          |
