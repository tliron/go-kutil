package problems

//
// Problematic
//

type Problematic interface {
	Problem() (string, string, string, int, int)
}
