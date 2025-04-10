package platform

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/redsift/go-cfg/dcfg"
	"github.com/redsift/go-siftjson"
)

const KEY_METADATA = "metadata"

var SiftMetadataV1Key = MapKey[siftjson.GUID, SiftMetadata](1, KEY_SIFTS, KEY_METADATA)

type SiftMetadataMap = dcfg.TypedMap[siftjson.GUID, SiftMetadata]

type SiftFlag uint

const (
	SIFT_FLAG_INTERNAL SiftFlag = 1 << iota
	SIFT_FLAG_VERBOSE_TAGS
)

type SiftMetadata struct {
	GUID  siftjson.GUID `json:"guid"`
	Name  string        `json:"name"`
	Flags SiftFlag      `json:"flags"`
}

func NewSiftMetadataMap(b dcfg.Backend) *SiftMetadataMap {
	res, _ := dcfg.NewTypedMap[siftjson.GUID, SiftMetadata](b, SiftMetadataV1Key)
	return (*SiftMetadataMap)(res)
}

func SiftMetadataServeHTTP(prefix string, m *SiftMetadataMap) http.Handler {
	prefix = "/" + strings.Trim(prefix, "/")
	mux := http.NewServeMux()
	mux.HandleFunc("GET "+prefix+"/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		data, err := m.Load(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
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
		data, err := m.GetKey(r.Context(), siftjson.GUID(guid))
		if err != nil {
			if errors.Is(err, dcfg.ErrNotFound) {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
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

		err := m.SetKey(r.Context(), siftjson.GUID(guid), content)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})

	mux.HandleFunc("POST "+prefix+"/{guid}/verbose", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		var verbose bool
		if _, err := fmt.Fscanf(r.Body, "%v", &verbose); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		guid := r.PathValue("guid")
		if len(guid) != 50 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		data, err := m.GetKey(r.Context(), siftjson.GUID(guid))
		if err != nil {
			if errors.Is(err, dcfg.ErrNotFound) {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		changed := false
		if verbose {
			changed = data.Flags&SIFT_FLAG_VERBOSE_TAGS == 0
			data.Flags |= SIFT_FLAG_VERBOSE_TAGS
		} else {
			changed = data.Flags&SIFT_FLAG_VERBOSE_TAGS != 0
			data.Flags &= ^SIFT_FLAG_VERBOSE_TAGS
		}

		if changed {
			if err := m.SetKey(r.Context(), siftjson.GUID(guid), data); err != nil {
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

		err := m.DelKey(r.Context(), siftjson.GUID(guid))
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
