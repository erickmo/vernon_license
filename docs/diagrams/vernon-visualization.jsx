import { useState } from "react";

const COLORS = {
  primary: "#4D2975",
  primaryLight: "#6B3FA0",
  primaryDark: "#3A1D5C",
  accent: "#26B8B0",
  accentLight: "#E8F8F7",
  amber: "#E9A800",
  amberLight: "#FFF8E1",
  success: "#22C55E",
  successLight: "#F0FDF4",
  error: "#EF4444",
  errorLight: "#FEF2F2",
  warning: "#F59E0B",
  dark: "#1A1A2E",
  gray900: "#111827",
  gray700: "#374151",
  gray500: "#6B7280",
  gray400: "#9CA3AF",
  gray300: "#D1D5DB",
  gray200: "#E5E7EB",
  gray100: "#F3F4F6",
  gray50: "#F9FAFB",
  white: "#FFFFFF",
};

const tabs = [
  { id: "overview", label: "System Overview" },
  { id: "entities", label: "Entity Map" },
  { id: "proposal", label: "Proposal Flow" },
  { id: "license", label: "License Lifecycle" },
  { id: "check", label: "License Check" },
  { id: "roles", label: "Role Policies" },
];

// === OVERVIEW TAB ===
function OverviewTab() {
  return (
    <div className="space-y-8">
      <div className="text-center mb-8">
        <h2 className="text-2xl font-bold" style={{ color: COLORS.primary }}>Vernon License — System Overview</h2>
        <p className="text-sm mt-1" style={{ color: COLORS.gray500 }}>Centralized licensing for all Vernon client apps</p>
      </div>

      <svg viewBox="0 0 900 520" className="w-full" style={{ maxWidth: 900 }}>
        {/* Vernon License System */}
        <rect x="280" y="20" width="340" height="80" rx="16" fill={COLORS.primary} />
        <text x="450" y="52" textAnchor="middle" fill="white" fontWeight="700" fontSize="16">Vernon License</text>
        <text x="450" y="72" textAnchor="middle" fill="rgba(255,255,255,0.7)" fontSize="11">Centralized Licensing System</text>
        <text x="450" y="88" textAnchor="middle" fill="rgba(255,255,255,0.5)" fontSize="10">API + Go WASM PWA</text>

        {/* Users */}
        <rect x="30" y="170" width="140" height="56" rx="10" fill={COLORS.accentLight} stroke={COLORS.accent} strokeWidth="1.5" />
        <text x="100" y="194" textAnchor="middle" fill={COLORS.accent} fontWeight="700" fontSize="12">Superuser</text>
        <text x="100" y="210" textAnchor="middle" fill={COLORS.gray500} fontSize="9">Products, Users, Audit</text>

        <rect x="30" y="240" width="140" height="56" rx="10" fill={COLORS.amberLight} stroke={COLORS.amber} strokeWidth="1.5" />
        <text x="100" y="264" textAnchor="middle" fill={COLORS.amber} fontWeight="700" fontSize="12">Project Owner</text>
        <text x="100" y="280" textAnchor="middle" fill={COLORS.gray500} fontSize="9">Approve, License, Provision</text>

        <rect x="30" y="310" width="140" height="56" rx="10" fill="#EFF6FF" stroke="#3B82F6" strokeWidth="1.5" />
        <text x="100" y="334" textAnchor="middle" fill="#3B82F6" fontWeight="700" fontSize="12">Sales</text>
        <text x="100" y="350" textAnchor="middle" fill={COLORS.gray500} fontSize="9">Proposals, Prospecting</text>

        {/* Arrow: Users → Vernon */}
        <path d="M170 268 L280 60" stroke={COLORS.gray300} strokeWidth="1.5" fill="none" markerEnd="url(#arrowGray)" />
        <text x="210" y="155" textAnchor="middle" fill={COLORS.gray400} fontSize="9" transform="rotate(-40, 210, 155)">manage via PWA</text>

        {/* Products box */}
        <rect x="290" y="140" width="130" height="50" rx="8" fill={COLORS.gray50} stroke={COLORS.gray200} strokeWidth="1" />
        <text x="355" y="162" textAnchor="middle" fill={COLORS.gray700} fontWeight="600" fontSize="11">Products</text>
        <text x="355" y="178" textAnchor="middle" fill={COLORS.gray400} fontSize="9">Dynamic — any product</text>

        {/* Companies box */}
        <rect x="440" y="140" width="130" height="50" rx="8" fill={COLORS.gray50} stroke={COLORS.gray200} strokeWidth="1" />
        <text x="505" y="162" textAnchor="middle" fill={COLORS.gray700} fontWeight="600" fontSize="11">Companies</text>
        <text x="505" y="178" textAnchor="middle" fill={COLORS.gray400} fontSize="9">PT Maju, CV Jaya...</text>

        {/* Projects */}
        <rect x="365" y="210" width="170" height="50" rx="8" fill={COLORS.amberLight} stroke={COLORS.amber} strokeWidth="1" />
        <text x="450" y="232" textAnchor="middle" fill={COLORS.amber} fontWeight="600" fontSize="11">Projects</text>
        <text x="450" y="248" textAnchor="middle" fill={COLORS.gray500} fontSize="9">Groups licenses + proposals</text>

        {/* Two paths */}
        <rect x="280" y="290" width="150" height="55" rx="10" fill="#EFF6FF" stroke="#3B82F6" strokeWidth="1.5" />
        <text x="355" y="312" textAnchor="middle" fill="#3B82F6" fontWeight="700" fontSize="11">Proposal Path</text>
        <text x="355" y="328" textAnchor="middle" fill={COLORS.gray500} fontSize="9">Sales → PO approve</text>

        <rect x="470" y="290" width="150" height="55" rx="10" fill={COLORS.amberLight} stroke={COLORS.amber} strokeWidth="1.5" />
        <text x="545" y="312" textAnchor="middle" fill={COLORS.amber} fontWeight="700" fontSize="11">Direct Path</text>
        <text x="545" y="328" textAnchor="middle" fill={COLORS.gray500} fontSize="9">PO creates directly</text>

        {/* License */}
        <rect x="350" y="375" width="200" height="60" rx="12" fill={COLORS.successLight} stroke={COLORS.success} strokeWidth="2" />
        <text x="450" y="400" textAnchor="middle" fill="#166534" fontWeight="700" fontSize="13">License</text>
        <text x="450" y="418" textAnchor="middle" fill={COLORS.gray500} fontSize="9">FL-XXXXXXXX · constraints · modules</text>

        {/* Arrows to License */}
        <line x1="355" y1="345" x2="420" y2="375" stroke="#3B82F6" strokeWidth="1.5" markerEnd="url(#arrowBlue)" />
        <line x1="545" y1="345" x2="480" y2="375" stroke={COLORS.amber} strokeWidth="1.5" markerEnd="url(#arrowAmber)" />

        {/* Provision arrow */}
        <path d="M550 405 Q650 405 700 405" stroke={COLORS.success} strokeWidth="2" fill="none" markerEnd="url(#arrowGreen)" strokeDasharray="6,3" />
        <text x="625" y="395" textAnchor="middle" fill={COLORS.success} fontSize="9" fontWeight="600">provision</text>

        {/* Client Apps — any Vernon product deployment */}
        <rect x="710" y="140" width="170" height="320" rx="12" fill={COLORS.gray50} stroke={COLORS.gray200} strokeWidth="1.5" strokeDasharray="4,3" />
        <text x="795" y="165" textAnchor="middle" fill={COLORS.gray700} fontWeight="700" fontSize="12">Any Client App</text>
        <text x="795" y="180" textAnchor="middle" fill={COLORS.gray400} fontSize="9">deployed Vernon products</text>

        <rect x="725" y="195" width="140" height="36" rx="6" fill="white" stroke={COLORS.accent} strokeWidth="1" />
        <text x="795" y="218" textAnchor="middle" fill={COLORS.accent} fontSize="10" fontWeight="600">Client App A</text>

        <rect x="725" y="240" width="140" height="36" rx="6" fill="white" stroke={COLORS.accent} strokeWidth="1" />
        <text x="795" y="263" textAnchor="middle" fill={COLORS.accent} fontSize="10" fontWeight="600">Client App B</text>

        <rect x="725" y="285" width="140" height="36" rx="6" fill="white" stroke={COLORS.accent} strokeWidth="1" />
        <text x="795" y="308" textAnchor="middle" fill={COLORS.accent} fontSize="10" fontWeight="600">Client App C</text>

        <text x="795" y="340" textAnchor="middle" fill={COLORS.gray400} fontSize="18">···</text>
        <text x="795" y="360" textAnchor="middle" fill={COLORS.gray400} fontSize="9">unlimited deployments</text>

        {/* Periodic check arrow */}
        <path d="M795 365 Q795 440 550 405" stroke={COLORS.accent} strokeWidth="1.5" fill="none" markerEnd="url(#arrowTeal)" strokeDasharray="4,3" />
        <text x="700" y="445" textAnchor="middle" fill={COLORS.accent} fontSize="9" fontWeight="600">periodic check (1h/6h/24h)</text>

        {/* Deny access */}
        <rect x="725" y="375" width="140" height="40" rx="6" fill={COLORS.errorLight} stroke={COLORS.error} strokeWidth="1" />
        <text x="795" y="393" textAnchor="middle" fill={COLORS.error} fontWeight="600" fontSize="9">If invalid →</text>
        <text x="795" y="406" textAnchor="middle" fill={COLORS.error} fontWeight="700" fontSize="10">DENY ACCESS</text>

        {/* Arrow definitions */}
        <defs>
          <marker id="arrowGray" viewBox="0 0 10 10" refX="9" refY="5" markerWidth="6" markerHeight="6" orient="auto"><path d="M 0 0 L 10 5 L 0 10 z" fill={COLORS.gray300} /></marker>
          <marker id="arrowBlue" viewBox="0 0 10 10" refX="9" refY="5" markerWidth="6" markerHeight="6" orient="auto"><path d="M 0 0 L 10 5 L 0 10 z" fill="#3B82F6" /></marker>
          <marker id="arrowAmber" viewBox="0 0 10 10" refX="9" refY="5" markerWidth="6" markerHeight="6" orient="auto"><path d="M 0 0 L 10 5 L 0 10 z" fill={COLORS.amber} /></marker>
          <marker id="arrowGreen" viewBox="0 0 10 10" refX="9" refY="5" markerWidth="6" markerHeight="6" orient="auto"><path d="M 0 0 L 10 5 L 0 10 z" fill={COLORS.success} /></marker>
          <marker id="arrowTeal" viewBox="0 0 10 10" refX="9" refY="5" markerWidth="6" markerHeight="6" orient="auto"><path d="M 0 0 L 10 5 L 0 10 z" fill={COLORS.accent} /></marker>
        </defs>
      </svg>
    </div>
  );
}

