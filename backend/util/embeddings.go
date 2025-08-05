package util

import (
	"bytes"
	"database/sql/driver"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
)

type Embedding []float64

func (e *Embedding) Scan(value any) error {
	if value == nil {
		*e = nil
		return nil
	}

	switch v := value.(type) {
	case []byte:
		// Handle byte array directly
		if len(v)%8 != 0 {
			return errors.New("invalid byte length for float64 array")
		}

		embedding := make(Embedding, len(v)/8)
		reader := bytes.NewReader(v)
		for i := 0; i < len(embedding); i++ {
			var bits uint64
			if err := binary.Read(reader, binary.LittleEndian, &bits); err != nil {
				return errors.Join(err, errors.New("error reading embedding bytes"))
			}
			embedding[i] = math.Float64frombits(bits)
		}
		*e = embedding
		return nil

	case string:
		// Handle text representations if needed (not typical for binary storage)
		return fmt.Errorf("unsupported type: %T", value)

	default:
		return fmt.Errorf("incompatible type for Embedding: %T", value)
	}
}

func (e *Embedding) Value() (driver.Value, error) {
	if e == nil || len(*e) == 0 {
		return nil, nil
	}

	b := make([]byte, len(*e)*8)
	for i, v := range *e {
		bits := math.Float64bits(v)
		binary.LittleEndian.PutUint64(b[i*8:], bits)
	}

	return b, nil
}

func (e *Embedding) CosineSimilarity(other Embedding) (float64, error) {
	if e == nil {
		return 0, nil
	}

	this := *e
	if len(this) != len(other) {
		return 0.0, errors.New("embedding dimensions must match")
	}

	dotProduct := float64(0)
	normE := float64(0)
	normO := float64(0)

	for i := 0; i < len(this); i++ {
		dotProduct += this[i] * other[i]
		normE += math.Pow(this[i], 2)
		normO += math.Pow(other[i], 2)
	}

	denominator := math.Sqrt(normE * normO)
	if denominator == 0.0 {
		return 0.0, nil
	}

	return dotProduct / denominator, nil
}
