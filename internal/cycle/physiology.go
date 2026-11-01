package cycle

import (
	"fmt"
	"time"
)

type BodyNote struct {
	Hormones     string
	HormoneHint  string
	Vitality     string
	VitalityHint string
	Badge        string
}

func Physiology(p Phase) BodyNote {
	switch p {
	case PhaseMenstruation:
		return BodyNote{
			Hormones:     "Prostaglandinas",
			HormoneHint:  "Contrações do miométrio expulsam o endométrio e explicam a cólica.",
			Vitality:     "Ritmo baixo",
			VitalityHint: "Priorize calor, sono e refeições regulares.",
			Badge:        "Menstruação",
		}
	case PhaseFollicular:
		return BodyNote{
			Hormones:     "Estrogênio",
			HormoneHint:  "O folículo cresce e a energia costuma subir ao longo desta fase.",
			Vitality:     "Energia crescente",
			VitalityHint: "Bom momento para treinos e tarefas que pedem foco.",
			Badge:        "Fase folicular",
		}
	case PhaseFertile:
		return BodyNote{
			Hormones:     "LH & Estrogênio",
			HormoneHint:  "Pico hormonal antecipa a ovulação e eleva a energia basal.",
			Vitality:     "Energia alta",
			VitalityHint: "Ótimo momento para treinos de força, apresentações e conexões sociais.",
			Badge:        "Janela fértil",
		}
	case PhaseOvulation:
		return BodyNote{
			Hormones:     "LH & Estrogênio",
			HormoneHint:  "Pico hormonal desencadeia a liberação do óvulo e eleva a energia basal.",
			Vitality:     "Energia máxima",
			VitalityHint: "Ótimo momento para treinos de força, apresentações e conexões sociais.",
			Badge:        "Fase ovulatória • Janela fértil",
		}
	case PhaseLuteal:
		return BodyNote{
			Hormones:     "Progesterona",
			HormoneHint:  "O corpo lúteo prepara o endométrio; humor e sono podem oscilar.",
			Vitality:     "Ritmo de recuperação",
			VitalityHint: "Reduza o ritmo se a TPM apertar; calor e magnésio costumam ajudar.",
			Badge:        "Fase lútea",
		}
	case PhasePregnancy:
		return BodyNote{
			Hormones:     "hCG & Progesterona",
			HormoneHint:  "A gestação mantém o endométrio e muda o metabolismo basal.",
			Vitality:     "Gestação",
			VitalityHint: "Siga o pré-natal; o app só acompanha a semana a partir da última menstruação.",
			Badge:        "Gravidez",
		}
	default:
		return BodyNote{Badge: PhaseLabel(p)}
	}
}

func ConfidenceWords(level int) string {
	switch level {
	case 4:
		return "Excelente (Nível 4 de 4)"
	case 3:
		return "Boa (Nível 3 de 4)"
	case 2:
		return "Moderada (Nível 2 de 4)"
	default:
		return "Baixa (Nível 1 de 4)"
	}
}

func ChanceWords(pct int) string {
	switch {
	case pct >= 70:
		return "Chance Alta"
	case pct >= 40:
		return "Chance média"
	case pct > 0:
		return "Chance baixa"
	default:
		return ""
	}
}

func PregnancySize(week int) string {
	switch {
	case week <= 4:
		return "Semente de papoula"
	case week <= 8:
		return "Framboesa"
	case week <= 12:
		return "Lima"
	case week <= 16:
		return "Limão"
	case week <= 20:
		return "Banana"
	case week <= 24:
		return "Milho"
	case week <= 28:
		return "Berinjela"
	case week <= 32:
		return "Repolho"
	case week <= 36:
		return "Melão"
	default:
		return "Melancia"
	}
}

func FertileIndex(today, start, end time.Time) (n, total int) {
	if end.Before(start) {
		return 0, 0
	}
	total = daysBetween(start, end) + 1
	n = daysBetween(start, today) + 1
	if n < 1 {
		n = 1
	}
	if n > total {
		n = total
	}
	return n, total
}

func IntervalRange(lengths []int) string {
	if len(lengths) == 0 {
		return fmt.Sprintf("%d dias (padrão)", DefaultCycleDays)
	}
	min, max := lengths[0], lengths[0]
	for _, n := range lengths[1:] {
		if n < min {
			min = n
		}
		if n > max {
			max = n
		}
	}
	if min == max {
		return fmt.Sprintf("%d dias", min)
	}
	return fmt.Sprintf("%d–%d dias", min, max)
}