// === ENTITIES TAB ===
function EntitiesTab() {
  const entities = [
    { name: "Company", color: "#3B82F6", bg: "#EFF6FF", desc: "Owns licenses", fields: "name, email, PIC, address" },
    { name: "Project", color: COLORS.amber, bg: COLORS.amberLight, desc: "Groups licenses + proposals", fields: "company_id, name, status" },
    { name: "Product", color: COLORS.primary, bg: "#F3E8FF", desc: "Dynamic — superuser adds any product", fields: "name, slug, modules, apps, pricing" },
    { name: "Proposal", color: "#8B5CF6", bg: "#F5F3FF", desc: "Versioned with changelog", fields: "project, product, version, status, modules, pricing" },
    { name: "License", color: COLORS.success, bg: COLORS.successLight, desc: "Active deployment", fields: "key, project, product, constraints, instance_url" },
    { name: "User", color: COLORS.gray700, bg: COLORS.gray100, desc: "superuser / project_owner / sales", fields: "name, email, role, is_active" },
  ];

  const relations = [
    { from: "Company", to: "Project", label: "1 → N", y: 0 },
    { from: "Project", to: "License", label: "1 → N", y: 1 },
    { from: "Project", to: "Proposal", label: "1 → N", y: 1 },
    { from: "Product", to: "License", label: "1 → N", y: 2 },
    { from: "Product", to: "Proposal", label: "1 → N", y: 2 },
    { from: "Proposal", to: "License", label: "approve →\nauto-create", y: 3 },
  ];

  return (
    <div className="space-y-6">
      <div className="text-center mb-6">
        <h2 className="text-2xl font-bold" style={{ color: COLORS.primary }}>Entity Relationship Map</h2>
      </div>

      <div className="grid grid-cols-3 gap-3">
        {entities.map((e) => (
          <div key={e.name} className="rounded-xl p-4 border" style={{ background: e.bg, borderColor: e.color + "40" }}>
            <div className="flex items-center gap-2 mb-1">
              <div className="w-3 h-3 rounded-full" style={{ background: e.color }} />
              <span className="font-bold text-sm" style={{ color: e.color }}>{e.name}</span>
            </div>
            <p className="text-xs mb-2" style={{ color: COLORS.gray500 }}>{e.desc}</p>
            <p className="text-xs font-mono px-2 py-1 rounded" style={{ background: "rgba(0,0,0,0.04)", color: COLORS.gray600 }}>{e.fields}</p>
          </div>
        ))}
      </div>

      <div className="rounded-xl p-5 border" style={{ borderColor: COLORS.gray200, background: COLORS.gray50 }}>
        <h3 className="font-bold text-sm mb-3" style={{ color: COLORS.gray700 }}>Relationships</h3>
        <div className="space-y-2">
          {relations.map((r, i) => (
            <div key={i} className="flex items-center gap-3 text-sm">
              <span className="font-mono font-bold px-2 py-0.5 rounded" style={{ background: COLORS.white, color: COLORS.primary, fontSize: 12 }}>{r.from}</span>
              <span style={{ color: COLORS.accent, fontWeight: 700, fontSize: 12 }}>{r.label.replace("\n", " ")}</span>
              <span className="font-mono font-bold px-2 py-0.5 rounded" style={{ background: COLORS.white, color: COLORS.success, fontSize: 12 }}>{r.to}</span>
            </div>
          ))}
        </div>
      </div>

      <div className="rounded-xl p-5 border" style={{ borderColor: COLORS.amber + "40", background: COLORS.amberLight }}>
        <h3 className="font-bold text-sm mb-2" style={{ color: COLORS.amber }}>Two Paths to Create License</h3>
        <div className="grid grid-cols-2 gap-4 mt-3">
          <div className="rounded-lg p-3 bg-white border" style={{ borderColor: "#3B82F640" }}>
            <div className="text-xs font-bold mb-1" style={{ color: "#3B82F6" }}>Path 1: Via Proposal</div>
            <p className="text-xs" style={{ color: COLORS.gray500 }}>Sales → Draft → Submit → PO edits/approves → License auto-created</p>
          </div>
          <div className="rounded-lg p-3 bg-white border" style={{ borderColor: COLORS.amber + "40" }}>
            <div className="text-xs font-bold mb-1" style={{ color: COLORS.amber }}>Path 2: Direct</div>
            <p className="text-xs" style={{ color: COLORS.gray500 }}>Project Owner creates license directly (skips proposal)</p>
          </div>
        </div>
      </div>
    </div>
  );
}

