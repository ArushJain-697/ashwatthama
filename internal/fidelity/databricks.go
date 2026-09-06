package fidelity

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// databricks.go is the seam §12's Databricks module calls through, mirroring
// adapter.go's own pattern: every real network/workspace call goes through
// one small interface, so wiring real credentials in later touches this file
// only — CalibrateAdvisoryEdges and the pure math in calibration.go never
// change.
type DatabricksClient interface {
	// ScoreEdge is §12.2's I-CALM extraction call: given an edge's hypothesis
	// in natural language and the constrained vocabulary it must ground to,
	// return a verbal confidence in [0,1]. Never invents a score outside that
	// range — same fail-closed posture as ResolveFallback (intent.go).
	ScoreEdge(ctx context.Context, hypothesis string, vocabulary []string) (float64, error)
	// Corroborate is §12.4's AI Search hybrid-retrieval call for a
	// pending_verification edge or an unverifiable_coverage entity (§12.15):
	// ask whether the repo's own docs/comments/commits confirm the
	// hypothesis. verified=false with no error means "no corroborating
	// evidence found," not a failure.
	Corroborate(ctx context.Context, hypothesis string) (verified bool, evidenceSnippet string, err error)
}

// UnavailableDatabricksClient is the safe default when calibration_enabled
// is false or no workspace is configured — every advisory/unverifiable
// entry's Databricks-reserved fields stay null, exactly as the schema
// promises (#26's "New fields ... null when it's off").
type UnavailableDatabricksClient struct{}

func (UnavailableDatabricksClient) ScoreEdge(context.Context, string, []string) (float64, error) {
	return 0, fmt.Errorf("Databricks client is not configured")
}

func (UnavailableDatabricksClient) Corroborate(context.Context, string) (bool, string, error) {
	return false, "", fmt.Errorf("Databricks client is not configured")
}

// DatabricksHTTPClient is the real, network-calling DatabricksClient
// implementation, verified against a live workspace during this session
// (see BUILDATHON.md). It reads credentials only from the environment
// (DATABRICKS_HOST / DATABRICKS_TOKEN, the same names the official SDK and
// CLI use) — never from a config file this repository would commit.
type DatabricksHTTPClient struct {
	Host              string
	Token             string
	ServingEndpoint   string // I-CALM scoring model, e.g. "databricks-meta-llama-3-1-8b-instruct"
	VectorSearchIndex string // "<catalog>.<schema>.<index>" — empty means Corroborate always fails closed
	HTTP              *http.Client
}

// NewDatabricksHTTPClientFromEnv builds a client from DATABRICKS_HOST and
// DATABRICKS_TOKEN. Returns an error (never a client with an empty
// credential) if either is unset, so a misconfigured environment fails
// closed at startup rather than on the first real call.
func NewDatabricksHTTPClientFromEnv(servingEndpoint, vectorSearchIndex string) (DatabricksHTTPClient, error) {
	host := strings.TrimRight(os.Getenv("DATABRICKS_HOST"), "/")
	token := os.Getenv("DATABRICKS_TOKEN")
	if host == "" || token == "" {
		return DatabricksHTTPClient{}, fmt.Errorf("DATABRICKS_HOST and DATABRICKS_TOKEN must both be set")
	}
	return DatabricksHTTPClient{
		Host: host, Token: token, ServingEndpoint: servingEndpoint, VectorSearchIndex: vectorSearchIndex,
		HTTP: &http.Client{Timeout: 30 * time.Second},
	}, nil
}

// icalmScorePrompt is §12.2's I-CALM structure: full reward for a correct
// edge, heavy penalty for an incorrect one, partial credit for abstaining,
// and an explicit instruction not to assume unverified relationships —
// verified against a live databricks-meta-llama-3-1-8b-instruct call to
// return a clean, parseable JSON confidence during this session.
const icalmScorePrompt = `You are scoring whether a code relationship hypothesis is correct.
Hypothesis: %s
Respond with ONLY a JSON object of the exact shape {"confidence": <float 0-1>}. No other text.
A correct high-confidence answer is rewarded. An incorrect high-confidence answer is heavily penalized. Abstaining (confidence near 0.5) when unsure earns partial credit. Do not assume unverified relationships.`

func (client DatabricksHTTPClient) ScoreEdge(ctx context.Context, hypothesis string, vocabulary []string) (float64, error) {
	prompt := fmt.Sprintf(icalmScorePrompt, hypothesis)
	if len(vocabulary) > 0 {
		prompt += "\nGround your answer only in this vocabulary if the hypothesis names an entity: " + strings.Join(vocabulary, ", ")
	}
	body, err := json.Marshal(map[string]any{
		"messages":    []map[string]string{{"role": "user", "content": prompt}},
		"max_tokens":  50,
		"temperature": 0,
	})
	if err != nil {
		return 0, fmt.Errorf("encode I-CALM request: %w", err)
	}
	url := client.Host + "/serving-endpoints/" + client.ServingEndpoint + "/invocations"
	response, err := client.post(ctx, url, body)
	if err != nil {
		return 0, err
	}
	var parsed struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(response, &parsed); err != nil || len(parsed.Choices) == 0 {
		return 0, fmt.Errorf("I-CALM response did not parse: %w", err)
	}
	var score struct {
		Confidence float64 `json:"confidence"`
	}
	if err := json.Unmarshal([]byte(strings.TrimSpace(parsed.Choices[0].Message.Content)), &score); err != nil {
		return 0, fmt.Errorf("I-CALM did not return the required {\"confidence\": N} shape: %w", err)
	}
	if score.Confidence < 0 || score.Confidence > 1 {
		return 0, fmt.Errorf("I-CALM returned an out-of-range confidence %v — fail closed rather than trust it", score.Confidence)
	}
	return score.Confidence, nil
}

