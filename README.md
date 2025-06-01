# protogen

`protogen` allows you to generate TypeScript types and [zod](https://github.com/colinhacks/zod) validators from [protobuf](https://protobuf.dev/) schemas.

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

## Why?

Why build (yet) another type generator when more mature options like [ts-proto](https://github.com/stephenh/ts-proto) exist? Honestly, to learn Go. 

The other primary motivator was to build a better communication system for shypz.io, a multiplayer space game I've been working on. Of course, multiplayer browser games involve a significant amount of communication between the server and many clients to receive player inputs and broadcast updates like bullet collisions. Although I'm not totally certain if protobuf is the right tool for the job, it was clear I needed some kind of compressed structure to minimize the size of data transfers (and JSON isn't it).

I could have used typed arrays, but coming up with a custom approach that was fast and general enough to work with many different kinds of objects was hard. Even encoding/decoding seemingly basic primitives like strings would require a fair bit of effort and a lot of code duplication across the server and client. `protogen` was born to leverage `protobuf` in a TypeScript environment while bypassing most downsides of code duplication (since I would not need to maintain generated files, only a single tool producing those files). 

## Quickstart

