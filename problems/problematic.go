package problems

//
// Problematic
//

type Problematic interface {
	Problem() (string, string, int, int)
}
