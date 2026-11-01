package symptoms

const (
	CategoryFluxo   = "fluxo"
	CategoryFisicos = "fisicos"
	CategoryHumor   = "humor"
	CategoryEstilo  = "estilo"
	CategorySexual  = "sexual"
)

type Symptom struct {
	ID       string
	Category string
	Label    string
}

type Group struct {
	Category string
	Label    string
	Items    []Symptom
}

var catalog = []Symptom{
	{ID: "fluxo_leve", Category: CategoryFluxo, Label: "Fluxo leve"},
	{ID: "fluxo_medio", Category: CategoryFluxo, Label: "Fluxo médio"},
	{ID: "fluxo_intenso", Category: CategoryFluxo, Label: "Fluxo intenso"},
	{ID: "coagulos", Category: CategoryFluxo, Label: "Coágulos"},
	{ID: "spotting", Category: CategoryFluxo, Label: "Spotting"},
	{ID: "sangramento_marrom", Category: CategoryFluxo, Label: "Sangramento marrom"},
	{ID: "ausencia_fluxo", Category: CategoryFluxo, Label: "Sem fluxo"},

	{ID: "colicas", Category: CategoryFisicos, Label: "Cólicas"},
	{ID: "dor_cabeca", Category: CategoryFisicos, Label: "Dor de cabeça"},
	{ID: "enxaqueca", Category: CategoryFisicos, Label: "Enxaqueca"},
	{ID: "acne", Category: CategoryFisicos, Label: "Acne"},
	{ID: "inchaco", Category: CategoryFisicos, Label: "Inchaço"},
	{ID: "seios_sensiveis", Category: CategoryFisicos, Label: "Seios sensíveis"},
	{ID: "dor_costas", Category: CategoryFisicos, Label: "Dor nas costas"},
	{ID: "dor_lombar", Category: CategoryFisicos, Label: "Dor lombar"},
	{ID: "dor_articulacoes", Category: CategoryFisicos, Label: "Dor nas articulações"},
	{ID: "nauseas", Category: CategoryFisicos, Label: "Náuseas"},
	{ID: "tontura", Category: CategoryFisicos, Label: "Tontura"},
	{ID: "ondas_calor", Category: CategoryFisicos, Label: "Ondas de calor"},
	{ID: "calafrios", Category: CategoryFisicos, Label: "Calafrios"},
	{ID: "suor", Category: CategoryFisicos, Label: "Suor"},
	{ID: "palpitacoes", Category: CategoryFisicos, Label: "Palpitações"},
	{ID: "prisao_ventre", Category: CategoryFisicos, Label: "Prisão de ventre"},
	{ID: "diarreia", Category: CategoryFisicos, Label: "Diarreia"},
	{ID: "gases", Category: CategoryFisicos, Label: "Gases"},
	{ID: "azia", Category: CategoryFisicos, Label: "Azia"},
	{ID: "apetite_alto", Category: CategoryFisicos, Label: "Apetite aumentado"},
	{ID: "apetite_baixo", Category: CategoryFisicos, Label: "Apetite diminuído"},
	{ID: "desejo_doce", Category: CategoryFisicos, Label: "Desejo por doce"},
	{ID: "desejo_salgado", Category: CategoryFisicos, Label: "Desejo por salgado"},
	{ID: "ovulacao_dor", Category: CategoryFisicos, Label: "Dor da ovulação"},
	{ID: "coceira", Category: CategoryFisicos, Label: "Coceira"},
	{ID: "corrimento", Category: CategoryFisicos, Label: "Corrimento"},
	{ID: "sensibilidade_pele", Category: CategoryFisicos, Label: "Pele sensível"},
	{ID: "cabelo_oleoso", Category: CategoryFisicos, Label: "Cabelo oleoso"},
	{ID: "cabelo_seco", Category: CategoryFisicos, Label: "Cabelo seco"},
	{ID: "olhos_inchados", Category: CategoryFisicos, Label: "Olhos inchados"},
	{ID: "nariz_entupido", Category: CategoryFisicos, Label: "Nariz entupido"},
	{ID: "dor_pelvica", Category: CategoryFisicos, Label: "Dor pélvica"},
	{ID: "caimbras", Category: CategoryFisicos, Label: "Cãibras"},
	{ID: "formigamento", Category: CategoryFisicos, Label: "Formigamento"},

	{ID: "ansiedade", Category: CategoryHumor, Label: "Ansiedade"},
	{ID: "felicidade", Category: CategoryHumor, Label: "Felicidade"},
	{ID: "irritabilidade", Category: CategoryHumor, Label: "Irritabilidade"},
	{ID: "tristeza", Category: CategoryHumor, Label: "Tristeza"},
	{ID: "choro", Category: CategoryHumor, Label: "Choro fácil"},
	{ID: "raiva", Category: CategoryHumor, Label: "Raiva"},
	{ID: "calma", Category: CategoryHumor, Label: "Calma"},
	{ID: "motivacao", Category: CategoryHumor, Label: "Motivação"},
	{ID: "desanimo", Category: CategoryHumor, Label: "Desânimo"},
	{ID: "estresse", Category: CategoryHumor, Label: "Estresse"},
	{ID: "concentracao_baixa", Category: CategoryHumor, Label: "Baixa concentração"},
	{ID: "confianca", Category: CategoryHumor, Label: "Confiança"},
	{ID: "inseguranca", Category: CategoryHumor, Label: "Insegurança"},
	{ID: "euforia", Category: CategoryHumor, Label: "Euforia"},
	{ID: "solidao", Category: CategoryHumor, Label: "Solidão"},
	{ID: "sensibilidade_emocional", Category: CategoryHumor, Label: "Sensibilidade emocional"},
	{ID: "fadiga", Category: CategoryHumor, Label: "Fadiga"},
	{ID: "energia_baixa", Category: CategoryHumor, Label: "Energia baixa"},
	{ID: "energia_alta", Category: CategoryHumor, Label: "Energia alta"},
	{ID: "insonia", Category: CategoryHumor, Label: "Insônia"},
	{ID: "sono_excessivo", Category: CategoryHumor, Label: "Sono excessivo"},

	{ID: "exercicio", Category: CategoryEstilo, Label: "Exercício"},
	{ID: "caminhada", Category: CategoryEstilo, Label: "Caminhada"},
	{ID: "yoga", Category: CategoryEstilo, Label: "Yoga"},
	{ID: "meditacao", Category: CategoryEstilo, Label: "Meditação"},
	{ID: "agua", Category: CategoryEstilo, Label: "Hidratação ok"},
	{ID: "pouca_agua", Category: CategoryEstilo, Label: "Pouca água"},
	{ID: "cafe", Category: CategoryEstilo, Label: "Café"},
	{ID: "alcool", Category: CategoryEstilo, Label: "Álcool"},
	{ID: "cigarro", Category: CategoryEstilo, Label: "Cigarro"},
	{ID: "medicamento", Category: CategoryEstilo, Label: "Medicamento"},
	{ID: "vitamina", Category: CategoryEstilo, Label: "Vitamina"},
	{ID: "anticoncepcional", Category: CategoryEstilo, Label: "Pílula tomada"},
	{ID: "pilula_esquecida", Category: CategoryEstilo, Label: "Pílula esquecida"},
	{ID: "viagem", Category: CategoryEstilo, Label: "Viagem"},
	{ID: "trabalho_pesado", Category: CategoryEstilo, Label: "Trabalho pesado"},
	{ID: "alimentacao_leve", Category: CategoryEstilo, Label: "Alimentação leve"},
	{ID: "alimentacao_pesada", Category: CategoryEstilo, Label: "Alimentação pesada"},

	{ID: "libido_alta", Category: CategorySexual, Label: "Libido alta"},
	{ID: "libido_baixa", Category: CategorySexual, Label: "Libido baixa"},
	{ID: "relacao_protegida", Category: CategorySexual, Label: "Relação com proteção"},
	{ID: "relacao_desprotegida", Category: CategorySexual, Label: "Relação sem proteção"},
	{ID: "muco_seco", Category: CategorySexual, Label: "Muco seco"},
	{ID: "muco_cremoso", Category: CategorySexual, Label: "Muco cremoso"},
	{ID: "muco_elastico", Category: CategorySexual, Label: "Muco elástico"},
	{ID: "teste_ovulacao_pos", Category: CategorySexual, Label: "Teste de ovulação +"},
	{ID: "teste_ovulacao_neg", Category: CategorySexual, Label: "Teste de ovulação −"},
	{ID: "teste_gravidez_pos", Category: CategorySexual, Label: "Teste de gravidez +"},
	{ID: "teste_gravidez_neg", Category: CategorySexual, Label: "Teste de gravidez −"},
}

