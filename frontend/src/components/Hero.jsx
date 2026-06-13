import { useState, useEffect } from 'react';
import { getDashboardStats } from '../lib/api';
import Illustration from './Illustration';

export default function Hero() {
  const [stats, setStats] = useState(null);

  useEffect(() => {
    getDashboardStats().then(setStats).catch(() => {});
  }, []);

  return (
    <section className="w-full" style={{ background: '#2B2BFF' }}>
      <div className="max-w-[1200px] mx-auto px-6 py-[80px]">
        <div className="flex flex-col lg:flex-row items-center lg:items-start gap-12 lg:gap-0">
          {/* Text side — 45% */}
          <div className="w-full lg:w-[45%] flex flex-col justify-center">
            <div className="font-mono text-[11px] font-bold text-[rgba(255,255,255,0.55)] uppercase mb-4"
              style={{ letterSpacing: '0.2em' }}>
              OPEN SOURCE · OJK
            </div>
            <h1 className="font-display font-black text-white uppercase leading-[0.92]"
              style={{
                fontSize: 'clamp(48px, 7vw, 96px)',
                letterSpacing: '-0.02em',
              }}>
              THE FRAUD<br />INTELLIGENCE<br />SYSTEM
            </h1>
            <button
              className="inline-flex items-center gap-2 font-mono text-xs font-bold uppercase tracking-wider px-6 py-3 mt-8 rounded-sm border-none cursor-pointer"
              style={{
                background: '#FFFFFF',
                color: '#2B2BFF',
                letterSpacing: '0.15em',
                alignSelf: 'flex-start',
                transition: 'transform 0.15s ease, box-shadow 0.15s ease',
              }}
              onMouseEnter={(e) => { e.currentTarget.style.transform = 'translateY(-1px)'; e.currentTarget.style.boxShadow = '0 8px 24px rgba(0,0,0,0.3)'; }}
              onMouseLeave={(e) => { e.currentTarget.style.transform = ''; e.currentTarget.style.boxShadow = ''; }}
            >
              SCAN MERCHANT →
            </button>
          </div>

          {/* Illustration side — 55% */}
          <div className="w-full lg:w-[55%] flex justify-center lg:justify-end">
            <div className="w-[320px] h-[340px] sm:w-[400px] sm:h-[420px]">
              <Illustration />
            </div>
          </div>
        </div>

        {/* Stat bar */}
        {stats && (
          <div className="mt-16 pt-6 border-t border-[rgba(255,255,255,0.15)] flex flex-wrap gap-x-10 gap-y-3">
            <div className="flex items-center gap-3">
              <span className="font-mono text-[28px] font-bold text-white">{stats.total_flagged.toLocaleString()}</span>
              <span className="font-mono text-[11px] font-bold text-[rgba(255,255,255,0.55)] uppercase tracking-widest">TERFLAG</span>
            </div>
            <div className="flex items-center gap-3">
              <span className="font-mono text-[28px] font-bold text-[#FF2D55]">{stats.high_risk_count.toLocaleString()}</span>
              <span className="font-mono text-[11px] font-bold text-[rgba(255,255,255,0.55)] uppercase tracking-widest">RISIKO TINGGI</span>
            </div>
            <div className="flex items-center gap-3">
              <span className="font-mono text-[28px] font-bold text-[#FF9F0A]">{stats.pending_ojk.toLocaleString()}</span>
              <span className="font-mono text-[11px] font-bold text-[rgba(255,255,255,0.55)] uppercase tracking-widest">→ OJK</span>
            </div>
          </div>
        )}
      </div>
    </section>
  );
}
