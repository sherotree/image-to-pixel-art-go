#!/usr/bin/env python3
import base64, pathlib, zlib
root = pathlib.Path(__file__).resolve().parents[1]
parts = sorted(root.glob("testdata/parity.part*.b64"), key=lambda p: p.name)
assert parts, "no parts"
b64 = "".join(p.read_text().strip() for p in parts)
data = zlib.decompress(base64.b64decode(b64))
out = root / "testdata" / "parity.json"
out.write_bytes(data)
print(f"wrote {out} ({len(data)} bytes) from {len(parts)} parts")
