"use client";

import { useEffect, useState } from "react";

type Weather = { temp: number; place: string };

const FALLBACK = { lat: 40.7128, lon: -74.006, place: "NEW YORK, NY" };

export function StatusReadout({ className = "" }: { className?: string }) {
  const [time, setTime] = useState<string | null>(null);
  const [weather, setWeather] = useState<Weather | null>(null);

  useEffect(() => {
    function tick() {
      const now = new Date();
      setTime(
        `${String(now.getHours()).padStart(2, "0")}:${String(now.getMinutes()).padStart(2, "0")}`
      );
    }
    tick();
    const timer = setInterval(tick, 20_000);
    return () => clearInterval(timer);
  }, []);

  useEffect(() => {
    let cancelled = false;

    async function load(lat: number, lon: number, place: string) {
      try {
        const res = await fetch(
          `https://api.open-meteo.com/v1/forecast?latitude=${lat}&longitude=${lon}&current=temperature_2m&temperature_unit=fahrenheit`
        );
        const json = await res.json();
        const temp = json?.current?.temperature_2m;
        if (!cancelled && typeof temp === "number") setWeather({ temp, place });
      } catch {
        // weather is decorative -- a failure just leaves the clock alone
      }
    }

    // Ask for a real position, but never block on it: the readout falls back
    // to a fixed city if the user denies or the browser has no geolocation.
    if (typeof navigator !== "undefined" && navigator.geolocation) {
      navigator.geolocation.getCurrentPosition(
        (pos) => load(pos.coords.latitude, pos.coords.longitude, "LOCAL"),
        () => load(FALLBACK.lat, FALLBACK.lon, FALLBACK.place),
        { timeout: 4000 }
      );
    } else {
      load(FALLBACK.lat, FALLBACK.lon, FALLBACK.place);
    }

    return () => {
      cancelled = true;
    };
  }, []);

  // Rendered only after mount so server and client markup can't disagree.
  if (!time) return <div className={className} aria-hidden />;

  return (
    <div className={`text-right font-mono text-[0.68rem] leading-[1.5] ${className}`}>
      <div className="text-text-muted">
        {time}
        {weather && <span className="text-text-faint"> — {Math.round(weather.temp)}°F</span>}
      </div>
      {weather && <div className="text-text-faint">{weather.place}</div>}
    </div>
  );
}