// === PROPOSAL FLOW TAB ===
function ProposalFlowTab() {
  const steps = [
    { actor: "Sales", action: "Create Draft", icon: "📝", color: "#3B82F6", bg: "#EFF6FF", desc: "Select product, modules, apps, constraints, pricing. Auto-calc from base_pricing." },
    { actor: "Sales", action: "Submit for Review", icon: "📤", color: "#8B5CF6", bg: "#F5F3FF", desc: "Draft → Submitted. Project Owner gets notification." },
    { actor: "Project Owner", action: "Review Changelog", icon: "🔍", color: COLORS.amber, bg: COLORS.amberLight, desc: "v2+: Sees exactly what changed (added/removed/changed). Auto-focus changelog tab." },
    { actor: "Project Owner", action: "Edit (optional)", icon: "✏️", color: COLORS.amber, bg: COLORS.amberLight, desc: "CAN modify pricing, modules, notes directly on submitted proposal. Tracked in audit." },
    { actor: "Project Owner", action: "Approve", icon: "✅", color: COLORS.success, bg: COLORS.successLight, desc: "→ Generate PDF → Auto-create/update License → Notify sales" },
    { actor: "Project Owner", action: "OR Reject", icon: "❌", color: COLORS.error, bg: COLORS.errorLight, desc: "With reason. Sales notified. Sales creates new version (v2, v3...)." },
  ];

  return (
    <div className="space-y-6">
      <div className="text-center mb-6">
        <h2 className="text-2xl font-bold" style={{ color: COLORS.primary }}>Proposal Negotiation Flow</h2>
        <p className="text-xs mt-1" style={{ color: COLORS.gray500 }}>Versioned proposals with auto-computed changelog for efficient review</p>
      </div>

      <div className="space-y-3">
        {steps.map((s, i) => (
          <div key={i} className="flex items-start gap-4">
            <div className="flex flex-col items-center" style={{ minWidth: 40 }}>
              <div className="w-10 h-10 rounded-full flex items-center justify-center text-lg" style={{ background: s.bg, border: `2px solid ${s.color}` }}>{s.icon}</div>
              {i < steps.length - 1 && <div className="w-0.5 h-6" style={{ background: COLORS.gray200 }} />}
            </div>
            <div className="flex-1 rounded-xl p-4 border" style={{ background: s.bg, borderColor: s.color + "30" }}>
              <div className="flex items-center gap-2 mb-1">
                <span className="text-xs font-bold px-2 py-0.5 rounded-full" style={{ background: s.color + "20", color: s.color }}>{s.actor}</span>
                <span className="font-bold text-sm" style={{ color: s.color }}>{s.action}</span>
              </div>
              <p className="text-xs" style={{ color: COLORS.gray600 }}>{s.desc}</p>
            </div>
          </div>
        ))}
      </div>

      <div className="rounded-xl p-5 border mt-4" style={{ background: "#F5F3FF", borderColor: "#8B5CF640" }}>
        <h3 className="font-bold text-sm mb-3" style={{ color: "#8B5CF6" }}>Changelog — Solves the Re-read Problem</h3>
        <div className="grid grid-cols-3 gap-3">
          <div className="rounded-lg p-3 bg-white">
            <div className="text-xs font-bold mb-1" style={{ color: COLORS.success }}>🟢 Added</div>
            <p className="text-xs" style={{ color: COLORS.gray500 }}>New modules, apps. Green border + "NEW" badge in PDF.</p>
          </div>
          <div className="rounded-lg p-3 bg-white">
            <div className="text-xs font-bold mb-1" style={{ color: COLORS.amber }}>🟡 Changed</div>
            <p className="text-xs" style={{ color: COLORS.gray500 }}>Values with old→new. Amber highlight + strikethrough old.</p>
          </div>
          <div className="rounded-lg p-3 bg-white">
            <div className="text-xs font-bold mb-1" style={{ color: COLORS.gray400 }}>⚪ Unchanged</div>
            <p className="text-xs" style={{ color: COLORS.gray500 }}>Listed as gray chips. Reviewer skips these.</p>
          </div>
        </div>
      </div>
    </div>
  );
}

