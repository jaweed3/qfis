import { useState, useEffect, useCallback } from 'react';
import RiskBadge from '../components/RiskBadge';
import ReportForm from '../components/ReportForm';
import { getReports } from '../lib/api';

export default function Reports() {
  const [reports, setReports] = useState([]);

  const fetchReports = useCallback(async () => {
    try {
      const data = await getReports();
      setReports(data);
    } catch (e) {
      console.error('Failed to fetch reports:', e);
    }
  }, []);

  useEffect(() => { fetchReports(); }, [fetchReports]);

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 gap-4 max-w-[900px] mx-auto animate-[fade-up_0.4s_ease-out]">
      <div>
        <div className="font-mono text-[10px] font-bold text-[rgba(255,255,255,0.4)] uppercase mb-3 tracking-widest">
          ANTRIAN OJK
        </div>
        <div className="bg-[#0D0D1F] border border-[rgba(43,43,255,0.2)] rounded-lg overflow-hidden">
          {reports.length === 0 && (
            <div className="py-6 text-center font-mono text-xs text-[rgba(255,255,255,0.25)]">Belum ada laporan.</div>
          )}
          {reports.map((r, i) => (
            <div
              key={r.id}
              className="flex justify-between items-center px-4.5 py-3.5 transition-colors"
              style={{ borderBottom: i < reports.length - 1 ? '1px solid rgba(255,255,255,0.04)' : 'none' }}
              onMouseEnter={(e) => { e.currentTarget.style.background = 'rgba(43,43,255,0.05)'; }}
              onMouseLeave={(e) => { e.currentTarget.style.background = ''; }}
            >
              <div>
                <div className="font-body text-[13px] text-white">{r.merchant_name || r.merchant_id}</div>
                <div className="font-mono text-[10px] text-[rgba(255,255,255,0.25)] mt-0.5">{r.merchant_id} · {r.time}</div>
              </div>
              <RiskBadge score={r.score} />
            </div>
          ))}
        </div>
      </div>

      <div>
        <div className="font-mono text-[10px] font-bold text-[rgba(255,255,255,0.4)] uppercase mb-3 tracking-widest">
          LAPORKAN TEMUAN
        </div>
        <ReportForm onSubmitted={fetchReports} />
        <div className="mt-3 px-4 py-3 bg-[#0D0D1F] border border-[rgba(43,43,255,0.2)] rounded-lg">
          <div className="font-body text-xs text-[rgba(255,255,255,0.4)] leading-relaxed">
            Laporan dengan skor ≥70 diteruskan otomatis ke OJK. Data reporter bersifat anonim.
          </div>
        </div>
      </div>
    </div>
  );
}
