import { useState } from "react";

const C = {
  primary: "#4D2975", accent: "#26B8B0", amber: "#E9A800",
  success: "#22C55E", error: "#EF4444", warning: "#F59E0B",
  dark: "#1A1A2E", g900: "#111827", g700: "#374151", g500: "#6B7280",
  g400: "#9CA3AF", g300: "#D1D5DB", g200: "#E5E7EB", g100: "#F3F4F6",
  g50: "#F9FAFB", white: "#FFFFFF",
};

const tabs = [
  { id: "arch", label: "Architecture" },
  { id: "push", label: "Push Flow" },
  { id: "pull", label: "Pull Flow" },
  { id: "endpoint", label: "Endpoint Spec" },
  { id: "scenarios", label: "Scenarios" },
];

function ArchTab() {
  return (
    <div className="space-y-6">
      <div className="text-center mb-4">
        <h2 className="text-xl font-bold" style={{ color: C.primary }}>Vernon = Remote Kill Switch + Config Registry</h2>
        <p className="text-xs mt-1" style={{ color: C.g500 }}>You own the code and servers. No crypto needed. Just on/off.</p>
      </div>

      <svg viewBox="0 0 760 380" className="w-full" style={{ maxWidth: 760 }}>
        <defs>
          <marker id="a1" viewBox="0 0 10 10" refX="9" refY="5" markerWidth="6" markerHeight="6" orient="auto"><path d="M 0 0 L 10 5 L 0 10 z" fill={C.accent} /></marker>
          <marker id="a2" viewBox="0 0 10 10" refX="9" refY="5" markerWidth="6" markerHeight="6" orient="auto"><path d="M 0 0 L 10 5 L 0 10 z" fill={C.primary} /></marker>
          <marker id="a3" viewBox="0 0 10 10" refX="9" refY="5" markerWidth="6" markerHeight="6" orient="auto"><path d="M 0 0 L 10 5 L 0 10 z" fill={C.error} /></marker>
          <marker id="a4" viewBox="0 0 10 10" refX="9" refY="5" markerWidth="6" markerHeight="6" orient="auto"><path d="M 0 0 L 10 5 L 0 10 z" fill={C.g300} /></marker>
        </defs>

        {/* Vernon */}
        <rect x="250" y="20" width="260" height="100" rx="16" fill={C.primary} />
        <text x="380" y="52" textAnchor="middle" fill="white" fontWeight="800" fontSize="16">Vernon License</text>
        <text x="380" y="72" textAnchor="middle" fill="rgba(255,255,255,0.7)" fontSize="11">Source of truth for all client apps</text>
        <text x="380" y="92" textAnchor="middle" fill="rgba(255,255,255,0.5)" fontSize="10">status: active | suspended | expired</text>

        {/* PO action */}
        <rect x="30" y="40" width="170" height="60" rx="10" fill={C.amber + "15"} stroke={C.amber} strokeWidth="1.5" />
        <text x="115" y="65" textAnchor="middle" fill={C.amber} fontWeight="700" fontSize="11">Project Owner</text>
        <text x="115" y="80" textAnchor="middle" fill={C.g500} fontSize="9">toggles on/off</text>
        <path d="M200 70 L250 70" stroke={C.amber} strokeWidth="2" markerEnd="url(#a2)" />

        {/* Push arrow */}
        <path d="M380 120 L380 165" stroke={C.accent} strokeWidth="2.5" markerEnd="url(#a1)" />
        <rect x="310" y="135" width="140" height="22" rx="4" fill={C.accent + "15"} />
        <text x="380" y="150" textAnchor="middle" fill={C.accent} fontWeight="700" fontSize="9">① PUSH on change</text>

        {/* Client apps area */}
        <rect x="80" y="175" width="600" height="180" rx="14" fill={C.g50} stroke={C.g200} strokeWidth="1.5" strokeDasharray="5,3" />
        <text x="380" y="197" textAnchor="middle" fill={C.g700} fontWeight="700" fontSize="12">Your Servers (you control everything)</text>

        {/* Client A - active */}
        <rect x="100" y="215" width="160" height="80" rx="10" fill={C.success + "10"} stroke={C.success} strokeWidth="1.5" />
        <circle cx="120" cy="235" r="6" fill={C.success} />
        <text x="140" y="239" fill={C.success} fontWeight="700" fontSize="11">Client App A</text>
        <text x="140" y="256" fill={C.g500} fontSize="9">status: active ✓</text>
        <text x="140" y="272" fill={C.g400} fontSize="8">app running normally</text>
        <text x="140" y="286" fill={C.g400} fontSize="8">users have full access</text>

        {/* Client B - active */}
        <rect x="300" y="215" width="160" height="80" rx="10" fill={C.success + "10"} stroke={C.success} strokeWidth="1.5" />
        <circle cx="320" cy="235" r="6" fill={C.success} />
        <text x="340" y="239" fill={C.success} fontWeight="700" fontSize="11">Client App B</text>
        <text x="340" y="256" fill={C.g500} fontSize="9">status: active ✓</text>
        <text x="340" y="272" fill={C.g400} fontSize="8">app running normally</text>

        {/* Client C - killed */}
        <rect x="500" y="215" width="160" height="80" rx="10" fill={C.error + "10"} stroke={C.error} strokeWidth="1.5" />
        <circle cx="520" cy="235" r="6" fill={C.error} />
        <text x="540" y="239" fill={C.error} fontWeight="700" fontSize="11">Client App C</text>
        <text x="540" y="256" fill={C.g500} fontSize="9">status: suspended ✗</text>
        <text x="540" y="272" fill={C.error} fontSize="8" fontWeight="600">ACCESS DENIED</text>
        <text x="540" y="286" fill={C.g400} fontSize="8">maintenance page shown</text>

        {/* Pull arrows */}
        <path d="M180 215 Q180 175 320 140" stroke={C.g300} strokeWidth="1" fill="none" markerEnd="url(#a4)" strokeDasharray="3,3" />
        <path d="M380 215 Q380 175 380 140" stroke={C.g300} strokeWidth="1" fill="none" markerEnd="url(#a4)" strokeDasharray="3,3" />
        <path d="M580 215 Q580 175 440 140" stroke={C.g300} strokeWidth="1" fill="none" markerEnd="url(#a4)" strokeDasharray="3,3" />
        <text x="600" y="178" fill={C.g400} fontSize="8">② pull backup (every 6h)</text>

        {/* Legend */}
        <rect x="100" y="325" width="460" height="35" rx="6" fill={C.white} stroke={C.g200} strokeWidth="1" />
        <circle cx="125" cy="342" r="5" fill={C.accent} /><text x="135" y="346" fill={C.g600} fontSize="9">Push: Vernon calls client on status change</text>
        <circle cx="355" cy="342" r="5" fill={C.g300} /><text x="365" y="346" fill={C.g600} fontSize="9">Pull: Client checks Vernon periodically</text>
      </svg>

      {/* Key points */}
      <div className="grid grid-cols-3 gap-3">
        <div className="rounded-xl p-4 border" style={{ background: C.accent + "08", borderColor: C.accent + "30" }}>
          <div className="text-sm font-bold mb-1" style={{ color: C.accent }}>Dead Simple</div>
          <p className="text-xs" style={{ color: C.g500 }}>No JWT, no crypto, no signatures. Just an HTTP call that returns <code className="px-1 py-0.5 rounded text-xs" style={{ background: C.g100 }}>{"{ active: true }"}</code> or <code className="px-1 py-0.5 rounded text-xs" style={{ background: C.g100 }}>{"{ active: false }"}</code></p>
        </div>
        <div className="rounded-xl p-4 border" style={{ background: C.primary + "08", borderColor: C.primary + "30" }}>
          <div className="text-sm font-bold mb-1" style={{ color: C.primary }}>You Own It All</div>
          <p className="text-xs" style={{ color: C.g500 }}>You control the code, the servers, and the network. No need to protect against your own infra. Tamper concern = zero.</p>
        </div>
        <div className="rounded-xl p-4 border" style={{ background: C.amber + "08", borderColor: C.amber + "30" }}>
          <div className="text-sm font-bold mb-1" style={{ color: C.amber }}>Instant Kill Switch</div>
          <p className="text-xs" style={{ color: C.g500 }}>PO suspends license → Vernon pushes to client → client shows maintenance page. Takes effect in seconds, not hours.</p>
        </div>
      </div>
    </div>
  );
}