// === LICENSE LIFECYCLE TAB ===
function LicenseLifecycleTab() {
  return (
    <div className="space-y-6">
      <div className="text-center mb-6">
        <h2 className="text-2xl font-bold" style={{ color: COLORS.primary }}>License Lifecycle</h2>
      </div>

      <svg viewBox="0 0 800 400" className="w-full" style={{ maxWidth: 800 }}>
        {/* Status nodes */}
        <rect x="50" y="160" width="120" height="50" rx="10" fill={COLORS.amberLight} stroke={COLORS.amber} strokeWidth="2" />
        <text x="110" y="190" textAnchor="middle" fill={COLORS.amber} fontWeight="700" fontSize="13">Trial</text>

        <rect x="250" y="80" width="120" height="50" rx="10" fill={COLORS.successLight} stroke={COLORS.success} strokeWidth="2" />
        <text x="310" y="110" textAnchor="middle" fill="#166534" fontWeight="700" fontSize="13">Active</text>

        <rect x="250" y="240" width="120" height="50" rx="10" fill={COLORS.errorLight} stroke={COLORS.error} strokeWidth="2" />
        <text x="310" y="270" textAnchor="middle" fill={COLORS.error} fontWeight="700" fontSize="13">Suspended</text>

        <rect x="480" y="160" width="120" height="50" rx="10" fill={COLORS.gray100} stroke={COLORS.gray400} strokeWidth="2" />
        <text x="540" y="190" textAnchor="middle" fill={COLORS.gray700} fontWeight="700" fontSize="13">Expired</text>

        <rect x="660" y="160" width="120" height="50" rx="10" fill={COLORS.gray50} stroke={COLORS.gray300} strokeWidth="1.5" strokeDasharray="4,3" />
        <text x="720" y="190" textAnchor="middle" fill={COLORS.gray400} fontWeight="700" fontSize="13">Archived</text>

        {/* Arrows */}
        <defs>
          <marker id="arr" viewBox="0 0 10 10" refX="9" refY="5" markerWidth="6" markerHeight="6" orient="auto"><path d="M 0 0 L 10 5 L 0 10 z" fill={COLORS.gray500} /></marker>
          <marker id="arrG" viewBox="0 0 10 10" refX="9" refY="5" markerWidth="6" markerHeight="6" orient="auto"><path d="M 0 0 L 10 5 L 0 10 z" fill={COLORS.success} /></marker>
          <marker id="arrR" viewBox="0 0 10 10" refX="9" refY="5" markerWidth="6" markerHeight="6" orient="auto"><path d="M 0 0 L 10 5 L 0 10 z" fill={COLORS.error} /></marker>
        </defs>

        {/* Trial → Active */}
        <path d="M170 175 Q210 130 250 110" stroke={COLORS.success} strokeWidth="2" fill="none" markerEnd="url(#arrG)" />
        <text x="195" y="130" fill={COLORS.success} fontSize="9" fontWeight="600">approve</text>

        {/* Active → Suspended */}
        <path d="M290 130 Q270 185 270 240" stroke={COLORS.error} strokeWidth="1.5" fill="none" markerEnd="url(#arrR)" />
        <text x="250" y="190" fill={COLORS.error} fontSize="9" fontWeight="600">suspend</text>

        {/* Suspended → Active */}
        <path d="M350 240 Q370 185 350 130" stroke={COLORS.success} strokeWidth="1.5" fill="none" markerEnd="url(#arrG)" />
        <text x="375" y="190" fill={COLORS.success} fontSize="9" fontWeight="600">unsuspend</text>

        {/* Active → Expired */}
        <path d="M370 105 Q430 140 480 175" stroke={COLORS.gray500} strokeWidth="1.5" fill="none" markerEnd="url(#arr)" />
        <text x="430" y="125" fill={COLORS.gray500} fontSize="9">expires_at / manual</text>

        {/* Expired → Active */}
        <path d="M480 170 Q420 100 370 100" stroke={COLORS.success} strokeWidth="1.5" fill="none" markerEnd="url(#arrG)" />
        <text x="440" y="90" fill={COLORS.success} fontSize="9" fontWeight="600">renew</text>

        {/* Any → Archived */}
        <line x1="600" y1="185" x2="660" y2="185" stroke={COLORS.gray300} strokeWidth="1.5" markerEnd="url(#arr)" strokeDasharray="4,3" />
        <text x="630" y="178" fill={COLORS.gray400} fontSize="9">archive</text>

        {/* Creation paths */}
        <rect x="50" y="20" width="180" height="40" rx="8" fill="#EFF6FF" stroke="#3B82F6" strokeWidth="1" />
        <text x="140" y="45" textAnchor="middle" fill="#3B82F6" fontWeight="600" fontSize="10">Proposal approved → auto-create</text>
        <path d="M140 60 L110 160" stroke="#3B82F6" strokeWidth="1.5" fill="none" markerEnd="url(#arr)" />

        <rect x="250" y="20" width="180" height="40" rx="8" fill={COLORS.amberLight} stroke={COLORS.amber} strokeWidth="1" />
        <text x="340" y="45" textAnchor="middle" fill={COLORS.amber} fontWeight="600" fontSize="10">PO direct create</text>
        <path d="M340 60 L320 80" stroke={COLORS.amber} strokeWidth="1.5" fill="none" markerEnd="url(#arr)" />

        {/* Provision box */}
        <rect x="480" y="60" width="160" height="50" rx="8" fill={COLORS.accentLight} stroke={COLORS.accent} strokeWidth="1.5" />
        <text x="560" y="82" textAnchor="middle" fill={COLORS.accent} fontWeight="700" fontSize="11">Provision</text>
        <text x="560" y="98" textAnchor="middle" fill={COLORS.gray500} fontSize="9">Vernon → calls client app</text>
        <path d="M370 95 L480 85" stroke={COLORS.accent} strokeWidth="1.5" fill="none" markerEnd="url(#arr)" strokeDasharray="4,3" />

        {/* Legend */}
        <rect x="50" y="340" width="700" height="45" rx="8" fill={COLORS.gray50} stroke={COLORS.gray200} strokeWidth="1" />
        <text x="70" y="367" fill={COLORS.gray500} fontSize="10" fontWeight="600">Legend:</text>
        <circle cx="140" cy="362" r="5" fill={COLORS.success} /><text x="150" y="367" fill={COLORS.gray600} fontSize="10">PO/Superuser action</text>
        <circle cx="310" cy="362" r="5" fill={COLORS.error} /><text x="320" y="367" fill={COLORS.gray600} fontSize="10">PO/Superuser action</text>
        <circle cx="470" cy="362" r="5" fill={COLORS.gray400} /><text x="480" y="367" fill={COLORS.gray600} fontSize="10">Automatic / system</text>
        <circle cx="620" cy="362" r="5" fill="#3B82F6" /><text x="630" y="367" fill={COLORS.gray600} fontSize="10">Via proposal</text>
      </svg>
    </div>
  );
}

