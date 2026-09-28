// Copyright 2025-2026 Nuvolaris Inc
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package main

import (
	"encoding/json"
	"testing"
)

func TestWithLegacyTrustantCatalog(t *testing.T) {
	keys := func(t *testing.T, body []byte) map[string]json.RawMessage {
		t.Helper()
		var m map[string]json.RawMessage
		if err := json.Unmarshal(body, &m); err != nil {
			t.Fatalf("invalid JSON: %v", err)
		}
		return m
	}

	t.Run("copies trustable when trustant is missing", func(t *testing.T) {
		got := keys(t, withLegacyTrustantCatalog([]byte(`{"message":"","trustable":{"modelsVersion":3}}`)))
		if string(got["trustant"]) != `{"modelsVersion":3}` {
			t.Fatalf("trustant = %s", got["trustant"])
		}
		if _, ok := got["trustable"]; !ok {
			t.Fatal("legacy trustable key was removed")
		}
	})

	t.Run("prefers trustant when both exist", func(t *testing.T) {
		in := `{"trustant":{"modelsVersion":5},"trustable":{"modelsVersion":3}}`
		if got := string(withLegacyTrustantCatalog([]byte(in))); got != in {
			t.Fatalf("body changed: %s", got)
		}
	})

	t.Run("leaves a trustant-only body unchanged", func(t *testing.T) {
		in := `{"trustant":{"modelsVersion":5},"ollama":{}}`
		if got := string(withLegacyTrustantCatalog([]byte(in))); got != in {
			t.Fatalf("body changed: %s", got)
		}
	})

	t.Run("passes a non-object body through", func(t *testing.T) {
		for _, in := range []string{`[1,2]`, `not json`, `null`} {
			if got := string(withLegacyTrustantCatalog([]byte(in))); got != in {
				t.Fatalf("body %q changed to %q", in, got)
			}
		}
	})
}