function PushTab() {
  const steps = [
    { actor: "Vernon", label: "Status changes", desc: "Project Owner suspends a license via Vernon app", color: C.amber, icon: "🔄" },
    { actor: "Vernon", label: "Push to client", desc: "POST {instance_url}/api/v1/vernon/sync with new status", color: C.primary, icon: "📡" },
    { actor: "Client", label: "Receives push", desc: "Client app's /vernon/sync endpoint receives the call", color: C.accent, icon: "📥" },
    { actor: "Client", label: "Updates local state", desc: "Writes new status to local config/DB. If suspended → immediately block users", color: C.error, icon: "⚡" },
    { actor: "Client", label: "Returns 200 OK", desc: "Vernon marks push as delivered. If client unreachable → retry 3x then notify PO", color: C.success, icon: "✓" },
  ];
  return (
    <div className="space-y-6">
      <div className="text-center mb-4">
        <h2 className="text-xl font-bold" style={{ color: C.primary }}>Push Flow — Instant Toggle</h2>
        <p className="text-xs mt-1" style={{ color: C.g500 }}>Vernon calls the client app immediately when status changes</p>
      </div>
      <div className="space-y-2">
        {steps.map((s, i) => (
          <div key={i} className="flex items-start gap-3">
            <div className="flex flex-col items-center" style={{ minWidth: 36 }}>
              <div className="w-9 h-9 rounded-full flex items-center justify-center text-base" style={{ background: s.color + "15", border: `2px solid ${s.color}` }}>{s.icon}</div>
              {i < steps.length - 1 && <div className="w-0.5 h-4" style={{ background: C.g200 }} />}
            </div>
            <div className="flex-1 rounded-xl p-3 border" style={{ background: s.color + "06", borderColor: s.color + "25" }}>
              <div className="flex items-center gap-2 mb-0.5">
                <span className="text-xs font-bold px-2 py-0.5 rounded-full" style={{ background: s.color + "15", color: s.color, fontSize: 10 }}>{s.actor}</span>
                <span className="text-sm font-bold" style={{ color: s.color }}>{s.label}</span>
              </div>
              <p className="text-xs" style={{ color: C.g500 }}>{s.desc}</p>
            </div>
          </div>
        ))}
      </div>
      <div className="rounded-xl p-4 border" style={{ background: C.g50, borderColor: C.g200 }}>
        <div className="text-xs font-bold mb-2" style={{ color: C.g700 }}>Push Payload</div>
        <pre className="text-xs p-3 rounded-lg overflow-x-auto" style={{ background: C.dark, color: C.accent, fontFamily: "monospace" }}>{`POST {instance_url}/api/v1/vernon/sync
Header: X-API-Key: {client_registration_code}

{
  "license_key": "FL-A1B2C3D4",
  "active": false,
  "reason": "suspended",
  "updated_at": "2026-03-24T10:00:00Z"
}`}</pre>
      </div>
    </div>
  );
}