// === LICENSE CHECK TAB ===
function LicenseCheckTab() {
  const steps = [
    { label: "Vernon provisions license", desc: "Calls client app with license_key, vernon_url, check_interval", icon: "🔗", side: "vernon" },
    { label: "Client app stores config", desc: "Saves license_key + vernon_url + interval locally", icon: "💾", side: "client" },
    { label: "Periodic check (every 1h/6h/24h)", desc: "GET /api/v1/client/license?key=FL-XXX with Origin header", icon: "🔄", side: "client" },
    { label: "Vernon validates", desc: "Check Origin vs instance_url. Return license data + constraints + status", icon: "✓", side: "vernon" },
    { label: "Client enforces", desc: "Update local constraints (max_users, modules, etc). If status ≠ active → DENY ACCESS", icon: "🛡️", side: "client" },
    { label: "Grace period", desc: "If Vernon unreachable for 3× interval → deny access. Prevents offline bypass.", icon: "⏳", side: "client" },
  ];

  return (
    <div className="space-y-6">
      <div className="text-center mb-6">
        <h2 className="text-2xl font-bold" style={{ color: COLORS.primary }}>License Check Protocol</h2>
        <p className="text-xs mt-1" style={{ color: COLORS.gray500 }}>How client apps validate their license against Vernon</p>
      </div>

      <div className="flex gap-4">
        <div className="flex-1">
          <div className="text-center text-xs font-bold mb-3 px-3 py-1.5 rounded-full" style={{ background: COLORS.primary + "10", color: COLORS.primary }}>Vernon License</div>
          {steps.filter(s => s.side === "vernon").map((s, i) => (
            <div key={i} className="rounded-xl p-4 border mb-3" style={{ background: "#F3E8FF", borderColor: COLORS.primary + "30" }}>
              <div className="text-lg mb-1">{s.icon}</div>
              <div className="text-xs font-bold mb-1" style={{ color: COLORS.primary }}>{s.label}</div>
              <div className="text-xs" style={{ color: COLORS.gray500 }}>{s.desc}</div>
            </div>
          ))}
        </div>
        <div className="flex flex-col items-center justify-center" style={{ minWidth: 40 }}>
          <div className="w-0.5 flex-1" style={{ background: COLORS.gray200 }} />
          <div className="text-xs font-bold my-2 px-2 py-1 rounded-full" style={{ background: COLORS.accentLight, color: COLORS.accent }}>HTTP</div>
          <div className="w-0.5 flex-1" style={{ background: COLORS.gray200 }} />
        </div>
        <div className="flex-1">
          <div className="text-center text-xs font-bold mb-3 px-3 py-1.5 rounded-full" style={{ background: COLORS.accent + "10", color: COLORS.accent }}>Client App (any Vernon product)</div>
          {steps.filter(s => s.side === "client").map((s, i) => (
            <div key={i} className="rounded-xl p-4 border mb-3" style={{ background: COLORS.accentLight, borderColor: COLORS.accent + "30" }}>
              <div className="text-lg mb-1">{s.icon}</div>
              <div className="text-xs font-bold mb-1" style={{ color: COLORS.accent }}>{s.label}</div>
              <div className="text-xs" style={{ color: COLORS.gray500 }}>{s.desc}</div>
            </div>
          ))}
        </div>
      </div>

      <div className="rounded-xl p-4 border" style={{ background: COLORS.errorLight, borderColor: COLORS.error + "30" }}>
        <div className="text-xs font-bold mb-1" style={{ color: COLORS.error }}>⛔ Access Denied When:</div>
        <div className="grid grid-cols-3 gap-2 mt-2">
          {["status ≠ active (suspended/expired)", "expires_at < now", "Vernon unreachable for 3× interval"].map((r, i) => (
            <div key={i} className="text-xs rounded-lg p-2 bg-white" style={{ color: COLORS.gray700 }}>{r}</div>
          ))}
        </div>
      </div>
    </div>
  );
}

