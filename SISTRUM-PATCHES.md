# Sistrum fork patches

## Refuse RIFF-wrapped legacy WavPack — fix/riff-wrapped-wavpack

- Commit: the commit introducing this entry (`git log -1 --format=%H -- SISTRUM-PATCHES.md`).
- What: inspect the first ten bytes of a RIFF data payload and refuse the
  official decoder suite's version-1/2/3 WavPack headers as unsupported.
  Probe before size-clamping warnings, so strict and tolerant mode agree.
- Why: the old streams retain a PCM WAV header, but their compressed data
  must never be interpreted as PCM or used to generate PCM oracle hashes.
- Clean-room: yes. Original recognition logic written from the eight official
  suite files and WaxFlow's MIT RIFF parser. No FFmpeg or other GPL/LGPL
  implementation was consulted. Existing tests may invoke oracle binaries.
- Tests: direct RIFF regression for all eight files in both modes, engine-level
  refusal regression, PCM magic-only control, and exact bounded read-cost test.
- Upstream: not sent; suitable for upstream submission.
