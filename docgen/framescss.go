package docgen

// pageCSS is the chrome around the frames, and it is deliberately quiet.
//
// Every colour is defined on bare :root and only overridden in the theme
// blocks. A colour whose one definition sits inside a media query never applies
// when the viewer's theme is unset, which is the default — and the page then
// renders one theme's text on the other theme's ground.
//
// The .tuikit-frame rule is the exception that proves the rest: it takes no
// tokens and does not change between themes, because the ANSI inside it was
// captured for a dark terminal and recolouring it would report colours the tool
// does not have.
const pageCSS = `
:root {
  --ground:  #fbfafb;
  --surface: #ffffff;
  --ink:     #22202a;
  --muted:   #63606f;
  --faint:   #8a8695;
  --line:    #e6e3ec;
  --accent:  #5b4bd6;
  --live:    #1d7a4c;
  --composed:#9a5b12;
  --sans: ui-sans-serif, -apple-system, "Segoe UI", system-ui, sans-serif;
  --mono: ui-monospace, "SF Mono", SFMono-Regular, Menlo, Consolas, "DejaVu Sans Mono", monospace;
}
@media (prefers-color-scheme: dark) {
  :root:not([data-theme="light"]) {
    --ground:  #131218;
    --surface: #1b1a22;
    --ink:     #e8e6f0;
    --muted:   #a29fb0;
    --faint:   #7b7889;
    --line:    #2b2936;
    --accent:  #a99cff;
    --live:    #5dd39a;
    --composed:#e0a95c;
  }
}
:root[data-theme="dark"] {
  --ground:  #131218;
  --surface: #1b1a22;
  --ink:     #e8e6f0;
  --muted:   #a29fb0;
  --faint:   #7b7889;
  --line:    #2b2936;
  --accent:  #a99cff;
  --live:    #5dd39a;
  --composed:#e0a95c;
}

* { box-sizing: border-box; }
body {
  margin: 0; padding: clamp(24px, 4vw, 56px) clamp(16px, 4vw, 48px) 80px;
  background: var(--ground); color: var(--ink);
  font: 16px/1.6 var(--sans);
}
header, section, figure { max-width: 1180px; margin-inline: auto; }
h1 { font-size: clamp(26px, 4vw, 38px); line-height: 1.08; letter-spacing: -0.02em;
     text-wrap: balance; margin: 0 0 10px; }
.lede { color: var(--muted); max-width: 64ch; margin: 0 0 12px; font-size: 17px; }
.facts { color: var(--faint); font-size: 13px; margin: 0 0 8px;
         font-variant-numeric: tabular-nums; }

section { margin-top: 56px; }
.group { border-top: 2px solid var(--ink); padding-top: 12px; margin-bottom: 28px; }
.group h2 { font-size: 20px; letter-spacing: -0.01em; margin: 0 0 3px; }
.group p { margin: 0; color: var(--muted); font-size: 15px; }

figure { margin: 0 0 40px; }
figcaption { margin-bottom: 10px; }
.id { display: flex; align-items: baseline; gap: 12px; flex-wrap: wrap; }
.id h3 { font-family: var(--mono); font-size: 14px; font-weight: 500; margin: 0; }
.keys { display: inline-flex; gap: 4px; }
kbd { font-family: var(--mono); font-size: 11px; line-height: 1; padding: 4px 6px;
      border: 1px solid var(--line); border-bottom-width: 2px; border-radius: 4px;
      background: var(--surface); color: var(--muted); }
figcaption > p { margin: 7px 0 0; color: var(--muted); max-width: 74ch; font-size: 15px; }
.dims { display: flex; align-items: center; gap: 9px; color: var(--faint);
        font-size: 12px; font-variant-numeric: tabular-nums; }
.chip { font-size: 10px; letter-spacing: 0.09em; text-transform: uppercase; font-weight: 600;
        padding: 3px 7px; border-radius: 999px;
        border: 1px solid currentColor; color: var(--accent); }
.chip.live { color: var(--live); }
.chip.composed { color: var(--composed); }

/* The frames keep a dark ground in BOTH themes, deliberately. See the comment
   above: the ANSI was captured for a dark terminal.

   The font stack is written out rather than taken from --mono for the same
   reason the ground is: a token can be redefined, and a webfont with a fallback
   for the glyphs it lacks gives box-drawing and braille different advance
   widths, which takes a frame apart. */
.caveat {
  margin: 10px 0 0; padding: 9px 12px;
  border-left: 2px solid var(--accent);
  background: var(--surface); color: var(--muted);
  font-size: 13px; line-height: 1.5; max-width: 62ch;
}
.caveat strong { color: var(--ink); font-weight: 600; }
.tuikit-frame {
  font-family: ui-monospace, "SF Mono", SFMono-Regular, Menlo, Consolas, "DejaVu Sans Mono", monospace;
  font-size: 12.5px; line-height: 1.25;
  font-variant-ligatures: none; font-feature-settings: "liga" 0, "calt" 0;
  white-space: pre; tab-size: 4;
  background: #0c0c0e; color: #d0d0d0;
  padding: 15px 17px; border-radius: 7px;
  overflow-x: auto; border: 1px solid #26262c;
}
@media (max-width: 640px) { .tuikit-frame { font-size: 10px; } }
`