// === ROLES TAB ===
function RolesTab() {
  const permissions = [
    { feature: "Companies & Projects", sub: "CRUD", sales: true, po: true, su: true },
    { feature: "View licenses", sub: "list, detail, export", sales: true, po: true, su: true },
    { feature: "Create proposal", sub: "draft + submit", sales: true, po: true, su: true },
    { feature: "Download approved PDF", sub: "", sales: true, po: true, su: true },
    { feature: "Edit submitted proposal", sub: "modify pricing, modules", sales: false, po: true, su: true },
    { feature: "Approve / reject proposal", sub: "", sales: false, po: true, su: true },
    { feature: "Preview PDF (any status)", sub: "", sales: false, po: true, su: true },
    { feature: "Create license directly", sub: "skip proposal", sales: false, po: true, su: true },
    { feature: "Suspend / activate", sub: "", sales: false, po: true, su: true },
    { feature: "Renew license", sub: "", sales: false, po: true, su: true },
    { feature: "Provision to client app", sub: "", sales: false, po: true, su: true },
    { feature: "Manage products", sub: "CRUD", sales: false, po: false, su: true },
    { feature: "Manage users", sub: "create PO/sales", sales: false, po: false, su: true },
    { feature: "Global audit log", sub: "", sales: false, po: false, su: true },
  ];

  return (
    <div className="space-y-6">
      <div className="text-center mb-6">
        <h2 className="text-2xl font-bold" style={{ color: COLORS.primary }}>Role-based Access Policies</h2>
      </div>

      <div className="flex gap-3 justify-center mb-4">
        {[
          { role: "Sales", color: "#3B82F6", bg: "#EFF6FF", desc: "Prospecting, proposals" },
          { role: "Project Owner", color: COLORS.amber, bg: COLORS.amberLight, desc: "Approve, license mgmt" },
          { role: "Superuser", color: COLORS.primary, bg: "#F3E8FF", desc: "Full access, system admin" },
        ].map((r) => (
          <div key={r.role} className="rounded-xl p-3 border text-center" style={{ background: r.bg, borderColor: r.color + "30", minWidth: 140 }}>
            <div className="text-sm font-bold" style={{ color: r.color }}>{r.role}</div>
            <div className="text-xs mt-0.5" style={{ color: COLORS.gray500 }}>{r.desc}</div>
          </div>
        ))}
      </div>

      <div className="rounded-xl border overflow-hidden" style={{ borderColor: COLORS.gray200 }}>
        <table className="w-full text-xs">
          <thead>
            <tr style={{ background: COLORS.primary }}>
              <th className="text-left p-3 text-white font-bold">Feature</th>
              <th className="text-center p-3 text-white font-bold" style={{ width: 90 }}>Sales</th>
              <th className="text-center p-3 text-white font-bold" style={{ width: 90 }}>PO</th>
              <th className="text-center p-3 text-white font-bold" style={{ width: 90 }}>Super</th>
            </tr>
          </thead>
          <tbody>
            {permissions.map((p, i) => (
              <tr key={i} style={{ background: i % 2 === 0 ? COLORS.white : COLORS.gray50 }}>
                <td className="p-3">
                  <div className="font-semibold" style={{ color: COLORS.gray900 }}>{p.feature}</div>
                  {p.sub && <div style={{ color: COLORS.gray400 }}>{p.sub}</div>}
                </td>
                <td className="text-center p-3">{p.sales ? <span style={{ color: COLORS.success }}>✅</span> : <span style={{ color: COLORS.gray300 }}>—</span>}</td>
                <td className="text-center p-3">{p.po ? <span style={{ color: COLORS.success }}>✅</span> : <span style={{ color: COLORS.gray300 }}>—</span>}</td>
                <td className="text-center p-3">{p.su ? <span style={{ color: COLORS.success }}>✅</span> : <span style={{ color: COLORS.gray300 }}>—</span>}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}

// === MAIN APP ===
export default function App() {
  const [activeTab, setActiveTab] = useState("overview");

  const renderTab = () => {
    switch (activeTab) {
      case "overview": return <OverviewTab />;
      case "entities": return <EntitiesTab />;
      case "proposal": return <ProposalFlowTab />;
      case "license": return <LicenseLifecycleTab />;
      case "check": return <LicenseCheckTab />;
      case "roles": return <RolesTab />;
      default: return null;
    }
  };

  return (
    <div className="min-h-screen" style={{ background: COLORS.gray100 }}>
      <div className="max-w-4xl mx-auto p-6">
        {/* Header */}
        <div className="flex items-center gap-3 mb-6">
          <div className="w-10 h-10 rounded-xl flex items-center justify-center text-white font-bold text-sm" style={{ background: COLORS.primary }}>VL</div>
          <div>
            <h1 className="text-lg font-bold" style={{ color: COLORS.primary }}>Vernon License</h1>
            <p className="text-xs" style={{ color: COLORS.gray500 }}>Feature, Flow & Policy Visualization</p>
          </div>
        </div>

        {/* Tabs */}
        <div className="flex gap-1 mb-6 p-1 rounded-xl" style={{ background: COLORS.gray200 }}>
          {tabs.map((t) => (
            <button
              key={t.id}
              onClick={() => setActiveTab(t.id)}
              className="flex-1 py-2 px-3 rounded-lg text-xs font-semibold transition-all"
              style={{
                background: activeTab === t.id ? COLORS.white : "transparent",
                color: activeTab === t.id ? COLORS.primary : COLORS.gray500,
                boxShadow: activeTab === t.id ? "0 1px 3px rgba(0,0,0,0.1)" : "none",
              }}
            >
              {t.label}
            </button>
          ))}
        </div>

        {/* Content */}
        <div className="rounded-2xl p-6" style={{ background: COLORS.white, boxShadow: "0 1px 3px rgba(0,0,0,0.08)" }}>
          {renderTab()}
        </div>
      </div>
    </div>
  );
}
