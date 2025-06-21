package snowflake

import (
	"fmt"
	"sync"
	"time"
)

const (
	// Epoch adalah waktu awal kustom (dalam milidetik).
	// Anda bisa mengatur ini ke tanggal tertentu, misalnya tanggal proyek dimulai.
	// Contoh: 1 Januari 2024 UTC
	epoch int64 = 1704067200000 // time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).UnixMilli()

	workerIDBits uint8 = 10 // Jumlah bit untuk worker ID (0-1023)
	sequenceBits uint8 = 12 // Jumlah bit untuk sequence number (0-4095)

	maxWorkerID int64 = -1 ^ (-1 << workerIDBits) // 1023
	maxSequence int64 = -1 ^ (-1 << sequenceBits) // 4095

	timestampShift = workerIDBits + sequenceBits // 22
	workerIDShift  = sequenceBits                // 12
)

type Generator interface {
	// NextID menghasilkan snowflake ID unik berikutnya.
	NextID() (int64, error)
}

// Generator adalah struct untuk menghasilkan snowflake ID.
type generator struct {
	mu            sync.Mutex
	lastTimestamp int64
	workerID      int64
	sequence      int64
}

func NewGenerator(workerID int64) (Generator, error) {
	if workerID < 0 || workerID > maxWorkerID {
		return nil, fmt.Errorf("worker ID must be between 0 and %d", maxWorkerID)
	}
	return &generator{
		workerID: workerID,
	}, nil
}

// NextID menghasilkan snowflake ID unik berikutnya.
func (sg *generator) NextID() (int64, error) {
	sg.mu.Lock()
	defer sg.mu.Unlock()

	currentTimestamp := time.Now().UnixMilli()

	// Jika jam mundur, kita perlu menunggu atau mengembalikan error.
	// Di sini kita akan mengembalikan error untuk kesederhanaan.
	if currentTimestamp < sg.lastTimestamp {
		// Anda bisa menambahkan logika untuk menunggu jika perbedaan waktunya kecil,
		// atau mengembalikan error jika terlalu besar.
		return 0, fmt.Errorf("clock moved backwards. refusing to generate id for %d milliseconds", sg.lastTimestamp-currentTimestamp)
	}

	if currentTimestamp == sg.lastTimestamp {
		sg.sequence = (sg.sequence + 1) & maxSequence
		// Jika sequence meluap (kembali ke 0), kita perlu menunggu milidetik berikutnya.
		if sg.sequence == 0 {
			currentTimestamp = sg.tilNextMillis(sg.lastTimestamp)
		}
	} else {
		// Waktu berbeda, reset sequence.
		sg.sequence = 0
	}

	sg.lastTimestamp = currentTimestamp

	// Gabungkan semua bagian untuk membuat ID.
	id := ((currentTimestamp - epoch) << timestampShift) |
		(sg.workerID << workerIDShift) |
		sg.sequence

	return id, nil
}

// tilNextMillis akan memblokir hingga milidetik berikutnya.
func (sg *generator) tilNextMillis(lastTimestamp int64) int64 {
	timestamp := time.Now().UnixMilli()
	for timestamp <= lastTimestamp {
		timestamp = time.Now().UnixMilli()
	}
	return timestamp
}

// ParseID mengurai snowflake ID dan mengembalikan komponen-komponennya.
// Ini berguna untuk debugging atau analisis.
func ParseID(id int64) (timestampMs int64, workerID int64, sequence int64, t time.Time) {
	timestampMs = (id >> timestampShift) + epoch
	workerID = (id >> workerIDShift) & maxWorkerID
	sequence = id & maxSequence
	t = time.UnixMilli(timestampMs)
	return
}
