package sequence

import (
	"time"

	"github.com/200sc/klangsynthese/audio"
	"github.com/200sc/klangsynthese/synth"
	"github.com/200sc/klangsynthese/wav"
)

type WaveGenerator struct {
	PitchPattern
	WavePattern
	VolumePattern
	HoldPattern
	Length
	Tick
	Loop
}

func NewWaveGenerator(opts ...Option) *WaveGenerator {
	wg := &WaveGenerator{}
	for _, opt := range opts {
		opt(wg)
	}
	return wg
}

func (wg *WaveGenerator) Generate() *Sequence {
	sq := &Sequence{}
	sq.Ticker = time.NewTicker(time.Duration(wg.Tick))
	sq.tickDuration = time.Duration(wg.Tick)
	sq.loop = bool(wg.Loop)
	sq.stopCh = make(chan error)
	if wg.Length == 0 {
		wg.Length = Length(len(wg.PitchPattern))
	}
	if len(wg.HoldPattern) == 0 {
		wg.HoldPattern = []time.Duration{sq.tickDuration}
	}
	sq.Pattern = make([]audio.Multi, wg.Length)

	controller := wav.NewController()

	pitchIndex := 0
	volumeIndex := 0
	waveIndex := 0
	holdIndex := 0
	for i := range sq.Pattern {
		p := wg.PitchPattern[pitchIndex]
		if p != synth.Rest {
			a, _ := controller.Wave(
				wg.WavePattern[waveIndex](
					p,
					wg.HoldPattern[holdIndex].Seconds(),
					wg.VolumePattern[volumeIndex],
				))
			sq.Pattern[i] = audio.Multi{[]audio.Audio{a}}
		} else {
			sq.Pattern[i] = audio.Multi{[]audio.Audio{}}
		}
		pitchIndex = (pitchIndex + 1) % len(wg.PitchPattern)
		volumeIndex = (volumeIndex + 1) % len(wg.VolumePattern)
		waveIndex = (waveIndex + 1) % len(wg.WavePattern)
		holdIndex = (holdIndex + 1) % len(wg.HoldPattern)
	}
	return sq
}
