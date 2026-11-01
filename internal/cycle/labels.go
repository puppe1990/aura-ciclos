package cycle

func PhaseLabel(p Phase) string {
	switch p {
	case PhaseMenstruation:
		return "Menstruação"
	case PhaseFollicular:
		return "Fase folicular"
	case PhaseFertile:
		return "Janela fértil"
	case PhaseOvulation:
		return "Ovulação"
	case PhaseLuteal:
		return "Fase lútea"
	case PhasePregnancy:
		return "Gravidez"
	default:
		return ""
	}
}

func ModeLabel(m Mode) string {
	switch m {
	case ModeTrying:
		return "Crescer a família"
	case ModePregnancy:
		return "Acompanhar gravidez"
	default:
		return "Acompanhar ciclo"
	}
}

func PregnancyNote(week int) string {
	switch {
	case week <= 0:
		return "Informe a data da última menstruação para acompanhar a semana."
	case week <= 4:
		return "O saco gestacional começa a se formar. Cansaço e seios sensíveis são comuns."
	case week <= 8:
		return "O coração do embrião já pulsa. Náusea e olfato aguçado podem aparecer."
	case week <= 12:
		return "Fim do primeiro trimestre se aproxima. Órgãos principais já estão esboçados."
	case week <= 16:
		return "O bebê cresce rápido. Algumas pessoas sentem os primeiros movimentos."
	case week <= 20:
		return "Meio da gestação. Ultrassom morfológico costuma acontecer nesta fase."
	case week <= 24:
		return "A pele do bebê ainda é fina. Viabilidade aumenta com cuidado neonatal."
	case week <= 28:
		return "Terceiro trimestre. Movimentos são mais claros; acompanhe o padrão."
	case week <= 32:
		return "O bebê ganha gordura e ouve sons. Cansaço e inchaço são frequentes."
	case week <= 36:
		return "O bebê desce na pelve em muitas gestações. Prepare a mala da maternidade."
	case week <= 40:
		return "Termo. O parto pode acontecer a qualquer momento; a data é uma estimativa."
	default:
		return "Passou da 40ª semana: a data estimada é só uma média. Fale com o pré-natal."
	}
}
