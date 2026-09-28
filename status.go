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
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// handleStatus proxies GET $AIP_BASE_URL/status?version=<appVersion> and
// returns the response body, adding `trustant` from the legacy `trustable`
// block when needed (see withLegacyTrustantCatalog). Used by the splash and applist pages
// to receive per-provider model catalogs and operator banner messages.
func handleStatus(w http.ResponseWriter, r *http.Request) {
	if expiredGuard(w) {
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if AIPBaseURL == "" {
		writeStatusError(w, "AIP_BASE_URL is not set")
		return
	}

	endpoint := AIPBaseURL + "/status?version=" + url.QueryEscape(appVersion)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(endpoint)
	if err != nil {
		writeStatusError(w, fmt.Sprintf("status fetch failed: %s", err))
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		writeStatusError(w, fmt.Sprintf("status read failed: %s", err))
		return
	}
	if resp.StatusCode != http.StatusOK {
		writeStatusError(w, fmt.Sprintf("status returned HTTP %d", resp.StatusCode))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	_, _ = w.Write(withLegacyTrustantCatalog(body))
}

// withLegacyTrustantCatalog copies the pre-rebrand `trustable` provider block
// to `trustant` when the ai-proxy has not been renamed yet, so every page can
// read `status.trustant`. `trustant` wins when both exist; the legacy key is
// kept. A body that is not a JSON object is returned unchanged.
func withLegacyTrustantCatalog(body []byte) []byte {
	var status map[string]json.RawMessage
	if err := json.Unmarshal(body, &status); err != nil || status == nil {
		return body
	}
	if _, ok := status["trustant"]; ok {
		return body
	}
	legacy, ok := status["trustable"]
	if !ok {
		return body
	}
	status["trustant"] = legacy
	out, err := json.Marshal(status)
	if err != nil {
		return body
	}
	return out
}

func writeStatusError(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(http.StatusBadGateway)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
