package db

import "log"

func GenerateFormulaViewModelPhases(fullFormulaIngredients []GetFullFormulaRow) []FullFormulaPhase {
	var formulaPhases = NewPhaseLinkedHashMap()
	var phase FullFormulaPhase

	for _, ingredient := range fullFormulaIngredients {
		phase, *formulaPhases = getOrCreatePhaseViewModel(*formulaPhases, ingredient)
		formulaIngredientModel := FullFormulaIngredient{
			Id:           ingredient.FormulaIngredientID,
			IngredientId: ingredient.IngredientID,
			Name:         ingredient.IngredientName,
			Inci:         ingredient.Inci,
			Percentage:   ingredient.Percentage,
			Cost:         float32((ingredient.Cost).Float64),
		}
		phase.FormulaIngredients = append(phase.FormulaIngredients, formulaIngredientModel)
		formulaPhases.Put(ingredient.PhaseID, phase)
	}

	formulaPhaseModels := new([]FullFormulaPhase)

	for _, key := range formulaPhases.Keys() {
		phase, _ := formulaPhases.Get(key)
		log.Println(phase.ID, phase.Name)
		*formulaPhaseModels = append(*formulaPhaseModels, phase)
	}

	return *formulaPhaseModels
}

func getOrCreatePhaseViewModel(phases PhaseLinkedHashMap, ingredient GetFullFormulaRow) (FullFormulaPhase, PhaseLinkedHashMap) {
	phase, exists := phases.Get(ingredient.PhaseID)
	if exists {
		return phase, phases
	}

	formulaIngredients := new([]FullFormulaIngredient)
	phase = FullFormulaPhase{
		ID:                 ingredient.PhaseID,
		Name:               ingredient.PhaseName,
		FormulaIngredients: *formulaIngredients,
		Description:        ingredient.PhaseDescription,
	}
	log.Println("created new Phase ViewModel: ", ingredient.PhaseName, ingredient.PhaseID)

	phases.Put(ingredient.PhaseID, phase)
	return phase, phases

}

type PhaseLinkedHashMap struct {
	keys   []int64
	values map[int64]FullFormulaPhase
}

func NewPhaseLinkedHashMap() *PhaseLinkedHashMap {
	return &PhaseLinkedHashMap{
		keys:   make([]int64, 0),
		values: make(map[int64]FullFormulaPhase),
	}
}

func (lhm *PhaseLinkedHashMap) Put(key int64, value FullFormulaPhase) {
	if _, ok := lhm.values[key]; !ok {
		lhm.keys = append(lhm.keys, key)
	}
	lhm.values[key] = value
}

func (lhm *PhaseLinkedHashMap) Get(key int64) (FullFormulaPhase, bool) {
	value, ok := lhm.values[key]
	return value, ok
}

func (lhm *PhaseLinkedHashMap) Keys() []int64 {
	return lhm.keys
}
