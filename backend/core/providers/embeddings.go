package providers

import (
	"database/sql/driver"
	"encoding/binary"
	"fmt"
	"math"

	"github.com/pkg/errors"
)

type Embedding []float64

// Scan implements the sql.Scanner interface for Embedding type.
// It converts a database value to an Embedding object.
//
//goland:noinspection GoMixedReceiverTypes as needed by the sql.Scanner interface.
func (e *Embedding) Scan(value any) error {
	if value == nil {
		*e = nil
		return nil
	}

	b, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("unsupported type: %T", value)
	}
	if len(b)%8 != 0 {
		return errors.New("invalid byte length for float64 array")
	}

	n := len(b) / 8
	embedding := make(Embedding, n)
	for i := 0; i < n; i++ {
		start := i * 8
		bits := binary.LittleEndian.Uint64(b[start : start+8])
		embedding[i] = math.Float64frombits(bits)
	}
	*e = embedding
	return nil
}

// Value implements the driver.Valuer interface for Embedding type.
// It converts an Embedding object to a database value.
//
//goland:noinspection GoMixedReceiverTypes as needed by the driver.Valuer interface.
func (e Embedding) Value() (driver.Value, error) {
	if len(e) == 0 {
		return nil, nil
	}

	b := make([]byte, len(e)*8)
	for i, v := range e {
		bits := math.Float64bits(v)
		binary.LittleEndian.PutUint64(b[i*8:], bits)
	}

	return b, nil
}

// CosineSimilarity calculates the cosine similarity between two embeddings.
// It returns a float32 value representing the cosine of the angle between the vectors.
//
//goland:noinspection GoMixedReceiverTypes See Scan and Value methods
func (e *Embedding) CosineSimilarity(other Embedding) float64 {
	if e == nil {
		panic("nil Embedding")
	}
	if len(*e) != len(other) {
		panic(fmt.Sprintf("embedding dimensions must match (this %d, other %d)", len(*e), len(other)))
	}

	this := *e
	dotProduct := float64(0)

	for i := range this {
		v1, v2 := this[i], other[i]
		dotProduct += v1 * v2
	}

	return dotProduct
}

// Normalize scales the embedding vector to have unit length.
// It returns a new Embedding object with normalized values.
//
//goland:noinspection GoMixedReceiverTypes See Scan and Value methods
func (e *Embedding) Normalize() Embedding {
	if e == nil {
		panic("nil Embedding")
	}

	this := *e
	magnitude := float64(0)
	for _, v := range this {
		magnitude += v * v
	}
	magnitude = math.Sqrt(magnitude)

	if magnitude == 0 {
		// Handle zero vector case to avoid division by zero.
		return this
	}

	normalized := make(Embedding, len(this))
	for i, v := range this {
		normalized[i] = v / magnitude
	}
	return normalized
}
