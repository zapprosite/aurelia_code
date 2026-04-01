#!/usr/bin/env python3
"""
Edge TTS for Aurélia - Free PT-BR voice synthesis
Uses Microsoft Edge TTS: natural voices, no API key needed
"""

import asyncio
import argparse
import sys
from pathlib import Path


async def synthesize(
    text: str, voice: str, output: str, rate: str = "+0%", pitch: str = "+0Hz"
):
    from edge_tts import Communicate

    cm = Communicate(text, voice, rate=rate, pitch=pitch)
    await cm.save(output)
    print(f"Saved: {output}")


def main():
    parser = argparse.ArgumentParser(description="Edge TTS for Aurélia")
    parser.add_argument("--text", "-t", required=True, help="Text to synthesize")
    parser.add_argument(
        "--voice", "-v", default="pt-BR-FranciscaNeural", help="Voice name"
    )
    parser.add_argument("--output", "-o", required=True, help="Output file (mp3)")
    parser.add_argument("--rate", default="+0%", help="Speech rate (e.g., +10%, -5%)")
    parser.add_argument("--pitch", default="+0Hz", help="Pitch adjustment")

    args = parser.parse_args()
    asyncio.run(synthesize(args.text, args.voice, args.output, args.rate, args.pitch))


if __name__ == "__main__":
    main()
