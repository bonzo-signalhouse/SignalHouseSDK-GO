package signalhouse

// VoiceService groups voice-backend resources under one namespace. Access
// via client.Voice.SipTrunks, client.Voice.SipProfiles, client.Voice.Calls,
// and client.Voice.Tokens.
//
// Routes live under /voice/* on the voice-backend service.
type VoiceService struct {
	client *Client

	// SipTrunks — peer-to-peer SIP connections to a PBX or carrier.
	SipTrunks *VoiceSipTrunksService

	// SipProfiles — single registerable SIP endpoints (desk phones, softphones, apps).
	SipProfiles *VoiceSipProfilesService

	// ProgrammableVoiceProfiles — route a set of the account's numbers through Signal House's own voice servers.
	ProgrammableVoiceProfiles *VoiceProgrammableVoiceProfilesService

	// Calls — outbound call origination and call log queries.
	Calls *VoiceCallsService

	// CallLogs — account-scoped call history: list/get + presigned recording.
	CallLogs *VoiceCallLogsService

	// Analytics — aggregated call metrics (summary tiles + byDate/byNumber/byCarrier breakdowns).
	Analytics *VoiceAnalyticsService

	// GlobalVoiceSettings — account-wide voice defaults (accepted regions, max spend per minute, E911).
	GlobalVoiceSettings *VoiceGlobalVoiceSettingsService

	// Tokens — mint ephemeral SIP credentials for the browser voice SDK.
	Tokens *VoiceTokensService
}
