#!/usr/bin/env python3
import argparse
import contextlib
import json
import math
import os
import shutil
import subprocess
import sys
import tempfile
import time
import wave
from pathlib import Path

import numpy as np
from openwakeword.model import Model

FRAME_SAMPLES = 1280
SAMPLE_RATE = 16000
CHANNELS = 1
SAMPLE_WIDTH = 2
WAKEWORD_MAP = {
    "jarvis": "hey_jarvis",
    "hey_jarvis": "hey_jarvis",
}


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="One-shot wake-word capture for Aurelia voice plane.")
    parser.add_argument("--device", default=os.environ.get("AURELIA_VOICE_DEVICE", "default"))
    parser.add_argument("--monitor-seconds", type=int, default=int(os.environ.get("AURELIA_VOICE_MONITOR_SECONDS", "3")))
    parser.add_argument("--followup-seconds", type=int, default=int(os.environ.get("AURELIA_VOICE_FOLLOWUP_SECONDS", "8")))
    parser.add_argument("--output-dir", default=os.environ.get("AURELIA_VOICE_DROP_PATH", ""))
    parser.add_argument("--model", default=os.environ.get("AURELIA_VOICE_WAKE_PHRASE", "jarvis"))
    parser.add_argument("--threshold", type=float, default=float(os.environ.get("AURELIA_VOICE_WAKE_THRESHOLD", "0.55")))
    parser.add_argument("--vad-threshold", type=float, default=float(os.environ.get("AURELIA_VOICE_VAD_THRESHOLD", "0.45")))
    parser.add_argument("--user-id", type=int, default=int(os.environ.get("AURELIA_VOICE_USER_ID", "0") or "0"))
    parser.add_argument("--chat-id", type=int, default=int(os.environ.get("AURELIA_VOICE_CHAT_ID", "0") or "0"))
    parser.add_argument("--source", default=os.environ.get("AURELIA_VOICE_SOURCE", "mic"))
    parser.add_argument("--requires-audio", action="store_true", default=os.environ.get("AURELIA_VOICE_REQUIRES_AUDIO", "").lower() == "true")
    parser.add_argument("--input-wav", default="")
    parser.add_argument("--debug", action="store_true")
    return parser.parse_args()


def resolve_model_name(name: str) -> str:
    normalized = name.strip().lower().replace(" ", "_")
    return WAKEWORD_MAP.get(normalized, normalized)


def arecord(device: str, seconds: int, dest: Path) -> None:
    cmd = [
        "arecord",
        "-q",
        "-D",
        device,
        "-f",
        "S16_LE",
        "-r",
        str(SAMPLE_RATE),
        "-c",
        str(CHANNELS),
        "-d",
        str(max(1, seconds)),
        "-t",
        "wav",
        str(dest),
    ]
    subprocess.run(cmd, check=True, stdout=subprocess.DEVNULL, stderr=subprocess.PIPE)


def read_wav(path: Path) -> np.ndarray:
    with contextlib.closing(wave.open(str(path), "rb")) as wav_file:
        sample_rate = wav_file.getframerate()
        channels = wav_file.getnchannels()
        sample_width = wav_file.getsampwidth()
        if sample_rate != SAMPLE_RATE:
            raise RuntimeError(f"unexpected sample rate {sample_rate}, expected {SAMPLE_RATE}")
        if channels != CHANNELS:
            raise RuntimeError(f"unexpected channel count {channels}, expected {CHANNELS}")
        if sample_width != SAMPLE_WIDTH:
            raise RuntimeError(f"unexpected sample width {sample_width}, expected {SAMPLE_WIDTH}")
        frames = wav_file.readframes(wav_file.getnframes())
    return np.frombuffer(frames, dtype=np.int16)


def write_wav(path: Path, pcm: np.ndarray) -> None:
    pcm = np.asarray(pcm, dtype=np.int16)
    with contextlib.closing(wave.open(str(path), "wb")) as wav_file:
        wav_file.setnchannels(CHANNELS)
        wav_file.setsampwidth(SAMPLE_WIDTH)
        wav_file.setframerate(SAMPLE_RATE)
        wav_file.writeframes(pcm.tobytes())


def detect_wakeword(samples: np.ndarray, model_name: str, threshold: float, vad_threshold: float) -> tuple[bool, float]:
    detector = Model(vad_threshold=vad_threshold)
    best_score = 0.0
    for start in range(0, len(samples), FRAME_SAMPLES):
        chunk = samples[start : start + FRAME_SAMPLES]
        if len(chunk) < FRAME_SAMPLES:
            chunk = np.pad(chunk, (0, FRAME_SAMPLES - len(chunk)))
        prediction = detector.predict(chunk)
        score = float(prediction.get(model_name, 0.0))
        best_score = max(best_score, score)
    return best_score >= threshold, best_score


def merge_audio(first: np.ndarray, second: np.ndarray) -> np.ndarray:
    if second.size == 0:
        return first
    return np.concatenate([first, second]).astype(np.int16)


def build_output_path(output_dir: Path) -> Path:
    output_dir.mkdir(parents=True, exist_ok=True)
    stamp = time.strftime("%Y%m%d-%H%M%S")
    return output_dir / f"voice-capture-{stamp}.wav"


def main() -> int:
    args = parse_args()
    model_name = resolve_model_name(args.model)
    output_dir = Path(args.output_dir).expanduser() if args.output_dir else None

    with tempfile.TemporaryDirectory(prefix="aurelia-voice-capture-") as temp_dir:
        temp_dir_path = Path(temp_dir)
        monitor_path = temp_dir_path / "monitor.wav"
        try:
            if args.input_wav:
                monitor_path = Path(args.input_wav).expanduser().resolve()
            else:
                arecord(args.device, args.monitor_seconds, monitor_path)

            monitor_pcm = read_wav(monitor_path)
            detected, score = detect_wakeword(monitor_pcm, model_name, args.threshold, args.vad_threshold)
            if not detected:
                if args.debug:
                    print(json.dumps({"detected": False, "score": round(score, 4), "model": model_name}), file=sys.stderr)
                return 0

            final_pcm = monitor_pcm
            if not args.input_wav and args.followup_seconds > 0:
                followup_path = temp_dir_path / "followup.wav"
                arecord(args.device, args.followup_seconds, followup_path)
                final_pcm = merge_audio(monitor_pcm, read_wav(followup_path))

            if output_dir is None:
                output_dir = temp_dir_path
            output_path = build_output_path(output_dir)
            write_wav(output_path, final_pcm)

            payload = {
                "detected": True,
                "audio_file": str(output_path),
                "user_id": args.user_id,
                "chat_id": args.chat_id,
                "requires_audio": args.requires_audio,
                "source": args.source,
                "delete_source_after": False,
                "score": round(score, 4),
                "model": model_name,
            }
            print(json.dumps(payload))
            return 0
        except subprocess.CalledProcessError as exc:
            message = exc.stderr.decode("utf-8", errors="ignore").strip() if exc.stderr else str(exc)
            print(f"voice capture failed: {message}", file=sys.stderr)
            return 1
        except Exception as exc:  # noqa: BLE001
            print(f"voice capture failed: {exc}", file=sys.stderr)
            return 1


if __name__ == "__main__":
    raise SystemExit(main())