function PullTab() {
  const steps = [
    { actor: "Client", label: "Periodic timer fires", desc: "Every check_interval (default 6h), client calls Vernon", color: C.accent, icon: "⏰" },
    { actor: "Client", label: "GET /client/license", desc: "Sends license_key, Vernon returns current status", color: C.primary, icon: "📡" },
    { actor: "Vernon", label: "Returns status", desc: "Simple JSON: active true/false + reason if suspended", color: C.primary, icon: "📦" },
    { actor: "Client", label: "Applies status", desc: "If active → continue. If not → block access. Cache result locally.", color: C.success, icon: "✓" },
    { actor: "Client", label: "Vernon unreachable?", desc: "Use cached status. If unreachable for 3× interval → keep last known status (you own the server, no bypass risk)", color: C.amber, icon: "📴" },
  ];
  return (
    <div className="space-y-6">
      <div className="text-center mb-4">
        <h2 className="text-xl font-bold" style={{ color: C.primary }}>Pull Flow — Periodic Backup Check</h2>
        <p className="text-xs mt-1" style={{ color: C.g500 }}>Client checks Vernon periodically as a safety net (in case push was missed)</p>
      </div>
      <div className="space-y-2">
        {steps.map((s, i) => (
          <div key={i} className="flex items-start gap-3">
            <div className="flex flex-col items-center" style={{ minWidth: 36 }}>
              <div className="w-9 h-9 rounded-full flex items-center justify-center text-base" style={{ background: s.color + "15", border: `2px solid ${s.color}` }}>{s.icon}</div>
              {i < steps.length - 1 && <div className="w-0.5 h-4" style={{ background: C.g200 }} />}
            </div>
            <div className="flex-1 rounded-xl p-3 border" style={{ background: s.color + "06", borderColor: s.color + "25" }}>
              <div className="flex items-center gap-2 mb-0.5">
                <span className="text-xs font-bold px-2 py-0.5 rounded-full" style={{ background: s.color + "15", color: s.color, fontSize: 10 }}>{s.actor}</span>
                <span className="text-sm font-bold" style={{ color: s.color }}>{s.label}</span>
              </div>
              <p className="text-xs" style={{ color: C.g500 }}>{s.desc}</p>
            </div>
          </div>
        ))}
      </div>
      <div className="rounded-xl p-4 border" style={{ background: C.g50, borderColor: C.g200 }}>
        <div className="text-xs font-bold mb-2" style={{ color: C.g700 }}>Pull Response</div>
        <pre className="text-xs p-3 rounded-lg overflow-x-auto" style={{ background: C.dark, color: C.accent, fontFamily: "monospace" }}>{`GET /api/v1/client/license?key=FL-A1B2C3D4

{
  "active": true,
  "license_key": "FL-A1B2C3D4",
  "check_interval": "6h",
  "updated_at": "2026-03-24T08:00:00Z"
}`}</pre>
      </div>
    </div>
  );
}

