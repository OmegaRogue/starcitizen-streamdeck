//go:generate go-enum -f=$GOFILE --marshal --nocase -t ../../files/zerolog.gotmpl
package sc

// CompareOperation
/*
ENUM(
Empty=""
GreaterThan
NotEquals
)
*/
type CompareOperation string
