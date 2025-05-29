package common

type Consumption interface {
	Add(Consumption)
	Summary() ConsumptionSummary
}

type ConsumptionSummary map[string]uint64

type UnknownConsumption struct{}

func (u *UnknownConsumption) Add(add Consumption) {
	// DO NOTHING
}

func (u *UnknownConsumption) Summary() ConsumptionSummary {
	return map[string]uint64{}
}

func NewSimpleConsumption(input, output uint64) ConsumptionSummary {
	return NewCachedConsumption(input, 0, output)
}

func NewCachedConsumption(input, cached, output uint64) ConsumptionSummary {
	return map[string]uint64{
		"input":  input,
		"cached": cached,
		"output": output,
	}
}
