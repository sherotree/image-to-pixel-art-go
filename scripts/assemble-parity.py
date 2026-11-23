#!/usr/bin/env python3
"""Assemble testdata/parity.json from zlib+base64 sidecar."""
import base64, pathlib, zlib
root = pathlib.Path(__file__).resolve().parents[1]
src = root / "testdata" / "parity.json.zlib.b64"
dst = root / "testdata" / "parity.json"
dst.write_bytes(zlib.decompress(base64.b64decode(src.read_text().strip())))
print(f"wrote {dst} ({dst.stat().st_size} bytes)")
