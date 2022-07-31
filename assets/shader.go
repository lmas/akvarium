package main

// Underwater shader.
// Writing shaders for ebiten is a little bit special, because of it's custom shader
// language called `Kage`, that looks like go and only allows for fragment shaders.
// https://ebiten.org/documents/shader.html

var Time float
var Resolution vec2

// Source: https://www.shadertoy.com/view/MdXGW7
// TODO: doesn't loop at all (https://news.ycombinator.com/item?id=30438541 for example)
// TODO: it's also got bubbles. I want some random bubbles..
func sunRay(coord, raySource, rayDirection vec2, seedA, seedB, speed float) vec4 {
	sourceToCoord := coord - raySource
	cosAngle := dot(normalize(sourceToCoord), rayDirection)
	val := (0.45 + 0.15*sin(cosAngle*seedA+Time*speed)) + (0.3 + 0.2*cos(-cosAngle*seedB+Time*speed))
	strength := (Resolution.x - length(sourceToCoord)) / Resolution.x
	return vec4(1.0, 1.0, 1.0, 1.0) * clamp(val, 0.0, 1.0) * clamp(strength, 0.5, 1.0)
}

func Fragment(position vec4, texCoord vec2, color vec4) vec4 {
	var fragColor vec4

	// Lotsa smaller rays, slowly rotating counter-clockwise?
	fragColor += sunRay(
		texCoord, vec2(Resolution.x*0.7, Resolution.y*-0.4), normalize(vec2(1.0, 0.2843)), 15.1869, 29.5428, 1.1,
	) * 0.4

	// A few larger-ish rays, moving faster clockwise?
	fragColor += sunRay(
		texCoord, vec2(Resolution.x*0.8, Resolution.y*-0.6), normalize(vec2(1.0, -0.0596)), 21.4852, 17.9246, 1.5,
	) * 0.5

	// Emulate light attenuation towards the depths, for the sun rays.
	// https://en.wikipedia.org/wiki/Attenuation
	fragColor *= (1 - smoothstep(0, Resolution.y, texCoord.y)) * 0.7

	// Apply smooth darkness towards the depths for whole screen
	// https://en.wikipedia.org/wiki/Smoothstep
	// TODO: could add a little "waving" to the bottom?
	fragColor += vec4(0, 0, 0, 1) * smoothstep(0, Resolution.y, texCoord.y)

	return fragColor
}
