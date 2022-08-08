
# About

This file documents parts of the art process and tracks sources for reference photos.


## Tracing

- Find subject on <unsplash.com> and download high quality image
- Rotate, crop and align image using Gimp
- Fuzzy select (shrink, grow and feather edges) and cut out subject
- Save as reference

- Import reference to Inkscape
- Trace bitmap > single scan > Brightness cutoff and with settings:
  threshold=0.990, speckles/corners/optimize to max (to cut down nodes)
- Path > simplify done once or twice to cut down amount of nodes
- Set background fill and stroke
- This is the background outline

- Trace again > multicolor > colors with settings:
  scans ~ 6 (amount of color layers, find a good amount)
  smooth (decrease amount of nodes), stack (avoids making gaps between colors), remove background
  speckles/corners/optimize max again
- Path > simplify a couple of times until good
- Clean up trace by removing/editing paths (like too small details, overlap etc.)
- Image done


## GIF clips

- Create source video with peek:
  Record as mp4, to prevent lag
  Set framerate to 30, to keep the size down

- Convert mp4 to gif using ffmpeg:
  ffmpeg -an -c:v gif -i input.mp4 -t 10s out.gif
  (you propably want more than 10 seconds of output)

- Verify file size and playback


## Reference photo sources

- Clown fish by Rachel Hisko
  https://unsplash.com/photos/rEM3cK8F1pk

- Underwater by Sime Basioli
  https://unsplash.com/photos/BRkikoNP0KQ

