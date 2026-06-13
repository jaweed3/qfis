import { useState } from 'react';
import Hero from './components/Hero';
import Dashboard from './pages/Dashboard';
import CheckMerchant from './pages/CheckMerchant';
import Reports from './pages/Reports';

const tabs = [
  { id: 'dashboard', label: 'DASHBOARD' },
  { id: 'check', label: 'CEK MERCHANT' },
  { id: 'reports', label: 'LAPORAN' },
];

export default function App() {
  const [activeTab, setActiveTab] = useState('dashboard');

  return (
    <div className="min-h-screen" style={{ background: '#07070F', color: '#FFFFFF' }}>
      <style>{`
        @import url('https://fonts.googleapis.com/css2?family=Playfair+Display:ital,wght@0,700;0,900;1,700&family=JetBrains+Mono:wght@400;600;700&family=Inter:wght@400;500&display=swap');
      `}</style>

      {/* Nav — Klein Blue translucent */}
      <nav className="h-14 px-6 flex items-center justify-between sticky top-0 z-100 border-b border-[rgba(255,255,255,0.1)]"
        style={{ background: 'rgba(43,43,255,0.95)', backdropFilter: 'blur(12px)' }}>
        <div className="flex items-center gap-3">
          <div className="flex items-center gap-2.5">
            <div className="w-2 h-2 rounded-full bg-white animate-[live-pulse_2s_ease-in-out_infinite]" />
            <span className="font-display font-bold text-lg text-white tracking-tight"
              style={{ fontSize: 18 }}>
              QFIS
            </span>
          </div>
        </div>

        <div className="flex gap-1">
          {tabs.map((t) => (
            <button
              key={t.id}
              onClick={() => setActiveTab(t.id)}
              className="font-mono text-[11px] font-bold uppercase tracking-wider px-3.5 py-1.5 rounded-sm cursor-pointer transition-all"
              style={{
                color: activeTab === t.id ? 'white' : 'rgba(255,255,255,0.45)',
                background: activeTab === t.id ? 'rgba(255,255,255,0.1)' : 'transparent',
                letterSpacing: '0.15em',
              }}
            >
              {t.label}
            </button>
          ))}
        </div>

        <div>
          <button
            className="font-mono text-[11px] font-bold uppercase tracking-wider px-4 py-1.5 rounded-sm cursor-pointer transition-all"
            style={{
              border: '1px solid rgba(255,255,255,0.4)',
              color: 'white',
              background: 'transparent',
              letterSpacing: '0.15em',
            }}
            onMouseEnter={(e) => { e.currentTarget.style.background = 'white'; e.currentTarget.style.color = '#2B2BFF'; }}
            onMouseLeave={(e) => { e.currentTarget.style.background = 'transparent'; e.currentTarget.style.color = 'white'; }}
          >
            INSTALL →
          </button>
        </div>
      </nav>

      {/* Hero — full-bleed Klein Blue */}
      <Hero />

      {/* Dashboard section */}
      <main className="max-w-[1200px] mx-auto px-6 py-8">
        {activeTab === 'dashboard' && <Dashboard />}
        {activeTab === 'check' && <CheckMerchant />}
        {activeTab === 'reports' && <Reports />}
      </main>

      {/* Footer */}
      <footer className="border-t border-[rgba(43,43,255,0.2)] px-6 py-4 flex justify-between items-center max-w-[1200px] mx-auto">
        <span className="font-mono text-[10px] font-bold text-[rgba(255,255,255,0.25)] uppercase tracking-widest">QFIS v0.1 — Hackathon Build</span>
        <span className="font-body text-[10px] text-[rgba(255,255,255,0.25)]">Data dari sumber publik OJK · Komdigi · OSS</span>
      </footer>
    </div>
  );
}