var byID = func() map[string]Symptom {
	m := make(map[string]Symptom, len(catalog))
	for _, s := range catalog {
		m[s.ID] = s
	}
	return m
}()

func All() []Symptom {
	out := make([]Symptom, len(catalog))
	copy(out, catalog)
	return out
}

func ByID(id string) (Symptom, bool) {
	s, ok := byID[id]
	return s, ok
}

func Filter(ids []string) []Symptom {
	out := make([]Symptom, 0, len(ids))
	for _, id := range ids {
		if s, ok := ByID(id); ok {
			out = append(out, s)
		}
	}
	return out
}

func Grouped() []Group {
	order := []struct{ id, label string }{
		{CategoryFluxo, "Fluxo"},
		{CategoryFisicos, "Sintomas físicos"},
		{CategoryHumor, "Humor e energia"},
		{CategoryEstilo, "Estilo de vida"},
		{CategorySexual, "Fertilidade e sexual"},
	}
	buckets := map[string][]Symptom{}
	for _, s := range catalog {
		buckets[s.Category] = append(buckets[s.Category], s)
	}
	out := make([]Group, 0, len(order))
	for _, o := range order {
		out = append(out, Group{Category: o.id, Label: o.label, Items: buckets[o.id]})
	}
	return out
}

func CategoryLabel(id string) string {
	for _, g := range Grouped() {
		if g.Category == id {
			return g.Label
		}
	}
	return id
}
