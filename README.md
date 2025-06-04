# protogen

`protogen` allows you to generate TypeScript types and [zod](https://github.com/colinhacks/zod) validators from [protobuf](https://protobuf.dev/) schemas.

## Quickstart

## Why?

Why build (yet) another type generator when more mature options like [ts-proto](https://github.com/stephenh/ts-proto) exist? 

Honestly, I wanted to learn Go and a CLI tool seems like a great starter project. 

The other primary motivator was to build a typed communication system for shypz.io, a multiplayer space game I've been working on. Given that multiplayer games involve a significant amount of messages between the server and client (to receive player inputs and broadcast physics updates, for example), it was clear to me that I needed a compressed structure to minimize the size of data transfers. My first thought was to use typed arrays. However, developing a protocol that was fast, efficient, and general enough to work with many different data types increasingly felt impractical. Even encoding/decoding seemingly basic primitives like variable length strings didn't feel right (or fast), and my solution involved a lot of code duplication across the server and client to properly receive and send events. Ultimately, I decided to use `protobuf` and make some kind of type generator. I'm not totally certain if this is even the right tool for the job, but I know JSON certainly isn't. 

Thus, `protogen` was born out of an excuse to learn a new programming langauge and leverage `protobuf` in a TypeScript environment. 

## Supported Syntax

|Syntax|Implemented|
|-------|-----------|
|Message|✅|
|Reserved|✅|
|Comments|❌|
|Enums|❌|
|Import|❌|
|`repeated`|❌|
|`optional`|❌|
|Nested messages|❌|
|`google.protobuf.Any`|❌|
|`oneof`|❌|
|Maps|❌|
|`option`|❌|
|`extend`|❌|
