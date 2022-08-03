
# Tracing

- Find subject on <unsplash.com> and download high quality image
- Rotate, crop and align image using Gimp/Krita
- Fuzzy select (shrink, grow and feather edges) and cut out subject
- Save as reference

- Import reference to Inkscape
- Trace bitmap > single scan > Brightness cutoff and with settings:
  threshold=0.990, speckles/corners/optimize to max (to cut down nodes)
- Path > simplify done once or twice to cut down amount of nodes
- Set background fill and stroke
- Background outline is done

- Trace again > multicolor > colors with settings:
  scans ~ 6 (amount of color layers, find a good amount)
  smooth (decrease amount of nodes), stack (avoids making gaps between colors), remove background
  speckles/corners/optimize max again
- Path > simplify a couple of times until good
- Clean up trace by removing/editing paths (like too small details, overlap etc.)
- Image done