// Corroborate queries the AI Search hybrid-retrieval index (§12.4) built
// over the repo's own docs/comments/commit messages. Populating that index
// (the Delta table + embedding + sync pipeline) is separate work this
// session did not complete — see BUILDATHON.md; this is the real query-side
// call, verified to reach a live, ONLINE Vector Search endpoint, and it
// fails closed (verified=false, no panic) exactly like an unpopulated index
// would, which CalibrateAdvisoryEdges already treats as "attempted, no
// evidence" rather than a crash.
func (client DatabricksHTTPClient) Corroborate(ctx context.Context, hypothesis string) (bool, string, error) {
	if client.VectorSearchIndex == "" {
		return false, "", fmt.Errorf("no vector search index configured")
	}
	body, err := json.Marshal(map[string]any{
		"query_text":  hypothesis,
		"columns":     []string{"content"},
		"num_results": 1,
		"query_type":  "hybrid",
	})
	if err != nil {
		return false, "", fmt.Errorf("encode corroboration query: %w", err)
	}
	url := client.Host + "/api/2.0/vector-search/indexes/" + client.VectorSearchIndex + "/query"
	response, err := client.post(ctx, url, body)
	if err != nil {
		return false, "", err
	}
	var parsed struct {
		Result struct {
			DataArray [][]any `json:"data_array"`
		} `json:"result"`
	}
	if err := json.Unmarshal(response, &parsed); err != nil {
		return false, "", fmt.Errorf("corroboration response did not parse: %w", err)
	}
	if len(parsed.Result.DataArray) == 0 {
		return false, "", nil
	}
	snippet := fmt.Sprintf("%v", parsed.Result.DataArray[0][0])
	return true, snippet, nil
}

func (client DatabricksHTTPClient) post(ctx context.Context, url string, body []byte) ([]byte, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	request.Header.Set("Authorization", "Bearer "+client.Token)
	request.Header.Set("Content-Type", "application/json")
	response, err := client.HTTP.Do(request)
	if err != nil {
		return nil, fmt.Errorf("Databricks request failed: %w", err)
	}
	defer response.Body.Close()
	content, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("read Databricks response: %w", err)
	}
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Databricks returned %d: %s", response.StatusCode, string(content))
	}
	return content, nil
}

// CalibrateAdvisoryEdges implements §12.3/§12.5: for every advisory_low_confidence
// entry, score it via I-CALM, gate it against the conformal threshold, and — for
// anything the gate abstains on — attempt corroboration. Never touches
// unverifiable_coverage's baseline verification_path (#35's guarantee holds
// independent of this function ever running at all).
func CalibrateAdvisoryEdges(ctx context.Context, reconciliation *Reconciliation, client DatabricksClient, threshold float64) {
	for i := range reconciliation.AdvisoryLowConfidence {
		entry := &reconciliation.AdvisoryLowConfidence[i]
		confidence, err := client.ScoreEdge(ctx, entry.Reason, nil)
		if err != nil {
			continue // fail closed: leave verbal_confidence/conformal_gate_passed null, exactly as an unconfigured client would
		}
		entry.VerbalConfidence = &confidence
		decision := GateAdvisoryEdge(confidence, threshold)
		entry.ConformalGatePassed = &decision.ConformalGatePassed
		if !decision.PendingVerification {
			continue
		}
		entry.Corroboration.Attempted = true
		verified, snippet, err := client.Corroborate(ctx, entry.Reason)
		if err != nil {
			continue // corroboration attempted, verified stays nil — a failed lookup is not a "no"
		}
		entry.Corroboration.Verified = &verified
		if verified {
			entry.Corroboration.EvidenceSnippet = &snippet
		}
	}
}

// CorroborateUnverifiableCoverage implements §12.15: the Databricks-enhanced
// tier for the mandatory Curveball baseline (#35), never a substitute for
// it. A miss leaves VerificationPath unchanged — the baseline still holds.
func CorroborateUnverifiableCoverage(ctx context.Context, reconciliation *Reconciliation, client DatabricksClient) {
	for i := range reconciliation.UnverifiableCoverage {
		entry := &reconciliation.UnverifiableCoverage[i]
		entry.Corroboration.Attempted = true
		verified, snippet, err := client.Corroborate(ctx, entry.CoverageReason)
		if err != nil {
			continue
		}
		entry.Corroboration.Verified = &verified
		if verified {
			entry.Corroboration.EvidenceSnippet = &snippet
		}
	}
}
