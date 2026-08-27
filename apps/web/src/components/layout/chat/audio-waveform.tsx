"use client"

import { useEffect, useRef, useState } from "react"

interface AudioWaveformProps {
  isRecording: boolean
  stream?: MediaStream | null
}

export function AudioWaveform({ isRecording, stream }: AudioWaveformProps) {
  const [levels, setLevels] = useState<number[]>(Array(24).fill(0.1))
  const analyserRef = useRef<AnalyserNode | null>(null)
  const animationRef = useRef<number | null>(null)
  const audioContextRef = useRef<AudioContext | null>(null)
  const historyRef = useRef<number[]>([])

  useEffect(() => {
    if (!isRecording || !stream) {
      setLevels(
        Array(24)
          .fill(0)
          .map(() => 0.05 + Math.random() * 0.1),
      )
      historyRef.current = []
      return
    }

    const audioContext = new (window.AudioContext || (window as any).webkitAudioContext)()
    audioContextRef.current = audioContext

    const analyser = audioContext.createAnalyser()
    analyser.fftSize = 256
    analyser.smoothingTimeConstant = 0.3
    analyserRef.current = analyser

    const source = audioContext.createMediaStreamSource(stream)
    source.connect(analyser)

    const dataArray = new Uint8Array(analyser.fftSize)

    const updateLevels = () => {
      if (!analyserRef.current) return

      analyserRef.current.getByteTimeDomainData(dataArray)

      // Calculate RMS (Root Mean Square) for better volume detection
      let sumSquares = 0
      for (let i = 0; i < dataArray.length; i++) {
        const normalized = (dataArray[i] - 128) / 128
        sumSquares += normalized * normalized
      }
      const rms = Math.sqrt(sumSquares / dataArray.length)

      const amplifiedLevel = Math.min(1, rms * 4)
      const avgLevel = Math.max(0.15, amplifiedLevel)

      historyRef.current.unshift(avgLevel)

      if (historyRef.current.length > 12) {
        historyRef.current.pop()
      }

      const newLevels = Array(24).fill(0.1)
      const center = 12

      for (let i = 0; i < historyRef.current.length; i++) {
        const level = historyRef.current[i]
        const leftIndex = center - 1 - i
        const rightIndex = center + i

        if (leftIndex >= 0) newLevels[leftIndex] = level
        if (rightIndex < 24) newLevels[rightIndex] = level
      }

      setLevels(newLevels)
      animationRef.current = requestAnimationFrame(updateLevels)
    }

    updateLevels()

    return () => {
      if (animationRef.current) {
        cancelAnimationFrame(animationRef.current)
      }
      if (audioContextRef.current) {
        audioContextRef.current.close()
      }
    }
  }, [isRecording, stream])

  return (
    <div className="flex items-center justify-center gap-0.75 h-8 flex-1 px-2">
      {levels.map((level, index) => (
        <div
          key={index}
          className="bg-zinc-400 rounded-full transition-all duration-75"
          style={{
            width: "3px",
            height: `${Math.max(4, level * 28)}px`,
          }}
        />
      ))}
    </div>
  )
}
