package wal

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"sync"

	"github.com/iotaledger/hive.go/logger"
	"github.com/iotaledger/wasp/packages/chain"
	"github.com/iotaledger/wasp/packages/iscp"
	"github.com/iotaledger/wasp/packages/state"
	"github.com/prometheus/client_golang/prometheus"
)

type WAL struct {
	dir      string
	log      *logger.Logger
	metrics  *walMetrics
	segments []*segment
	mu       sync.RWMutex //nolint
}

type chainWAL struct {
	*WAL
	chainID   *iscp.ChainID
	lastIndex uint32
}

func New(log *logger.Logger, dir string) *WAL {
	return &WAL{log: log, dir: dir, metrics: newWALMetrics()}
}

var _ chain.WAL = &chainWAL{}

type segmentFile interface {
	Stat() (os.FileInfo, error)
	io.Writer
	io.Closer
	io.Reader
}

type segment struct {
	segmentFile
	index   uint32
	dir     string
	corrupt bool
	name    string
}

func (w *WAL) NewChainWAL(chainID *iscp.ChainID) (chain.WAL, error) {
	if w == nil {
		return &defaultWAL{}, nil
	}
	w.dir = filepath.Join(w.dir, chainID.Base58())
	if err := os.MkdirAll(w.dir, 0o777); err != nil {
		return nil, fmt.Errorf("create dir: %w", err)
	}
	// read all segments in log
	f, err := os.Open(w.dir)
	if err != nil {
		return nil, fmt.Errorf("could not open wal: %w", err)
	}
	var segments []*segment
	files, _ := f.ReadDir(-1)
	for _, file := range files {
		w.metrics.segments.Inc()
		index, _ := strconv.ParseUint(file.Name(), 10, 32)
		segments = append(segments, &segment{index: uint32(index), dir: w.dir})
	}
	sort.SliceStable(segments, func(i, j int) bool {
		return segments[i].index < segments[j].index
	})
	var lastIndex uint32
	w.segments = segments
	if len(segments) > 0 {
		last := segments[len(segments)-1]
		w.metrics.latestSegment.Set(float64(last.index))
		lastIndex = last.index
	}
	return &chainWAL{w, chainID, lastIndex}, nil
}

func (w *chainWAL) Write(blocks ...state.Block) {
	w.mu.Lock()
	defer w.mu.Unlock()

	for _, block := range blocks {
		err := w.write(block)
		if err != nil {
			w.metrics.failedWrites.Inc()
			continue
		}
		w.metrics.latestSegment.Set(float64(block.BlockIndex()))
	}
}

func (w *chainWAL) write(block state.Block) error {
	var index uint32 = 1
	if len(w.segments) > 0 {
		index = w.segments[len(w.segments)-1].index + 1
	}
	segment, err := w.createSegment(index)
	if err != nil {
		return err
	}
	n, err := segment.Write(block.Bytes())
	if err != nil || len(block.Bytes()) != n {
		segment.corrupt = true
		return err
	}
	w.metrics.segments.Inc()
	return segment.Close()
}

func (w *chainWAL) createSegment(i uint32) (*segment, error) {
	segName := segmentName(w.dir, i)
	f, err := os.OpenFile(segName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o666)
	if err != nil {
		return nil, fmt.Errorf("could not create segment: %w", err)
	}
	s := &segment{index: i, segmentFile: f, dir: w.dir, name: segName}
	w.segments = append(w.segments, s)
	return s, nil
}

func segmentName(dir string, index uint32) string {
	return filepath.Join(dir, fmt.Sprintf("%010d", index))
}

func (w *chainWAL) Read() ([]state.Block, error) {
	blocks := make([]state.Block, len(w.segments))
	for _, segment := range w.segments {
		if err := segment.load(); err != nil {
			w.log.Debug(err)
			w.metrics.failedReads.Inc()
			continue
		}
		stat, err := segment.Stat()
		if err != nil {
			w.log.Debug(err)
			w.metrics.failedReads.Inc()
			continue
		}
		blockBytes := make([]byte, stat.Size())
		bufr := bufio.NewReader(segment)
		n, err := bufr.Read(blockBytes)
		if err != nil || int64(n) != stat.Size() {
			w.log.Debug("Error reading segment: %v", err)
			w.metrics.failedReads.Inc()
			continue
		}
		block, err := state.BlockFromBytes(blockBytes)
		if err != nil {
			w.log.Debug("Invalid block bytes")
			w.metrics.failedReads.Inc()
			continue
		}
		blocks = append(blocks, block)
	}
	return blocks, nil
}

func (s *segment) load() error {
	segName := segmentName(s.dir, s.index)
	f, err := os.OpenFile(segName, os.O_RDONLY, 0o666)
	if err != nil {
		return fmt.Errorf("error opening segment: %w", err)
	}
	s.segmentFile = f
	return nil
}

type defaultWAL struct{}

var _ chain.WAL = &defaultWAL{}

func (w *defaultWAL) Write(_ ...state.Block) {}

func (w *defaultWAL) Read() ([]state.Block, error) {
	return []state.Block{}, nil
}

func NewDefault() chain.WAL {
	return &defaultWAL{}
}

type walMetrics struct {
	segments      prometheus.Counter
	failedWrites  prometheus.Counter
	failedReads   prometheus.Counter
	latestSegment prometheus.Gauge
}

var once sync.Once

func newWALMetrics() *walMetrics {
	m := &walMetrics{}

	m.segments = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "wasp_wal_total_segments",
		Help: "Total number of segment files",
	})

	m.failedWrites = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "wasp_wal_failed_writes",
		Help: "Total number of writes to WAL that failed",
	})

	m.failedReads = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "wasp_wal_failed_reads",
		Help: "Total number of reads failed while replaying WAL",
	})

	m.latestSegment = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "wasp_wal_latest_segment",
		Help: "Last segment created",
	})

	registerMetrics := func() {
		prometheus.MustRegister(
			m.segments,
			m.failedWrites,
			m.failedReads,
			m.latestSegment,
		)
	}
	once.Do(registerMetrics)
	return m
}
