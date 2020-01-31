meta METAR
==========

Highlight critical conditions in METARs for VFR pilots. This tool does not substitute an official pre-flight weather briefing.

Critical Conditions
-------------------

1. Visibility below 6000 m (as 5000 m is legal minimum. The weather might drop below at any time)
1. Convective Clouds (CB and TCB, as they can indicate a risk for thunderstorms and wind shears)
1. Critical weather conditions with:
  SH Showers
  TS Thunderstorm
  GR Hail
  GS Small Hail
  BR Mist
  DU Dust
  FG Fog
  FU Smoke
  HZ Haze
  PY Spray
  SA Sand
  VA Volcanic Ash
1. Clouds below 1500 feet (legal limits in most areas)
1. Windspeeds above 12kt (crosswinds of more than 12kt are challenging to fly with small airplanes)

Usage
-----

To get live METAR data, an account and API key from https://www.checkwx.com/ is required.

Use the environment variable `M_API_KEY` for the checkwx API key.
