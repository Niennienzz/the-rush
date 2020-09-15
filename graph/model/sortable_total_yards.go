package model

type SortableTotalYardsByTeamResponse []*TotalYardsByTeamResponse

// Impl the sort.Interface interface.
func (x SortableTotalYardsByTeamResponse) Len() int {
	return len(x)
}

func (x SortableTotalYardsByTeamResponse) Less(i, j int) bool {
	return x[i].TotalYards < x[j].TotalYards
}

func (x SortableTotalYardsByTeamResponse) Swap(i, j int) {
	x[i], x[j] = x[j], x[i]
}
