package term

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"fmt"
	"image"
	"strings"
)

// kittyChunk is the largest payload kitty accepts in one escape sequence.
const kittyChunk = 4096

// EncodeKitty encodes an image for the kitty graphics protocol, placed at the
// cursor and composited BEHIND the text.
//
// Three flags carry the whole difference from Sixel:
//
//   - f=32 is 32-bit RGBA, so the alpha channel survives and a rounded corner
//     blends with whatever is actually under it rather than with a background
//     the caller had to guess.
//   - z=-1 puts the image behind the text layer. The cells drawn over it stay
//     real text: selectable, in the reader's font, and redrawn by the ordinary
//     frame without the image being touched.
//   - C=1 leaves the cursor alone, so placing a picture does not disturb the
//     serialization happening around it.
//
// The payload is zlib-compressed (o=z) because a panel is mostly flat color
// and RGBA is four bytes a pixel; on a slow link that ratio is the difference
// between a redraw and a stutter.
//
// # On Unicode placeholders
//
// kitty can also anchor an image to ordinary printable cells, which would make
// placement survive a line-diffing renderer for free. It needs a 297-entry
// diacritic table to address rows and columns, and it is the better long-term
// answer. This is direct placement: simpler, no table, and it is re-sent when
// the frame around it changes. See issue 29.
func EncodeKitty(img *image.RGBA, id, cols, rows int) string {
	return encodeKitty(img, id, cols, rows, 2)
}

// EncodeKittyVerbose is the same bytes with replies turned back on (q=0).
//
// For diagnosis only. q=2 exists because a reply arriving mid-frame is read as
// a keystroke — but it also means a terminal REFUSING an image says so to
// nobody, and "the picture did not appear" is then indistinguishable from "the
// picture appeared behind an opaque background". One of those is a bug in this
// encoder and the other is not.
func EncodeKittyVerbose(img *image.RGBA, id, cols, rows int) string {
	return encodeKitty(img, id, cols, rows, 0)
}

func encodeKitty(img *image.RGBA, id, cols, rows, quiet int) string {
	b := img.Bounds()
	w, h := b.Dx(), b.Dy()
	if w <= 0 || h <= 0 {
		return ""
	}

	// Rows may be padded, so the pixel data is copied row by row rather than
	// handed over as img.Pix — a stride wider than the image would shear it.
	raw := make([]byte, 0, w*h*4)
	for y := b.Min.Y; y < b.Max.Y; y++ {
		off := img.PixOffset(b.Min.X, y)
		raw = append(raw, img.Pix[off:off+w*4]...)
	}

	var zbuf bytes.Buffer
	zw := zlib.NewWriter(&zbuf)
	if _, err := zw.Write(raw); err != nil {
		return ""
	}
	if err := zw.Close(); err != nil {
		return ""
	}
	payload := base64.StdEncoding.EncodeToString(zbuf.Bytes())

	var out strings.Builder
	first := true
	for len(payload) > 0 {
		n := min(kittyChunk, len(payload))
		chunk := payload[:n]
		payload = payload[n:]

		out.WriteString("\x1b_G")
		if first {
			// a=T transmits and displays in one step. q=2 suppresses both the
			// success and the failure reply: a terminal answering into the
			// input stream mid-frame would arrive as a keystroke.
			// c and r say the footprint in CELLS, and the terminal scales the
			// image into exactly that rectangle.
			//
			// Without them the terminal divides the pixel size by its own idea
			// of a cell — which is the one number nobody agrees on. Ghostty
			// reports the cell in device pixels, iTerm2 in points, and a
			// picture sized from the wrong one comes out at the wrong scale and
			// overlapping the text above it. Stating the footprint removes the
			// measurement from the answer entirely: it is the caller who knows
			// how many cells the picture is FOR, and that number is exact.
			fmt.Fprintf(&out, "a=T,f=32,o=z,s=%d,v=%d,c=%d,r=%d,i=%d,z=-1,C=1,q=%d,",
				w, h, cols, rows, id, quiet)
			first = false
		}
		// m=1 means another chunk follows; m=0 is the last one.
		//
		// The separating comma belongs to the keys BEFORE it, not to m. Written
		// as ",m=1" it produced `ESC_G,m=1;` on every chunk after the first —
		// a leading comma, which is malformed, so a terminal took the first
		// chunk and refused the rest. Small images have one chunk and worked
		// perfectly; only images past 4096 bytes of payload vanished, which is
		// a bug that hides itself until the picture gets interesting.
		if len(payload) > 0 {
			out.WriteString("m=1")
		} else {
			out.WriteString("m=0")
		}
		out.WriteByte(';')
		out.WriteString(chunk)
		out.WriteString("\x1b\\")
	}
	return out.String()
}

// DeleteKitty removes a placed image by id.
//
// Needed because a kitty image outlives the frame that placed it. Nothing else
// in tuikit has a lifetime longer than one canvas, which is exactly why this
// has to be explicit: a picture is the one thing that stays behind.
func DeleteKitty(id int) string { return fmt.Sprintf("\x1b_Ga=d,d=i,i=%d,q=2\x1b\\", id) }
