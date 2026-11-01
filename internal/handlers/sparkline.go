package handlers

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/puppe1990/aura-ciclos/internal/models"
)

type sparklineView struct {
	Path   string
	Area   string
	LastX  float64
	LastY  float64
	Latest string
	Delta  string
	Empty  bool
}

func bbtSparkline(logs []models.DailyLog) sparklineView {
	type pt struct {
		day timeOrd
		v   float64
	}
	var pts []pt
	for _, l := range logs {
		if l.BBTCelsius == nil {
			continue
		}
		pts = append(pts, pt{day: timeOrd(l.LoggedOn.Unix()), v: *l.BBTCelsius})
	}
	sort.Slice(pts, func(i, j int) bool { return pts[i].day < pts[j].day })
	if len(pts) > 8 {
		pts = pts[len(pts)-8:]
	}
	if len(pts) == 0 {
		return sparklineView{Empty: true}
	}
	minV, maxV := pts[0].v, pts[0].v
	for _, p := range pts {
		if p.v < minV {
			minV = p.v
		}
		if p.v > maxV {
			maxV = p.v
		}
	}
	if maxV-minV < 0.15 {
		maxV = minV + 0.15
	}
	const w, h, pad = 300.0, 60.0, 6.0
	step := w
	if len(pts) > 1 {
		step = w / float64(len(pts)-1)
	}
	var b strings.Builder
	for i, p := range pts {
		x := step * float64(i)
		y := pad + (h-2*pad)*(1-(p.v-minV)/(maxV-minV))
		if i == 0 {
			fmt.Fprintf(&b, "M %.1f %.1f", x, y)
		} else {
			fmt.Fprintf(&b, " L %.1f %.1f", x, y)
		}
	}
	path := b.String()
	last := pts[len(pts)-1]
	lastX := step * float64(len(pts)-1)
	lastY := pad + (h-2*pad)*(1-(last.v-minV)/(maxV-minV))
	delta := ""
	if len(pts) > 1 {
		d := last.v - pts[len(pts)-2].v
		sign := "+"
		if d < 0 {
			sign = ""
		}
		delta = fmt.Sprintf("(%s%.1f°)", sign, d)
	}
	return sparklineView{
		Path:   path,
		Area:   path + fmt.Sprintf(" L %.1f %.1f L 0 %.1f Z", lastX, h, h),
		LastX:  lastX,
		LastY:  lastY,
		Latest: strings.ReplaceAll(fmt.Sprintf("%.2f °C", last.v), ".", ","),
		Delta:  delta,
	}
}

type timeOrd int64

func meanLengthLabel(lengths []int) string {
	if len(lengths) == 0 {
		return "—"
	}
	sum := 0
	for _, n := range lengths {
		sum += n
	}
	avg := float64(sum) / float64(len(lengths))
	return fmt.Sprintf("%.1f dias", math.Round(avg*10)/10)
}

type lengthBar struct {
	Days    int
	Label   string
	Height  int
	Current bool
	X       int
	LabelX  int
	ValueY  int
	BarY    int
	Fill    string
	Opacity string
}

func cycleLengthBars(starts []models.Period, lengths []int) []lengthBar {
	if len(lengths) == 0 {
		return nil
	}
	meses := []string{"Jan", "Fev", "Mar", "Abr", "Mai", "Jun", "Jul", "Ago", "Set", "Out", "Nov", "Dez"}
	n := len(lengths)
	out := make([]lengthBar, n)
	maxDays := 28
	for _, d := range lengths {
		if d > maxDays {
			maxDays = d
		}
	}
	gap := 48
	if n > 0 {
		gap = 240 / n
	}
	for i, days := range lengths {
		h := int(float64(days) / float64(maxDays) * 60)
		if h < 18 {
			h = 18
		}
		x := 30 + i*gap
		label := fmt.Sprintf("C%d", i+1)
		if i+1 < len(starts) {
			m := starts[i].StartedOn.UTC().Month()
			label = meses[m-1]
		} else if i < len(starts) {
			m := starts[i].StartedOn.UTC().Month()
			label = meses[m-1]
		}
		current := i == n-1
		fill := "#c94a38"
		op := "0.8"
		if current {
			fill = "#a83223"
			op = "1"
		}
		out[i] = lengthBar{
			Days:    days,
			Label:   label,
			Height:  h,
			Current: current,
			X:       x,
			LabelX:  x + 11,
			BarY:    105 - h,
			ValueY:  105 - h - 7,
			Fill:    fill,
			Opacity: op,
		}
	}
	return out
}

func avatarLetter(name string) string {
	name = strings.TrimSpace(name)
	if name == "" || name == "você" {
		return "A"
	}
	r := []rune(name)
	return strings.ToUpper(string(r[0]))
}

func waterLiters(glasses int) string {
	if glasses <= 0 {
		return "0,0"
	}
	return strings.ReplaceAll(fmt.Sprintf("%.1f", float64(glasses)*0.25), ".", ",")
}

func sleepQuality(hours float64) string {
	switch {
	case hours <= 0:
		return ""
	case hours >= 7:
		return "Boa"
	case hours >= 6:
		return "Ok"
	default:
		return "Curta"
	}
}

func sleepLabel(hours float64) string {
	if hours <= 0 {
		return "—"
	}
	h := int(hours)
	m := int(math.Round((hours - float64(h)) * 60))
	if m == 0 {
		return fmt.Sprintf("%dh", h)
	}
	return fmt.Sprintf("%dh %02dm", h, m)
}
