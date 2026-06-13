import { useState, useEffect, useCallback } from 'react';
import StatCard from '../components/StatCard';
import NetworkGraph from '../components/NetworkGraph';
import MerchantDetail from '../components/MerchantDetail';
import RiskBadge, { riskColor } from '../components/RiskBadge';
import { getDashboardStats, getMerchantNetwork, getReports } from '../lib/api';

export default function Dashboard() {
  const [stats, setStats] = useState(null);
  const [merchants, setMerchants] = useState([]);
  const [reports, setReports] = useState([]);
  const [selected, setSelected] = useState(null);
  const [liveCount, setLiveCount] = useState(0);
  const [loading, setLoading] = useState(true);

  const fetchData = useCallback(async () => {
    try {
      const [s, m, r] = await Promise.all([
        getDashboardStats(),
        getMerchantNetwork(),
        getReports(),
      ]);
      setStats(s);
      setMerchants(m);
      setReports(r);
      setLiveCount(s.reports_today);
    } catch (e) {
      console.error('Dashboard fetch error:', e);
    }
    setLoading(false);
  }, []);

  useEffect(() => { fetchData(); }, [fetchData]);

  useEffect(() => {
    const interval = setInterval(() => {
      if (Math.random() > 0.7) setLiveCount((c) => c + 1);
    }, 4000);
    return () => clearInterval(interval);
  }, []);

  const selectedMerchant = merchants.find((m) => m.id === selected);

  if (loading) {
    return <div className="text-center font-mono text-[rgba(255,255,255,0.4)] py-10">Loading...</div>;
  }

  return (
    <div className="animate-[fade-up_0.4s_ease-out]">
      {/* Stats */}
      <div className="flex gap-3 mb-6 flex-wrap">
        <StatCard label="Total Terflag" value={stats.total_flagged} accent="#2B2BFF" />
        <StatCard label="Risiko Tinggi" value={stats.high_risk_count} accent="#FF2D55" />
        <StatCard label="Laporan Hari Ini" value={liveCount} accent="#FF9F0A" />
        <StatCard label="Pending OJK" value={stats.pending_ojk} accent="#2B2BFF" />
      </div>

      {/* Network + Detail */}
      <div className="grid grid-cols-1 lg:grid-cols-[1fr_300px] gap-4 mb-4">
        <div className="bg-[#0D0D1F] border border-[rgba(43,43,255,0.2)] rounded-lg overflow-hidden relative">
          <div className="absolute top-3 left-4 z-10 font-mono text-[10px] font-bold text-[rgba(255,255,255,0.4)] uppercase tracking-widest">
            JARINGAN MERCHANT — klik node
          </div>
          <div className="h-[420px]">
            <NetworkGraph merchants={merchants} selected={selected} onSelect={setSelected} />
          </div>
          <div className="px-4 py-2 border-t border-[rgba(43,43,255,0.2)] flex gap-5">
            {[['#FF2D55', 'Risiko Tinggi (≥70)'], ['#FF9F0A', 'Waspada (30–69)'], ['#30D158', 'Aman (<30)']].map(([c, l]) => (
              <div key={l} className="flex items-center gap-1.5">
                <div className="w-2 h-2 rounded-full" style={{ background: c }} />
                <span className="font-mono text-[10px] text-[rgba(255,255,255,0.4)]">{l}</span>
              </div>
            ))}
          </div>
        </div>

        <MerchantDetail merchant={selectedMerchant} onClose={() => setSelected(null)} />
      </div>

      {/* Recent Reports */}
      <div className="bg-[#0D0D1F] border border-[rgba(43,43,255,0.2)] rounded-lg overflow-hidden">
        <div className="px-5 py-3 border-b border-[rgba(43,43,255,0.2)]">
          <span className="font-mono text-[10px] font-bold text-[rgba(255,255,255,0.4)] uppercase tracking-widest">LAPORAN TERBARU</span>
        </div>
        {reports.map((r, i) => (
          <div
            key={r.id}
            className="flex items-center px-5 gap-4 transition-colors"
            style={{
              height: 52,
              borderBottom: i < reports.length - 1 ? '1px solid rgba(255,255,255,0.04)' : 'none',
            }}
            onMouseEnter={(e) => { e.currentTarget.style.background = 'rgba(43,43,255,0.05)'; }}
            onMouseLeave={(e) => { e.currentTarget.style.background = ''; }}
          >
            <div className="font-mono text-[11px] text-[rgba(255,255,255,0.4)] w-20">{r.id}</div>
            <div className="flex-1 min-w-0">
              <div className="font-body text-[13px] text-white truncate">{r.merchant_name || r.merchant_id}</div>
              <div className="font-mono text-[10px] text-[rgba(255,255,255,0.25)] mt-0.5">{r.merchant_id} · {r.time}</div>
            </div>
            <RiskBadge score={r.score} />
          </div>
        ))}
        {reports.length === 0 && (
          <div className="py-6 text-center font-mono text-xs text-[rgba(255,255,255,0.25)]">Tidak ada laporan hari ini.</div>
        )}
      </div>
    </div>
  );
}
