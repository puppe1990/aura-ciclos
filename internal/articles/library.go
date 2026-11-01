package articles

type Article struct {
	Slug       string
	Title      string
	Summary    string
	Body       string
	ReviewedBy string
}

var library = []Article{
	{
		Slug:       "o-que-e-a-janela-fertil",
		Title:      "O que é a janela fértil",
		Summary:    "Cinco dias antes da ovulação e o dia seguinte concentram as chances de concepção.",
		ReviewedBy: "Dra. Marina Alves, ginecologia",
		Body: `A janela fértil não é um único dia. O óvulo vive cerca de 12 a 24 horas depois da ovulação, mas espermatozoides podem permanecer viáveis no trato reprodutivo por até cinco dias. Por isso o app marca do dia −5 ao dia +1 em torno da ovulação estimada.

Ciclos de 28 dias colocam a ovulação perto do dia 14. Ciclos mais longos empurram essa data; ciclos curtos a antecipam. Temperatura basal e testes de LH refinam o calendário, mas nenhum método caseiro substitui acompanhamento clínico quando há tentativa de gravidez com dificuldade.

A chance relativa de concepção é mais alta no dia da ovulação e no dia anterior. Fora da janela, a chance cai para perto de zero naquele ciclo.`,
	},
	{
		Slug:       "ciclo-irregular-nao-e-falha",
		Title:      "Ciclo irregular não é falha",
		Summary:    "Variação maior que quatro dias entre ciclos pede paciência e mais registros, não pânico.",
		ReviewedBy: "Dr. Paulo Nogueira, endocrinologia",
		Body: `Um ciclo é irregular quando o intervalo entre menstruações oscila bastante — no Aura Ciclos, desvio padrão de quatro dias ou mais entre os últimos intervalos. Estresse, treino pesado, viagem, COVID recente, amamentação e síndrome dos ovários policísticos são causas comuns.

O algoritmo usa a média dos últimos seis intervalos e baixa a confiança quando a variação é grande. Isso não é um diagnóstico. Se os ciclos passarem de 35 dias com frequência, se a menstruação sumir por três meses ou se o fluxo for tão intenso a ponto de limitar o dia, converse com uma ginecologista.

Registrar a data de início, mesmo sem o fim do fluxo, já melhora a previsão do mês seguinte.`,
	},
	{
		Slug:       "colica-quando-procurar-ajuda",
		Title:      "Cólica: o que é esperado e o que não é",
		Summary:    "Dor que cede com calor e analgésico comum é frequente; dor que paralisa merece consulta.",
		ReviewedBy: "Dra. Marina Alves, ginecologia",
		Body: `A cólica menstrual (dismenorreia) vem das prostaglandinas que contraem o útero para expulsar o endométrio. Calor local, movimento leve e anti-inflamatórios combinados com alimentação regular costumam reduzir a intensidade.

Procure atendimento se a dor for nova e crescente, se vier com desmaio, febre, corrimento fétido, ou se analgésicos habituais não fizerem efeito. Endometriose e adenomiose são hipóteses clínicas — o app não as diagnostica.

Anotar intensidade e dias no registro ajuda a mostrar o padrão na consulta, em vez de reconstruir de memória.`,
	},
	{
		Slug:       "tpm-humor-e-sono",
		Title:      "TPM, humor e sono",
		Summary:    "A fase lútea muda neurotransmissores; rastrear humor evita surpresa todo mês.",
		ReviewedBy: "Dra. Helena Costa, psiquiatria",
		Body: `Na fase lútea a progesterona sobe e, em algumas pessoas, o humor desce. Irritabilidade, choro fácil, ansiedade e sono fragmentado são relatos clássicos da tensão pré-menstrual. Isso dura os dias que antecedem o fluxo e tende a aliviar quando a menstruação começa.

O que o app faz: cruza humor e energia com o dia do ciclo para você ver se o padrão se repete. O que o app não faz: tratar TDPM, que é a forma intensa e incapacitante e pede avaliação psiquiátrica e ginecológica.

Cafeína tarde da noite, álcool e restrição extrema de carboidrato pioram o quadro em muita gente. Não é falha de caráter.`,
	},
	{
		Slug:       "temperatura-basal-sem-misterio",
		Title:      "Temperatura basal sem mistério",
		Summary:    "Três dias seguidos 0,2 °C acima da cobertura das seis noites anteriores sugerem ovulação.",
		ReviewedBy: "Dra. Marina Alves, ginecologia",
		Body: `A temperatura basal (BBT) sobe depois da ovulação por ação da progesterona. O Aura Ciclos usa uma regra simples: média das seis medições anteriores (coverline) e três manhãs consecutivas pelo menos 0,2 °C acima. A ovulação é marcada no dia anterior ao início dessa subida.

Meça sempre no mesmo horário, ainda na cama, com o mesmo termômetro. Gripe, álcool, noites curtas e acordar em horário diferente contaminam a curva. Teste de LH positivo no mesmo ciclo tem prioridade sobre o calendário e sobre a BBT.

A BBT confirma ovulação depois que ela aconteceu. Não serve para “pegar o dia” em tempo real — serve para calibrar o próximo ciclo.`,
	},
	{
		Slug:       "privacidade-do-ciclo",
		Title:      "Privacidade do seu ciclo",
		Summary:    "Modo anônimo não pede nome nem e-mail; os dados ficam no servidor da app, sem E2EE milagroso.",
		ReviewedBy: "Equipe Aura Ciclos",
		Body: `Ciclo, humor e libido são dados sensíveis. O Aura Ciclos oferece um modo anônimo: a conta nasce com um e-mail interno e um código de recuperação. Guarde o código — sem ele não há como reabrir a conta.

Os registros ficam em SQLite no servidor da aplicação, protegidos por sessão e senha (ou código). Não há criptografia ponta a ponta neste produto: quem opera o servidor consegue, em tese, ler o banco. Essa honestidade importa mais do que um selo falso.

Não compartilhe prints do ciclo em redes. No modo anônimo o app não mostra e-mail na interface. Em qualquer modo, o logout encerra a sessão neste navegador.`,
	},
	{
		Slug:       "gravidez-semana-a-semana",
		Title:      "Acompanhar gravidez semana a semana",
		Summary:    "A idade gestacional conta a partir da última menstruação; o parto estimado soma 280 dias.",
		ReviewedBy: "Dra. Marina Alves, ginecologia",
		Body: `A idade gestacional clínica começa no primeiro dia da última menstruação (DUM), mesmo que a concepção tenha ocorrido cerca de duas semanas depois. Por isso “semana 4” ainda é tão cedo: o embrião tem por volta de duas semanas.

O Aura Ciclos no modo gravidez mostra a semana a partir da DUM e uma data estimada somando 280 dias. É a regra de Naegele, não um compromisso com o calendário do parto. Ultrassom no primeiro trimestre recalibra melhor.

O app descreve marcos gerais (órgãos, movimentos, viabilidade) como educação, não como laudo. Sangramento, dor forte, febre ou diminuição de movimentos depois que eles já existiam pedem serviço de saúde, não o chat da Maya.`,
	},
	{
		Slug:       "pilula-e-lembretes",
		Title:      "Pílula, água e o lembrete que você controla",
		Summary:    "Esquecer a pílula muda a proteção; o app lembra no horário que você escolher.",
		ReviewedBy: "Dr. Paulo Nogueira, endocrinologia",
		Body: `Anticoncepcional oral combinado depende de tomada regular. Atrasos longos e comprimidos pulados reduzem a eficácia — a bula do seu método diz quantas horas de margem você tem e quando usar backup.

No Aura Ciclos os lembretes de pílula, água e início de menstruação são horários que você grava na conta. Eles aparecem na tela do ciclo. Notificação push do sistema operacional só funciona se o navegador permitir; o registro do horário existe mesmo sem push.

Se você esqueceu a pílula, registre “pílula esquecida” no diário. Isso não substitui a bula nem a orientação da prescritora.`,
	},
}

func All() []Article {
	out := make([]Article, len(library))
	copy(out, library)
	return out
}

func BySlug(slug string) (Article, bool) {
	for _, a := range library {
		if a.Slug == slug {
			return a, true
		}
	}
	return Article{}, false
}
