package platform

import (
	"context"
	"encoding/json"
	"errors"
	"maps"
	"net/http"
	"slices"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/redsift/go-cfg/dcfg"
	"github.com/redsift/go-siftjson"
)

const KEY_METADATA = "metadata"

var SiftMetadataV1Key = MapKey[siftjson.GUID, SiftMetadata](1, KEY_SIFTS, KEY_METADATA)

type SiftMetadataMap = dcfg.TypedMap[siftjson.GUID, SiftMetadata]

type SiftFlag uint

const (
	SIFT_FLAG_INTERNAL SiftFlag = 1 << iota
	SIFT_FLAG_TEST
)

type SiftMetadata struct {
	GUID        siftjson.GUID     `json:"guid"`
	Name        string            `json:"name"`
	DisplayName string            `json:"display_name,omitempty"`
	Flags       SiftFlag          `json:"flags"`
	Tags        map[string]string `json:"tags,omitempty"`
}

func (m SiftMetadata) StatsdTags() []string {
	tags := make([]string, 0, len(m.Tags))
	for _, key := range slices.Sorted(maps.Keys(m.Tags)) {
		tags = append(tags, key+":"+m.Tags[key])
	}
	return tags
}

func NewSiftMetadataMap(b dcfg.Backend) *SiftMetadataMap {
	res, _ := dcfg.NewTypedMap[siftjson.GUID, SiftMetadata](b, SiftMetadataV1Key)
	return (*SiftMetadataMap)(res)
}

type SiftMetadataService struct {
	lock  sync.Mutex
	data  atomic.Pointer[map[siftjson.GUID]*atomic.Pointer[SiftMetadata]]
	store *SiftMetadataMap
}

func NewSiftMetadataService(ctx context.Context, store *SiftMetadataMap) (*SiftMetadataService, error) {
	s := &SiftMetadataService{
		store: store,
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	if err := store.Subscribe(ctx, s.subscriptionHandler); err != nil {
		return nil, err
	}

	data, err := s.store.Load(ctx)
	if err != nil {
		return nil, err
	}

	m := map[siftjson.GUID]*atomic.Pointer[SiftMetadata]{}
	for guid, meta := range data {
		m[guid] = new(atomic.Pointer[SiftMetadata])
		m[guid].Store(&meta)
	}

	s.data.Store(&m)

	return s, nil
}

func (s *SiftMetadataService) Get(guid siftjson.GUID) (p *atomic.Pointer[SiftMetadata], ok bool) {
	p, ok = (*s.data.Load())[guid]
	s.lock.Lock()
	return
}

func (s *SiftMetadataService) Set(ctx context.Context, guid siftjson.GUID, m SiftMetadata) error {
	data := *s.data.Load()
	if _, ok := data[guid]; !ok {
		s.lock.Lock()
		defer s.lock.Unlock()
		data = *s.data.Load()
		if _, ok := data[guid]; !ok {
			next := maps.Clone(data)
			p := new(atomic.Pointer[SiftMetadata])
			p.Store(&m)
			next[guid] = p
			s.data.Store(&next)
		}
	}

	return s.store.SetKey(ctx, guid, m)
}

func (s *SiftMetadataService) MakeHTTPHandler(prefix string) http.Handler {
	prefix = "/" + strings.Trim(prefix, "/")
	mux := http.NewServeMux()
	mux.HandleFunc("GET "+prefix+"/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")

		var (
			m    = *s.data.Load()
			data = make(map[siftjson.GUID]SiftMetadata, len(m))
		)

		for guid, p := range m {
			data[guid] = *p.Load()
		}

		json.NewEncoder(w).Encode(data)
	})

	mux.HandleFunc("GET "+prefix+"/{guid}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		guid := r.PathValue("guid")
		if len(guid) != 50 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		p, ok := (*s.data.Load())[siftjson.GUID(guid)]
		data := p.Load()

		if !ok || data == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(data)
	})

	mux.HandleFunc("POST "+prefix+"/{guid}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		var content SiftMetadata
		if err := json.NewDecoder(r.Body).Decode(&content); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		guid := r.PathValue("guid")
		if len(guid) != 50 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if content.GUID == "" {
			content.GUID = siftjson.GUID(guid)
		} else if content.GUID != siftjson.GUID(guid) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := s.Set(r.Context(), siftjson.GUID(guid), content); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})

	mux.HandleFunc("PUT "+prefix+"/{guid}/tags/{key}/{value}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")

		guid := r.PathValue("guid")
		if len(guid) != 50 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		key := r.PathValue("key")
		if len(key) < 2 || len(key) > 20 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		value := r.PathValue("value")
		if len(key) < 2 || len(key) > 100 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		p, ok := (*s.data.Load())[siftjson.GUID(guid)]
		data := p.Load()

		if !ok || data == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if old, ok := data.Tags[key]; !ok && value != old {
			if data.Tags == nil {
				data.Tags = make(map[string]string)
			}
			data.Tags[key] = value

			if err := s.Set(r.Context(), siftjson.GUID(guid), *data); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		json.NewEncoder(w).Encode(data)
	})

	mux.HandleFunc("DELETE "+prefix+"/{guid}/tags/{key}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")

		guid := r.PathValue("guid")
		if len(guid) != 50 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		key := r.PathValue("key")
		if len(key) < 2 || len(key) > 20 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		data, err := s.store.GetKey(r.Context(), siftjson.GUID(guid))
		if err != nil {
			if errors.Is(err, dcfg.ErrNotFound) {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if _, ok := data.Tags[key]; ok {
			delete(data.Tags, key)

			if err := s.store.SetKey(r.Context(), siftjson.GUID(guid), data); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		json.NewEncoder(w).Encode(data)
	})

	mux.HandleFunc("DELETE "+prefix+"/{guid}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")

		guid := r.PathValue("guid")
		if len(guid) != 50 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err := s.store.DelKey(r.Context(), siftjson.GUID(guid))
		if err != nil {
			if errors.Is(err, dcfg.ErrNotFound) {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})

	return mux
}

func (s *SiftMetadataService) subscriptionHandler(
	updated map[siftjson.GUID]SiftMetadata,
	removed []siftjson.GUID,
	err error,
) bool {
	s.lock.Lock()
	defer s.lock.Unlock()

	m := *s.data.Load()
	next := maps.Clone(m)

	for guid, meta := range updated {
		p, ok := m[guid]
		if !ok {
			p = new(atomic.Pointer[SiftMetadata])
			next[guid] = p
		}
		p.Store(&meta)
	}

	for _, guid := range removed {
		if p := next[guid]; p != nil {
			p.Store(nil)
		}
		delete(next, guid)
	}

	s.data.Store(&next)

	return true
}
