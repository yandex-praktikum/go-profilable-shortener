package auth

import (
	"testing"

	"github.com/gofrs/uuid"
)

func BenchmarkEncodeUID(b *testing.B) {
	uid := uuid.Must(uuid.NewV4())

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = EncodeUID(uid)
	}
}
