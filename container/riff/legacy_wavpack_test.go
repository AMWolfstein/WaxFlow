package riff

import (
	"archive/zip"
	"bytes"
	"io"
	"testing"

	"github.com/colespringer/waxflow/audio"
	"github.com/colespringer/waxflow/codec/pcm"
	"github.com/colespringer/waxflow/container"
	"github.com/colespringer/waxflow/internal/testutil"
	"github.com/colespringer/waxflow/waxerr"
)

func TestLegacyWavPackUnsupported(t *testing.T) {
	zr, err := zip.OpenReader(testutil.VectorPath(t, "wavpack/test_suite.zip"))
	if err != nil {
		t.Fatal(err)
	}
	defer zr.Close()
	names := []string{
		"vers-10.wv", "vers-20.wv", "vers-20-lossy.wv",
		"vers-30.wv", "vers-30-fast.wv", "vers-30-lossy.wv",
		"vers-397.wv", "vers-397-hybrid.wv",
	}
	for _, name := range names {
		t.Run(name, func(t *testing.T) {
			rc, err := zr.Open("test_suite/legacy/" + name)
			if err != nil {
				t.Fatal(err)
			}
			raw, err := io.ReadAll(rc)
			rc.Close()
			if err != nil {
				t.Fatal(err)
			}
			for _, strict := range []bool{false, true} {
				_, err := NewDemuxer(container.BytesSource(raw), &DemuxerOptions{Strict: strict})
				if waxerr.CodeOf(err) != waxerr.CodeUnsupportedFormat {
					t.Fatalf("strict=%v: %v, want unsupported", strict, err)
				}
				const want = "wavpack: unsupported stream version: only version-4 streams (0x402..0x410) are supported"
				if err.Error() != want {
					t.Fatalf("strict=%v: %q, want %q", strict, err, want)
				}
			}
		})
	}
}

func TestPCMWavPackMagicAloneIsNotLegacy(t *testing.T) {
	// Ordinary PCM may spell the magic. Require the legacy size and version too.
	wire := append([]byte("wvpk"), make([]byte, 60)...)
	var wav bytes.Buffer
	cfg := pcm.Config{Encoding: pcm.SignedInt, Bits: 16}
	f := audio.Format{Rate: 44100, Channels: 2, Layout: audio.DefaultLayout(2), Type: audio.Int, BitDepth: 16}
	muxWAV(t, &wav, cfg, f, wire, 16, nil)
	d, err := NewDemuxer(container.BytesSource(wav.Bytes()), nil)
	if err != nil {
		t.Fatal(err)
	}
	var pkt container.Packet
	if err := d.ReadPacket(&pkt); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(pkt.Data, wire) {
		t.Fatal("PCM payload changed")
	}
}
