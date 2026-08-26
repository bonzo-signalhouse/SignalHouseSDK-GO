package signalhouse

import "context"

// VoiceTokensService handles WebRTC voice token API calls. Mints ephemeral
// SIP identities the browser/mobile voice SDK uses to register against the
// SignalHouse voice fabric. Wraps voice-backend's /voice/tokens endpoint.
//
// Customer flow:
//  1. Backend calls Tokens.CreateToken({Identity: "alice", TTL: 1800}) and
//     returns the resulting token to the customer's browser/app.
//  2. Browser SDK registers using the embedded SipCredentials.
//  3. Outbound calls from the SDK go through OpenSIPS → Asterisk → carrier;
//     the server-side Calls.CreateCall({ToIdentity: "alice", ...}) flow can
//     also ring this identity for click-to-call patterns.
type VoiceTokensService struct {
	client *Client
}

// VoiceIncomingGrants are inbound call grants on a voice token.
type VoiceIncomingGrants struct {
	// Allow permits inbound calls to this identity.
	Allow *bool `json:"allow,omitempty"`
}

// VoiceOutgoingGrants are outbound call grants on a voice token.
type VoiceOutgoingGrants struct {
	// AllowedNumbers are E.164 numbers this identity may present as caller ID.
	AllowedNumbers []string `json:"allowedNumbers,omitempty"`
}

// VoiceGrantsBlock holds the voice-specific grants attached to a token.
type VoiceGrantsBlock struct {
	Incoming *VoiceIncomingGrants `json:"incoming,omitempty"`
	Outgoing *VoiceOutgoingGrants `json:"outgoing,omitempty"`
}

// VoiceGrants is the capability grants object on a token.
type VoiceGrants struct {
	Voice *VoiceGrantsBlock `json:"voice,omitempty"`
}

// CreateTokenParams is the body for CreateToken. All fields optional.
type CreateTokenParams struct {
	// Identity is the customer-facing logical identity name (e.g. "alice").
	// Persisted as requested_identity so server-initiated
	// Calls.CreateCall({ToIdentity: "alice"}) can resolve it later.
	// Defaults to "anonymous".
	Identity string `json:"identity,omitempty"`
	// TTL is the token TTL in seconds. Clamped to the server's max.
	// Defaults to VOICE_TOKEN_DEFAULT_TTL_SEC.
	TTL int `json:"ttl,omitempty"`
	// Grants are the capability grants attached to the token.
	Grants *VoiceGrants `json:"grants,omitempty"`
}

// CreateToken mints an ephemeral voice token + SIP credentials. The response
// payload includes the JWT, the server-generated SIP username, an expiry
// timestamp, and the embedded SIP credentials (username, password, domain,
// wss_url). The plaintext password is returned ONCE and is not recoverable.
// POST /voice/tokens
func (s *VoiceTokensService) CreateToken(ctx context.Context, params CreateTokenParams, opts *RequestOptions) (*Response, error) {
	return s.client.doJSON(ctx, "POST", "/voice/tokens", params, opts)
}
