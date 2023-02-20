# `proto-proxyman`

`proto-proxy` reads a `google.protobuf.FileDescriptorSet` from STDIN and
outputs the [Proxyman](https://docs.proxyman.io/advanced-features/protobuf)
configuration to STDOUT, allowing you to easily proxy protobuf-encoded requests
and responses.

You can then import the generated config file into Proxyman via `Tools >
Protobuf > More > Import Settings...`.

## `google.protobuf.FileDescriptorSet`

A `.desc` file is an encoded `google.protobuf.FileDescriptorSet` message.

If you're using `buf` you can generate one with `buf build
--as-file-descriptor-set -o -`.

If you're using `protoc` the `--descriptor_set_in` and `--descriptor_set_out`
flags will be what you're looking for!
