//go:generate go-enum -f=$GOFILE --marshal --nocase -t ../../files/zerolog.gotmpl
package sc

// Rank is an enumeration of Elite Dangerous Ranks
/*
ENUM(
Empty=""
GreaterThan
NotEquals
)
*/
type CompareOperation string
