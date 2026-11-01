package maya

import (
	"strings"

	"github.com/puppe1990/aura-ciclos/internal/cycle"
)

const disclaimer = " Isso não substitui avaliação de uma profissional de saúde."

type Context struct {
	Phase     cycle.Phase
	CycleDay  int
	FertileOn bool
	ChancePct int
}

func Reply(question string, ctx Context) string {
	q := strings.TrimSpace(strings.ToLower(fold(question)))
	if q == "" {
		return "Pode mandar sua pergunta sobre ciclo, cólica, atraso, ovulação ou TPM. Eu respondo com o que o app já sabe do seu ritmo."
	}
	switch {
	case containsAny(q, "colica", "cólica", "dor na barriga", "dor menstrual"):
		return crampReply(ctx)
	case containsAny(q, "atras", "nao desceu", "não desceu", "menstruacao nao", "menstruação não"):
		return lateReply()
	case containsAny(q, "ovul", "fertil", "fértil", "conceb", "engrav"):
		return ovulationReply(ctx)
	case containsAny(q, "tpm", "irritad", "chor", "ansied"):
		return tpmReply(ctx)
	case containsAny(q, "cabeca", "cabeça", "enxaqueca"):
		return "Dor de cabeça no ciclo costuma acompanhar queda de estrogênio ou tensão. Hidratação, refeição regular e descanso em ambiente escuro ajudam; se a dor for súbita ou com aura, procure atendimento." + disclaimer
	default:
		return "Ainda não tenho uma resposta específica para isso, mas você pode registrar o sintoma no diário e ler os artigos revisados em Apoio. Se a dúvida for clínica, fale com uma profissional de saúde."
	}
}

func crampReply(ctx Context) string {
	base := "Cólica menstrual vem da contração do útero. Calor local, descanso e um anti-inflamatório que a sua médica já tenha indicado costumam aliviar."
	if ctx.Phase == cycle.PhaseMenstruation {
		return base + " Você está na menstruação agora — seja gentil com o ritmo de hoje." + disclaimer
	}
	return base + disclaimer
}

func lateReply() string {
	return "Atraso menstrual pode vir de estresse, viagem, mudança de peso, ciclo irregular ou gravidez. Se o atraso passar de alguns dias e houver relação sem proteção, um teste de gravidez tira a dúvida mais rápido do que esperar." + disclaimer
}

func ovulationReply(ctx Context) string {
	if ctx.FertileOn {
		return "Você está na janela fértil. A ovulação é o dia de maior chance de concepção; os cinco dias anteriores também importam por causa da vida do espermatozoide." + disclaimer
	}
	return "A ovulação costuma ocorrer cerca de 14 dias antes da próxima menstruação. A janela fértil cobre os cinco dias anteriores e o dia seguinte. O app recalcula isso quando você registra mais ciclos." + disclaimer
}

func tpmReply(ctx Context) string {
	base := "A TPM aparece na fase lútea, quando a progesterona sobe e o humor, o sono e o apetite mudam. Registrar irritabilidade, choro e energia ajuda a ver o padrão."
	if ctx.Phase == cycle.PhaseLuteal {
		return base + " Você está na fase lútea agora." + disclaimer
	}
	return base + disclaimer
}

func containsAny(q string, needles ...string) bool {
	for _, n := range needles {
		if strings.Contains(q, fold(n)) {
			return true
		}
	}
	return false
}

func fold(s string) string {
	repl := strings.NewReplacer(
		"á", "a", "à", "a", "ã", "a", "â", "a",
		"é", "e", "ê", "e",
		"í", "i",
		"ó", "o", "ô", "o", "õ", "o",
		"ú", "u", "ü", "u",
		"ç", "c",
	)
	return repl.Replace(strings.ToLower(s))
}
