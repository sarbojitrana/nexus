"use client";

import { useEffect, useRef, useState } from "react";

function formatTime(seconds: number) {
  if (!Number.isFinite(seconds)) return "0:00";
  const m = Math.floor(seconds / 60);
  const s = Math.floor(seconds % 60);
  return `${m}:${String(s).padStart(2, "0")}`;
}

export function VideoPlayer({ src }: { src: string }) {
  const videoRef = useRef<HTMLVideoElement>(null);
  const [isPlaying, setIsPlaying] = useState(false);
  const [isMuted, setIsMuted] = useState(true);
  const [progress, setProgress] = useState(0);
  const [duration, setDuration] = useState(0);
  const [hasEnded, setHasEnded] = useState(false);

  useEffect(() => {
    const video = videoRef.current;
    if (!video) return;

    const onTime = () => setProgress(video.currentTime);
    const onMeta = () => setDuration(video.duration);
    const onPlay = () => {
      setIsPlaying(true);
      setHasEnded(false);
    };
    const onPause = () => setIsPlaying(false);
    const onEnded = () => {
      setIsPlaying(false);
      setHasEnded(true);
    };

    video.addEventListener("timeupdate", onTime);
    video.addEventListener("loadedmetadata", onMeta);
    video.addEventListener("play", onPlay);
    video.addEventListener("pause", onPause);
    video.addEventListener("ended", onEnded);
    return () => {
      video.removeEventListener("timeupdate", onTime);
      video.removeEventListener("loadedmetadata", onMeta);
      video.removeEventListener("play", onPlay);
      video.removeEventListener("pause", onPause);
      video.removeEventListener("ended", onEnded);
    };
  }, []);

  function togglePlay() {
    const video = videoRef.current;
    if (!video) return;
    if (video.paused) video.play();
    else video.pause();
  }

  function toggleMute() {
    const video = videoRef.current;
    if (!video) return;
    video.muted = !video.muted;
    setIsMuted(video.muted);
  }

  function replay() {
    const video = videoRef.current;
    if (!video) return;
    video.currentTime = 0;
    video.play();
  }

  function seek(e: React.ChangeEvent<HTMLInputElement>) {
    const video = videoRef.current;
    if (!video) return;
    const next = Number(e.target.value);
    video.currentTime = next;
    setProgress(next);
  }

  const pct = duration > 0 ? (progress / duration) * 100 : 0;

  return (
    <div className="group relative border border-border bg-black">
      <video
        ref={videoRef}
        src={src}
        muted
        playsInline
        preload="metadata"
        onClick={togglePlay}
        className="max-h-[520px] w-full cursor-pointer object-contain"
      />

      <div className="absolute inset-x-0 bottom-0 flex items-center gap-2.5 bg-gradient-to-t from-black/85 to-transparent px-3 pt-8 pb-2.5">
        <button
          onClick={hasEnded ? replay : togglePlay}
          aria-label={hasEnded ? "Replay" : isPlaying ? "Pause" : "Play"}
          className="shrink-0 font-mono text-[0.8rem] text-white hover:text-accent-strong"
        >
          {hasEnded ? "↻" : isPlaying ? "❚❚" : "▶"}
        </button>

        <span className="shrink-0 font-mono text-[0.66rem] text-white/80 tabular-nums">
          {formatTime(progress)} / {formatTime(duration)}
        </span>

        <div className="relative flex h-4 flex-1 items-center">
          <div className="absolute inset-x-0 h-[3px] bg-white/25" />
          <div className="absolute h-[3px] bg-accent" style={{ width: `${pct}%` }} />
          <input
            type="range"
            min={0}
            max={duration || 0}
            step={0.1}
            value={progress}
            onChange={seek}
            aria-label="Seek"
            className="relative h-4 w-full cursor-pointer appearance-none bg-transparent
              [&::-webkit-slider-thumb]:h-3 [&::-webkit-slider-thumb]:w-3
              [&::-webkit-slider-thumb]:appearance-none [&::-webkit-slider-thumb]:bg-accent
              [&::-moz-range-thumb]:h-3 [&::-moz-range-thumb]:w-3 [&::-moz-range-thumb]:border-0
              [&::-moz-range-thumb]:bg-accent"
          />
        </div>

        <button
          onClick={toggleMute}
          aria-label={isMuted ? "Unmute" : "Mute"}
          className="shrink-0 font-mono text-[0.78rem] text-white hover:text-accent-strong"
        >
          {isMuted ? "🔇" : "🔊"}
        </button>
      </div>
    </div>
  );
}
