package providers

import (
	"database/sql/driver"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
)

type Embeddings []float32

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

func (e *Embeddings) Value() (driver.Value, error) {
	if e == nil || len(*e) == 0 {
		return nil, nil
	}

	b := make([]byte, len(*e)*4)
	for i, v := range *e {
		bits := math.Float32bits(v)
		binary.LittleEndian.PutUint32(b[i*4:], bits)
	}

	return b, nil
}

func (e *Embeddings) CosineSimilarity(other Embeddings) (float32, error) {
	if e == nil || len(*e) != len(other) {
		return 0.0, errors.New("embedding dimensions must match")
	}

	dotProduct := float32(0)
	normE, normO := float32(0), float32(0)

	for i, v := range *e {
		dotProduct += v * other[i]
		normE += v * v
		normO += other[i] * other[i]
	}

	denominator := math.Sqrt(float64(normE * normO))
	if denominator == 0.0 {
		return 0.0, nil
	}
	return dotProduct / float32(denominator), nil
}