function EndpointTab() {
  return (
    <div className="space-y-6">
      <div className="text-center mb-4">
        <h2 className="text-xl font-bold" style={{ color: C.primary }}>Endpoint Spec</h2>
        <p className="text-xs mt-1" style={{ color: C.g500 }}>Two endpoints — one on Vernon, one on each client app</p>
      </div>

      <div className="grid grid-cols-2 gap-4">
        <div className="rounded-xl border overflow-hidden" style={{ borderColor: C.primary + "30" }}>
          <div className="p-3" style={{ background: C.primary, color: C.white }}>
            <div className="text-xs font-bold">Vernon API (pull endpoint)</div>
            <div className="text-xs opacity-70 mt-0.5">Client calls this periodically</div>
          </div>
          <div className="p-4 space-y-3">
            <pre className="text-xs p-2 rounded" style={{ background: C.g50, color: C.g700, fontFamily: "monospace" }}>{`GET /api/v1/client/license
  ?key=FL-XXXXXXXX`}</pre>
            <div className="text-xs font-bold" style={{ color: C.g700 }}>Response:</div>
            <pre className="text-xs p-2 rounded" style={{ background: C.g50, color: C.g700, fontFamily: "monospace" }}>{`{
  "active": true,
  "license_key": "FL-A1B2C3D4",
  "check_interval": "6h",
  "updated_at": "2026-03-24T..."
}`}</pre>
            <div className="text-xs" style={{ color: C.g500 }}>No auth needed. Just license_key as identifier. Origin header check optional (extra safety for SaaS).</div>
          </div>
        </div>

        <div className="rounded-xl border overflow-hidden" style={{ borderColor: C.accent + "30" }}>
          <div className="p-3" style={{ background: C.accent, color: C.white }}>
            <div className="text-xs font-bold">Client App (push endpoint)</div>
            <div className="text-xs opacity-70 mt-0.5">Vernon calls this on status change</div>
          </div>
          <div className="p-4 space-y-3">
            <pre className="text-xs p-2 rounded" style={{ background: C.g50, color: C.g700, fontFamily: "monospace" }}>{`POST /api/v1/vernon/sync
Header: X-API-Key: {key}
`}</pre>
            <div className="text-xs font-bold" style={{ color: C.g700 }}>Body:</div>
            <pre className="text-xs p-2 rounded" style={{ background: C.g50, color: C.g700, fontFamily: "monospace" }}>{`{
  "license_key": "FL-A1B2C3D4",
  "active": false,
  "reason": "suspended",
  "updated_at": "2026-03-24T..."
}`}</pre>
            <div className="text-xs" style={{ color: C.g500 }}>Secured by client_registration_code. Client must implement this endpoint.</div>
          </div>
        </div>
      </div>

      <div className="rounded-xl p-4 border" style={{ background: C.amber + "08", borderColor: C.amber + "30" }}>
        <div className="text-xs font-bold mb-2" style={{ color: C.amber }}>What got removed vs the old design</div>
        <div className="grid grid-cols-2 gap-3">
          <div>
            <div className="text-xs font-bold mb-1" style={{ color: C.error }}>Removed ✗</div>
            <div className="space-y-1">
              {["JWT signing/verification", "RS256 key pairs", "Environment fingerprinting", "Modules/apps/constraints in response", "Grace period deny logic", "Origin header validation (optional now)"].map((x, i) => (
                <div key={i} className="text-xs flex items-center gap-1" style={{ color: C.g500 }}>
                  <span style={{ color: C.error }}>✗</span> {x}
                </div>
              ))}
            </div>
          </div>
          <div>
            <div className="text-xs font-bold mb-1" style={{ color: C.success }}>Kept ✓</div>
            <div className="space-y-1">
              {["License key as identifier", "check_interval (configurable)", "client_registration_code for push auth", "Push on change (instant)", "Pull as backup (periodic)", "Retry + notification on push failure"].map((x, i) => (
                <div key={i} className="text-xs flex items-center gap-1" style={{ color: C.g500 }}>
                  <span style={{ color: C.success }}>✓</span> {x}
                </div>
              ))}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

function ScenariosTab() {
  const scenarios = [
    {
      title: "Client doesn't pay → suspend",
      steps: ["PO opens Vernon, clicks Suspend on license", "Vernon updates status to 'suspended'", "Vernon POSTs to client's /vernon/sync: { active: false, reason: 'suspended' }", "Client app immediately shows 'Akun ditangguhkan, hubungi admin'", "Users can't access anything"],
      color: C.error,
    },
    {
      title: "Client pays → reactivate",
      steps: ["PO clicks Activate on the license", "Vernon pushes { active: true } to client", "Client app instantly resumes normal operation", "Users regain full access"],
      color: C.success,
    },
    {
      title: "License expires naturally",
      steps: ["Vernon cron detects expires_at < now", "Auto-changes status to 'expired'", "Pushes { active: false, reason: 'expired' } to client", "Client shows 'Lisensi kedaluwarsa'", "PO gets notification"],
      color: C.amber,
    },
    {
      title: "Push fails (client unreachable)",
      steps: ["Vernon tries to push → connection refused", "Retry 3 times (1min, 5min, 15min)", "All retries fail → send notification to PO: 'Push failed for FL-XXX'", "Client will still catch it on next periodic pull (within check_interval)", "No data loss — Vernon is the source of truth"],
      color: C.primary,
    },
    {
      title: "Vernon is temporarily down",
      steps: ["Client's periodic pull fails", "Client keeps last known status (cached locally)", "Since you own the server, Vernon being down = your infra issue, not a bypass attack", "Fix Vernon, client resumes normal checks on next interval", "No deny — trusted environment"],
      color: C.g500,
    },
  ];

  return (
    <div className="space-y-6">
      <div className="text-center mb-4">
        <h2 className="text-xl font-bold" style={{ color: C.primary }}>Real-world Scenarios</h2>
      </div>
      <div className="space-y-4">
        {scenarios.map((s, i) => (
          <div key={i} className="rounded-xl border overflow-hidden" style={{ borderColor: s.color + "30" }}>
            <div className="px-4 py-3 flex items-center gap-2" style={{ background: s.color + "10" }}>
              <div className="w-6 h-6 rounded-full flex items-center justify-center text-white text-xs font-bold" style={{ background: s.color }}>{i + 1}</div>
              <h3 className="text-sm font-bold" style={{ color: s.color }}>{s.title}</h3>
            </div>
            <div className="p-4">
              <div className="space-y-2">
                {s.steps.map((step, j) => (
                  <div key={j} className="flex items-start gap-2">
                    <div className="w-5 h-5 rounded-full flex items-center justify-center text-xs font-bold flex-shrink-0 mt-0.5" style={{ background: C.g100, color: C.g500 }}>{j + 1}</div>
                    <p className="text-xs" style={{ color: C.g600 }}>{step}</p>
                  </div>
                ))}
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}

export default function App() {
  const [tab, setTab] = useState("arch");
  const render = () => {
    switch (tab) {
      case "arch": return <ArchTab />;
      case "push": return <PushTab />;
      case "pull": return <PullTab />;
      case "endpoint": return <EndpointTab />;
      case "scenarios": return <ScenariosTab />;
    }
  };

  return (
    <div className="min-h-screen" style={{ background: C.g100 }}>
      <div className="max-w-4xl mx-auto p-6">
        <div className="flex items-center gap-3 mb-5">
          <div className="w-10 h-10 rounded-xl flex items-center justify-center text-white font-bold text-sm" style={{ background: C.primary }}>VL</div>
          <div>
            <h1 className="text-lg font-bold" style={{ color: C.primary }}>License Validation — Simplified</h1>
            <p className="text-xs" style={{ color: C.g500 }}>Push on change + periodic pull. Just on/off. You own the infra.</p>
          </div>
        </div>

        <div className="flex gap-1 mb-5 p-1 rounded-xl" style={{ background: C.g200 }}>
          {tabs.map((t) => (
            <button key={t.id} onClick={() => setTab(t.id)}
              className="flex-1 py-2 px-2 rounded-lg text-xs font-semibold transition-all"
              style={{ background: tab === t.id ? C.white : "transparent", color: tab === t.id ? C.primary : C.g500, boxShadow: tab === t.id ? "0 1px 3px rgba(0,0,0,0.1)" : "none" }}>
              {t.label}
            </button>
          ))}
        </div>

        <div className="rounded-2xl p-6" style={{ background: C.white, boxShadow: "0 1px 3px rgba(0,0,0,0.08)" }}>
          {render()}
        </div>
      </div>
    </div>
  );
}
