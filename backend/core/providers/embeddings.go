package providers

import (
	"database/sql/driver"
	"encoding/binary"
	"fmt"
	"math"

	"github.com/pkg/errors"
)

type Embeddings []float32

//goland:noinspection GoMixedReceiverTypes as needed by the sql.Scanner interface.
func (e *Embeddings) Scan(value any) error {
	if value == nil {
		*e = nil
		return nil
	}

	b, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("unsupported type: %T", value)
	}
	if len(b)%4 != 0 {
		return errors.New("invalid byte length for float32 array")
	}

	n := len(b) / 4
	embedding := make(Embeddings, n)
	for i := 0; i < n; i++ {
		start := i * 4
		bits := binary.LittleEndian.Uint32(b[start : start+4])
		embedding[i] = math.Float32frombits(bits)
	}
	*e = embedding
	return nil
}

//goland:noinspection GoMixedReceiverTypes as needed by the driver.Valuer interface.
func (e Embeddings) Value() (driver.Value, error) {
	if len(e) == 0 {
		return nil, nil
	}

	b := make([]byte, len(e)*4)
	for i, v := range e {
		bits := math.Float32bits(v)
		binary.LittleEndian.PutUint32(b[i*4:], bits)
	}

	return b, nil
}

//goland:noinspection GoMixedReceiverTypes See Scan and Value methods
func (e *Embeddings) CosineSimilarity(other Embeddings) (float32, error) {
	if e == nil {
		return 0.0, errors.New("nil Embeddings")
	}
	if len(*e) != len(other) {
		return 0.0, errors.Errorf("embedding dimensions must match (this %d, other %d)", len(*e), len(other))
	}

	dotProduct := float32(0)
	normE, normO := float32(0), float32(0)

	for i := range *e {
		v1, v2 := (*e)[i], other[i]
		dotProduct += v1 * v2
		normE += v1 * v1
		normO += v2 * v2
	}

	productOfNorms := normE * normO
	if productOfNorms == 0 {
		return 0.0, nil
	}

	denominator := float32(math.Sqrt(float64(productOfNorms)))
	return dotProduct / denominator, nil
}
